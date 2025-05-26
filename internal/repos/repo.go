package repos

import (
	"diary/internal/models"

	"gorm.io/gorm"
)

// Главный интерфейс объединяет все подрепозитории
type Repository interface {
	EntryRepository
}

// --- Комбинирующий репозиторий ---

type repository struct {
	entryRepo EntryRepository
}

// Прокси-методы EntryRepository

func (r *repository) Create(entry *models.Entry) error {
	return r.entryRepo.Create(entry)
}

func (r *repository) Read(id string) (*models.Entry, error) {
	return r.entryRepo.Read(id)
}

func (r *repository) Update(entry *models.Entry) error {
	return r.entryRepo.Update(entry)
}

func (r *repository) Delete(id string) error {
	return r.entryRepo.Delete(id)
}

func (r *repository) List() ([]*models.Entry, error) {
	return r.entryRepo.List()
}

// --- Конструктор комбинирующего репозитория ---

func NewRepository(db *gorm.DB) Repository {
	return &repository{
		entryRepo: NewEntryRepository(db),
	}
}
