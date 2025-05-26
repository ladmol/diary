package repos

import (
	"diary/internal/models"

	"gorm.io/gorm"
)

// --- Интерфейсы ---

type EntryRepository interface {
	Create(entry *models.Entry) error
	Read(id string) (*models.Entry, error)
	Update(entry *models.Entry) error
	Delete(id string) error
	List() ([]*models.Entry, error)
}

// Главный интерфейс объединяет все подрепозитории
type Repository interface {
	EntryRepository
}

// --- Структуры реализации ---

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
