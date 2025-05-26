package services

import (
	"diary/internal/models"
	"diary/internal/repos"
)

// Главный интерфейс объединяет все подсервисы
type Service interface {
	EntryService
}

// --- Комбинирующий сервис ---

type service struct {
	entryService EntryService
}

// Прокси-методы EntryService

func (s *service) CreateEntry(entry *models.Entry) error {
	return s.entryService.CreateEntry(entry)
}

func (s *service) GetEntryByID(id string) (*models.Entry, error) {
	return s.entryService.GetEntryByID(id)
}

func (s *service) UpdateEntry(entry *models.Entry) error {
	return s.entryService.UpdateEntry(entry)
}

func (s *service) DeleteEntry(id string) error {
	return s.entryService.DeleteEntry(id)
}

func (s *service) ListEntries() ([]*models.Entry, error) {
	return s.entryService.ListEntries()
}

// --- Конструктор комбинирующего сервиса ---

func NewService(repo repos.Repository) Service {
	return &service{
		entryService: NewEntryService(repo),
	}
}
