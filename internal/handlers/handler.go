package handlers

import (
	"diary/internal/services"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/supertokens/supertokens-golang/recipe/session"
)

// Главный интерфейс объединяет все подобработчики
type Handler interface {
	EntryHandler
	RegisterRoutes(r *chi.Mux)
}

// --- Комбинирующий обработчик ---

type handler struct {
	entryHandler EntryHandler
}

// Регистрация маршрутов для всего приложения
func (h *handler) RegisterRoutes(r *chi.Mux) {
	r.Route("/api/entries", func(r chi.Router) {
		// Все маршруты требуют аутентификации
		// Заворачиваем каждый маршрут в session.VerifySession для проверки авторизации
		r.Post("/", session.VerifySession(nil, h.CreateEntry))
		r.Get("/{id}", session.VerifySession(nil, h.GetEntry))
		r.Put("/{id}", session.VerifySession(nil, h.UpdateEntry))
		r.Delete("/{id}", session.VerifySession(nil, h.DeleteEntry))
		r.Get("/", session.VerifySession(nil, h.ListEntries))
	})
}

// Прокси-методы EntryHandler

func (h *handler) CreateEntry(w http.ResponseWriter, r *http.Request) {
	h.entryHandler.CreateEntry(w, r)
}

func (h *handler) GetEntry(w http.ResponseWriter, r *http.Request) {
	h.entryHandler.GetEntry(w, r)
}

func (h *handler) UpdateEntry(w http.ResponseWriter, r *http.Request) {
	h.entryHandler.UpdateEntry(w, r)
}

func (h *handler) DeleteEntry(w http.ResponseWriter, r *http.Request) {
	h.entryHandler.DeleteEntry(w, r)
}

func (h *handler) ListEntries(w http.ResponseWriter, r *http.Request) {
	h.entryHandler.ListEntries(w, r)
}

// --- Конструктор комбинирующего обработчика ---

func NewHandler(service services.Service) Handler {
	return &handler{
		entryHandler: NewEntryHandler(service),
	}
}
