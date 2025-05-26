package repos

import (
	"diary/internal/models"

	"gorm.io/gorm"
)

// --- Entry Repository Interface ---

type EntryRepository interface {
	Create(entry *models.Entry) error
	Read(id string) (*models.Entry, error)
	Update(entry *models.Entry) error
	Delete(id string) error
	List() ([]*models.Entry, error)
}

// --- Entry Repository Implementation ---

type entryRepository struct {
	db *gorm.DB
}

func NewEntryRepository(db *gorm.DB) EntryRepository {
	return &entryRepository{db: db}
}

// --- CRUD Entry ---

func (r *entryRepository) Create(entry *models.Entry) error {
	return r.db.Create(entry).Error
}

func (r *entryRepository) Read(id string) (*models.Entry, error) {
	var entry models.Entry
	if err := r.db.First(&entry, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &entry, nil
}

func (r *entryRepository) Update(entry *models.Entry) error {
	return r.db.Save(entry).Error
}

func (r *entryRepository) Delete(id string) error {
	return r.db.Delete(&models.Entry{}, "id = ?", id).Error
}

func (r *entryRepository) List() ([]*models.Entry, error) {
	var entries []*models.Entry
	if err := r.db.Order("created_at DESC").Find(&entries).Error; err != nil {
		return nil, err
	}
	return entries, nil
}
