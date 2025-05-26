package router

import (
	"diary/internal/handlers"
	"diary/internal/repos"
	"diary/internal/services"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"
)

// SetupRouter настраивает все маршруты приложения и возвращает готовый к использованию роутер
func SetupRouter(db *gorm.DB) *chi.Mux {
	// Создаем экземпляр chi роутера
	r := chi.NewRouter()

	// Добавляем базовые middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Настройка CORS, если необходимо
	// r.Use(cors.Handler(cors.Options{
	//     AllowedOrigins:   []string{"https://*", "http://*"},
	//     AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	//     AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
	//     ExposedHeaders:   []string{"Link"},
	//     AllowCredentials: false,
	//     MaxAge:           300,
	// }))

	// Инициализация репозиториев, сервисов и обработчиков
	repository := repos.NewRepository(db)
	service := services.NewService(repository)
	handler := handlers.NewHandler(service)

	// Регистрация API маршрутов
	handler.RegisterRoutes(r)

	// Можно добавить статический файловый сервер, если он нужен
	// fileServer := http.FileServer(http.Dir("./static"))
	// r.Handle("/static/*", http.StripPrefix("/static", fileServer))

	// Маршрут для проверки здоровья сервера
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	return r
}

// InitializeServer создает и настраивает HTTP сервер с заданным роутером
func InitializeServer(router *chi.Mux, address string) *http.Server {
	return &http.Server{
		Addr:    address,
		Handler: router,
	}
}
