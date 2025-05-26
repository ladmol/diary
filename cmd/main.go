package main

import (
	"context"
	"diary/internal/models"
	"diary/internal/router"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {	// Инициализация подключения к базе данных
	dsn := getEnv("DATABASE_URL", "host=localhost user=postgres password=postgres dbname=diary port=5432 sslmode=disable TimeZone=Europe/Moscow")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных: %v", err)
	}
	
	// Автоматическая миграция схемы базы данных с помощью GORM
	log.Println("Запуск автоматической миграции схемы...")
	if err := db.AutoMigrate(&models.Entry{}); err != nil {
		log.Fatalf("Ошибка при миграции схемы: %v", err)
	}
	log.Println("Миграция схемы успешно завершена")

	// Инициализация SuperTokens
	err = supertokens.Init(supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: getEnv("SUPERTOKENS_CONNECTION_URI", "http://localhost:3567"),
			APIKey:        getEnv("SUPERTOKENS_API_KEY", ""),
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "Дневник",
			APIDomain:     getEnv("API_DOMAIN", "http://localhost:8080"),
			WebsiteDomain: getEnv("WEBSITE_DOMAIN", "http://localhost:3000"),
		},
		RecipeList: []supertokens.Recipe{
			session.Init(&sessmodels.TypeInput{
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
			// Можно добавить другие рецепты SuperTokens (emailpassword, thirdparty и т.д.)
		},
	})
	if err != nil {
		log.Fatalf("Ошибка инициализации SuperTokens: %v", err)
	}

	// Настройка маршрутов
	r := router.SetupRouter(db)

	// Настройка и запуск сервера
	port := getEnv("PORT", "8080")
	server := router.InitializeServer(r, fmt.Sprintf(":%s", port))

	// Graceful shutdown
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("Сервер запущен на порту %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Ошибка запуска сервера: %v", err)
		}
	}()

	// Ожидание сигнала завершения
	<-done
	log.Print("Завершение работы сервера...")

	// Даем 5 секунд для завершения текущих запросов
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Ошибка завершения работы сервера: %v", err)
	}

	log.Print("Сервер успешно завершил работу")
}

// getEnv получает значение переменной окружения или возвращает значение по умолчанию
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
