package handlers

import (
	"diary/internal/models"
	"diary/internal/services"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/supertokens/supertokens-golang/recipe/session"
)

// --- Entry Handler Interface ---

type EntryHandler interface {
	CreateEntry(w http.ResponseWriter, r *http.Request)
	GetEntry(w http.ResponseWriter, r *http.Request)
	UpdateEntry(w http.ResponseWriter, r *http.Request)
	DeleteEntry(w http.ResponseWriter, r *http.Request)
	ListEntries(w http.ResponseWriter, r *http.Request)
}

// --- Entry Handler Implementation ---

type entryHandler struct {
	service services.EntryService
}

func NewEntryHandler(service services.EntryService) EntryHandler {
	return &entryHandler{service: service}
}

// --- Request/Response Structs ---

type EntryRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type EntryResponse struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
}

// --- Entry Handlers ---

func (h *entryHandler) CreateEntry(w http.ResponseWriter, r *http.Request) {
	// Получаем userID из сессии
	sessionContainer := session.GetSessionFromRequestContext(r.Context())
	userIDStr := sessionContainer.GetUserID()
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Декодируем запрос
	var req EntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Создаем запись
	entry := &models.Entry{
		UserID:  userID,
		Title:   req.Title,
		Content: req.Content,
	}

	if err := h.service.CreateEntry(entry); err != nil {
		http.Error(w, "Failed to create entry", http.StatusInternalServerError)
		return
	}

	// Возвращаем ответ
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]interface{}{
		"id":      entry.ID.String(),
		"message": "Entry created successfully",
	})
}

func (h *entryHandler) GetEntry(w http.ResponseWriter, r *http.Request) {
	// Получаем ID записи из URL
	entryID := chi.URLParam(r, "id")

	// Получаем запись
	entry, err := h.service.GetEntryByID(entryID)
	if err != nil {
		http.Error(w, "Entry not found", http.StatusNotFound)
		return
	}

	// Проверка доступа (только владелец может видеть свою запись)
	sessionContainer := session.GetSessionFromRequestContext(r.Context())
	userIDStr := sessionContainer.GetUserID()
	if entry.UserID.String() != userIDStr {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	// Формируем ответ
	response := EntryResponse{
		ID:        entry.ID.String(),
		UserID:    entry.UserID.String(),
		Title:     entry.Title,
		Content:   entry.Content,
		CreatedAt: entry.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	render.JSON(w, r, response)
}

func (h *entryHandler) UpdateEntry(w http.ResponseWriter, r *http.Request) {
	// Получаем ID записи из URL
	entryID := chi.URLParam(r, "id")

	// Получаем существующую запись
	existingEntry, err := h.service.GetEntryByID(entryID)
	if err != nil {
		http.Error(w, "Entry not found", http.StatusNotFound)
		return
	}

	// Проверка доступа (только владелец может обновлять свою запись)
	sessionContainer := session.GetSessionFromRequestContext(r.Context())
	userIDStr := sessionContainer.GetUserID()
	if existingEntry.UserID.String() != userIDStr {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	// Декодируем запрос
	var req EntryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Обновляем запись
	existingEntry.Title = req.Title
	existingEntry.Content = req.Content

	if err := h.service.UpdateEntry(existingEntry); err != nil {
		http.Error(w, "Failed to update entry", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, map[string]interface{}{
		"message": "Entry updated successfully",
	})
}

func (h *entryHandler) DeleteEntry(w http.ResponseWriter, r *http.Request) {
	// Получаем ID записи из URL
	entryID := chi.URLParam(r, "id")

	// Получаем существующую запись для проверки владельца
	existingEntry, err := h.service.GetEntryByID(entryID)
	if err != nil {
		http.Error(w, "Entry not found", http.StatusNotFound)
		return
	}

	// Проверка доступа (только владелец может удалять свою запись)
	sessionContainer := session.GetSessionFromRequestContext(r.Context())
	userIDStr := sessionContainer.GetUserID()
	if existingEntry.UserID.String() != userIDStr {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	// Удаляем запись
	if err := h.service.DeleteEntry(entryID); err != nil {
		http.Error(w, "Failed to delete entry", http.StatusInternalServerError)
		return
	}

	render.JSON(w, r, map[string]interface{}{
		"message": "Entry deleted successfully",
	})
}

func (h *entryHandler) ListEntries(w http.ResponseWriter, r *http.Request) {
	// Получаем список всех записей
	entries, err := h.service.ListEntries()
	if err != nil {
		http.Error(w, "Failed to retrieve entries", http.StatusInternalServerError)
		return
	}

	// Получаем userID для фильтрации записей (пользователь должен видеть только свои записи)
	sessionContainer := session.GetSessionFromRequestContext(r.Context())
	userIDStr := sessionContainer.GetUserID()

	// Фильтруем записи и формируем ответ
	var response []EntryResponse
	for _, entry := range entries {
		if entry.UserID.String() == userIDStr {
			response = append(response, EntryResponse{
				ID:        entry.ID.String(),
				UserID:    entry.UserID.String(),
				Title:     entry.Title,
				Content:   entry.Content,
				CreatedAt: entry.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			})
		}
	}

	render.JSON(w, r, response)
}
