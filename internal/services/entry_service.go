package services

import (
	"diary/internal/models"
	"diary/internal/repos"

	"github.com/google/uuid"
)

// --- Entry Service Interface ---

type EntryService interface {
	CreateEntry(entry *models.Entry) error
	GetEntryByID(id string) (*models.Entry, error)
	UpdateEntry(entry *models.Entry) error
	DeleteEntry(id string) error
	ListEntries() ([]*models.Entry, error)
}

// --- Entry Service Implementation ---

type entryService struct {
	repo repos.EntryRepository
}

func NewEntryService(repo repos.EntryRepository) EntryService {
	return &entryService{repo: repo}
}

// --- Business Logic Entry ---

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
