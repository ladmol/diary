package repos

import (
	"diary/internal/models"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type EntryRepositoryTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo EntryRepository
}

func (suite *EntryRepositoryTestSuite) SetupSuite() {
	// Создаем in-memory SQLite базу данных для тестов
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Require().NoError(err)

	// Автомиграция
	err = db.AutoMigrate(&models.Entry{})
	suite.Require().NoError(err)

	suite.db = db
	suite.repo = NewEntryRepository(db)
}

func (suite *EntryRepositoryTestSuite) TearDownTest() {
	// Очищаем таблицу после каждого теста
	suite.db.Exec("DELETE FROM entries")
}

func (suite *EntryRepositoryTestSuite) TestCreate() {
	// Arrange
	entry := &models.Entry{
		ID:        uuid.New(),
		UserID:    uuid.New(),
		Title:     "Test Title",
		Content:   "Test Content",
		CreatedAt: time.Now(),
	}

	// Act
	err := suite.repo.Create(entry)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotEqual(suite.T(), uuid.Nil, entry.ID)

	// Проверяем, что запись действительно сохранена в БД
	var savedEntry models.Entry
	err = suite.db.First(&savedEntry, "id = ?", entry.ID.String()).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), entry.Title, savedEntry.Title)
	assert.Equal(suite.T(), entry.Content, savedEntry.Content)
	assert.Equal(suite.T(), entry.UserID, savedEntry.UserID)
}

func (suite *EntryRepositoryTestSuite) TestCreateWithEmptyTitle() {
	// Arrange
	entry := &models.Entry{
		ID:      uuid.New(),
		UserID:  uuid.New(),
		Title:   "", // Пустой заголовок
		Content: "Test Content",
	}

	// Act & Assert
	// Поскольку в модели поле title помечено как not null,
	// создание записи с пустым title может завершиться ошибкой в зависимости от настроек БД
	err := suite.repo.Create(entry)
	// В SQLite пустая строка допустима, но в реальной БД может быть ограничение
	assert.NoError(suite.T(), err)
}

func (suite *EntryRepositoryTestSuite) TestRead() {
	// Arrange - создаем тестовую запись
	entry := &models.Entry{
		ID:      uuid.New(),
		UserID:  uuid.New(),
		Title:   "Test Title",
		Content: "Test Content",
	}
	err := suite.db.Create(entry).Error
	suite.Require().NoError(err)

	// Act
	foundEntry, err := suite.repo.Read(entry.ID.String())

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), foundEntry)
	assert.Equal(suite.T(), entry.ID, foundEntry.ID)
	assert.Equal(suite.T(), entry.Title, foundEntry.Title)
	assert.Equal(suite.T(), entry.Content, foundEntry.Content)
	assert.Equal(suite.T(), entry.UserID, foundEntry.UserID)
}

func (suite *EntryRepositoryTestSuite) TestReadNotFound() {
	// Arrange
	nonExistentID := uuid.New().String()

	// Act
	foundEntry, err := suite.repo.Read(nonExistentID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), foundEntry)
	assert.Equal(suite.T(), gorm.ErrRecordNotFound, err)
}

func (suite *EntryRepositoryTestSuite) TestReadInvalidID() {
	// Arrange
	invalidID := "invalid-uuid"

	// Act
	foundEntry, err := suite.repo.Read(invalidID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), foundEntry)
}

func (suite *EntryRepositoryTestSuite) TestUpdate() {
	// Arrange - создаем тестовую запись
	entry := &models.Entry{
		ID:      uuid.New(),
		UserID:  uuid.New(),
		Title:   "Original Title",
		Content: "Original Content",
	}
	err := suite.db.Create(entry).Error
	suite.Require().NoError(err)

	// Обновляем данные
	entry.Title = "Updated Title"
	entry.Content = "Updated Content"

	// Act
	err = suite.repo.Update(entry)

	// Assert
	assert.NoError(suite.T(), err)

	// Проверяем, что изменения сохранены
	var updatedEntry models.Entry
	err = suite.db.First(&updatedEntry, "id = ?", entry.ID.String()).Error
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Updated Title", updatedEntry.Title)
	assert.Equal(suite.T(), "Updated Content", updatedEntry.Content)
}

func (suite *EntryRepositoryTestSuite) TestUpdateNonExistent() {
	// Arrange
	entry := &models.Entry{
		ID:      uuid.New(),
		UserID:  uuid.New(),
		Title:   "Test Title",
		Content: "Test Content",
	}

	// Act - пытаемся обновить несуществующую запись
	err := suite.repo.Update(entry)

	// Assert
	// GORM создаст новую запись если ID не найден при Save()
	assert.NoError(suite.T(), err)

	// Проверяем, что запись была создана
	var foundEntry models.Entry
	err = suite.db.First(&foundEntry, "id = ?", entry.ID.String()).Error
	assert.NoError(suite.T(), err)
}

func (suite *EntryRepositoryTestSuite) TestDelete() {
	// Arrange - создаем тестовую запись
	entry := &models.Entry{
		ID:      uuid.New(),
		UserID:  uuid.New(),
		Title:   "Test Title",
		Content: "Test Content",
	}
	err := suite.db.Create(entry).Error
	suite.Require().NoError(err)

	// Act
	err = suite.repo.Delete(entry.ID.String())

	// Assert
	assert.NoError(suite.T(), err)

	// Проверяем, что запись удалена
	var deletedEntry models.Entry
	err = suite.db.First(&deletedEntry, "id = ?", entry.ID.String()).Error
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), gorm.ErrRecordNotFound, err)
}

func (suite *EntryRepositoryTestSuite) TestDeleteNonExistent() {
	// Arrange
	nonExistentID := uuid.New().String()

	// Act
	err := suite.repo.Delete(nonExistentID)

	// Assert
	// GORM не возвращает ошибку при удалении несуществующей записи
	assert.NoError(suite.T(), err)
}

func (suite *EntryRepositoryTestSuite) TestDeleteInvalidID() {
	// Arrange
	invalidID := "invalid-uuid"

	// Act
	err := suite.repo.Delete(invalidID)

	// Assert
	// GORM обработает недопустимый ID без ошибки
	assert.NoError(suite.T(), err)
}

func (suite *EntryRepositoryTestSuite) TestList() {
	// Arrange - создаем несколько тестовых записей
	entries := []*models.Entry{
		{
			ID:        uuid.New(),
			UserID:    uuid.New(),
			Title:     "First Entry",
			Content:   "First Content",
			CreatedAt: time.Now().Add(-2 * time.Hour),
		},
		{
			ID:        uuid.New(),
			UserID:    uuid.New(),
			Title:     "Second Entry",
			Content:   "Second Content",
			CreatedAt: time.Now().Add(-1 * time.Hour),
		},
		{
			ID:        uuid.New(),
			UserID:    uuid.New(),
			Title:     "Third Entry",
			Content:   "Third Content",
			CreatedAt: time.Now(),
		},
	}

	for _, entry := range entries {
		err := suite.db.Create(entry).Error
		suite.Require().NoError(err)
	}

	// Act
	result, err := suite.repo.List()

	// Assert
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 3)

	// Проверяем сортировку по created_at DESC
	assert.Equal(suite.T(), "Third Entry", result[0].Title)  // Самая новая
	assert.Equal(suite.T(), "Second Entry", result[1].Title) // Средняя
	assert.Equal(suite.T(), "First Entry", result[2].Title)  // Самая старая
}

func (suite *EntryRepositoryTestSuite) TestListEmpty() {
	// Act
	result, err := suite.repo.List()

	// Assert
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), result)
}

func (suite *EntryRepositoryTestSuite) TestCreateMultipleEntries() {
	// Arrange
	userID := uuid.New()
	entries := []*models.Entry{
		{
			ID:      uuid.New(),
			UserID:  userID,
			Title:   "Entry 1",
			Content: "Content 1",
		},
		{
			ID:      uuid.New(),
			UserID:  userID,
			Title:   "Entry 2",
			Content: "Content 2",
		},
	}

	// Act
	for _, entry := range entries {
		err := suite.repo.Create(entry)
		assert.NoError(suite.T(), err)
	}

	// Assert
	result, err := suite.repo.List()
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), result, 2)

	// Проверяем, что оба пользователя одинаковые
	assert.Equal(suite.T(), userID, result[0].UserID)
	assert.Equal(suite.T(), userID, result[1].UserID)
}

func TestEntryRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(EntryRepositoryTestSuite))
}

// Дополнительные единичные тесты для специфических случаев

func TestNewEntryRepository(t *testing.T) {
	// Arrange
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Act
	repo := NewEntryRepository(db)

	// Assert
	assert.NotNil(t, repo)
	assert.Implements(t, (*EntryRepository)(nil), repo)
}

func TestEntryRepositoryWithNilDB(t *testing.T) {
	// Act
	repo := NewEntryRepository(nil)

	// Assert
	assert.NotNil(t, repo)

	// Попытка создать запись с nil DB должна вызвать панику
	entry := &models.Entry{
		ID:      uuid.New(),
		UserID:  uuid.New(),
		Title:   "Test",
		Content: "Test",
	}

	assert.Panics(t, func() {
		repo.Create(entry)
	})
}
