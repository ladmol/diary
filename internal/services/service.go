package services

import (
	"diary/internal/models"
	"diary/internal/repos"

	"github.com/google/uuid"
)

// --- Интерфейсы ---

type EntryService interface {
	CreateEntry(entry *models.Entry) error
	GetEntryByID(id string) (*models.Entry, error)
	UpdateEntry(entry *models.Entry) error
	DeleteEntry(id string) error
	ListEntries() ([]*models.Entry, error)
}

// Главный интерфейс объединяет все подсервисы
type Service interface {
	EntryService
}

// --- Структуры реализации ---

type entryService struct {
	repo repos.EntryRepository
}

func NewEntryService(repo repos.EntryRepository) EntryService {
	return &entryService{repo: repo}
}

// --- Бизнес-логика Entry ---

func (s *entryService) CreateEntry(entry *models.Entry) error {
	if entry.ID == uuid.Nil {
		entry.ID = uuid.New()
	}
	return s.repo.Create(entry)
}

func (s *entryService) GetEntryByID(id string) (*models.Entry, error) {
	return s.repo.Read(id)
}

func (s *entryService) UpdateEntry(entry *models.Entry) error {
	return s.repo.Update(entry)
}

func (s *entryService) DeleteEntry(id string) error {
	return s.repo.Delete(id)
}

func (s *entryService) ListEntries() ([]*models.Entry, error) {
	return s.repo.List()
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
