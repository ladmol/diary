package repos

import (
	"diary/internal/models"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// --- Mock Entry Repository ---

type MockEntryRepository struct {
	mock.Mock
}

func (m *MockEntryRepository) Create(entry *models.Entry) error {
	args := m.Called(entry)
	return args.Error(0)
}

func (m *MockEntryRepository) Read(id string) (*models.Entry, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Entry), args.Error(1)
}

func (m *MockEntryRepository) Update(entry *models.Entry) error {
	args := m.Called(entry)
	return args.Error(0)
}

func (m *MockEntryRepository) Delete(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockEntryRepository) List() ([]*models.Entry, error) {
	args := m.Called()
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Entry), args.Error(1)
}

// --- Repository Test Suite ---

type RepositoryTestSuite struct {
	suite.Suite
	mockEntryRepo *MockEntryRepository
	repo          Repository
}

func (suite *RepositoryTestSuite) SetupTest() {
	suite.mockEntryRepo = new(MockEntryRepository)

	// Создаем репозиторий с мок-зависимостью
	suite.repo = &repository{
		entryRepo: suite.mockEntryRepo,
	}
}

func (suite *RepositoryTestSuite) TearDownTest() {
	suite.mockEntryRepo.AssertExpectations(suite.T())
}

func (suite *RepositoryTestSuite) TestCreate() {
	// Arrange
	entry := &models.Entry{
		ID:      uuid.New(),
		UserID:  uuid.New(),
		Title:   "Test Title",
		Content: "Test Content",
	}

	suite.mockEntryRepo.On("Create", entry).Return(nil)

	// Act
	err := suite.repo.Create(entry)

	// Assert
	assert.NoError(suite.T(), err)
}

func (suite *RepositoryTestSuite) TestCreateError() {
	// Arrange
	entry := &models.Entry{
		ID:      uuid.New(),
		UserID:  uuid.New(),
		Title:   "Test Title",
		Content: "Test Content",
	}

	expectedError := errors.New("database error")
	suite.mockEntryRepo.On("Create", entry).Return(expectedError)

	// Act
	err := suite.repo.Create(entry)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedError, err)
}

func (suite *RepositoryTestSuite) TestRead() {
	// Arrange
	entryID := uuid.New().String()
	expectedEntry := &models.Entry{
		ID:      uuid.MustParse(entryID),
		UserID:  uuid.New(),
		Title:   "Test Title",
		Content: "Test Content",
	}

	suite.mockEntryRepo.On("Read", entryID).Return(expectedEntry, nil)

	// Act
	entry, err := suite.repo.Read(entryID)

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedEntry, entry)
}

func (suite *RepositoryTestSuite) TestReadError() {
	// Arrange
	entryID := uuid.New().String()
	expectedError := gorm.ErrRecordNotFound

	suite.mockEntryRepo.On("Read", entryID).Return(nil, expectedError)

	// Act
	entry, err := suite.repo.Read(entryID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), entry)
	assert.Equal(suite.T(), expectedError, err)
}

func (suite *RepositoryTestSuite) TestUpdate() {
	// Arrange
	entry := &models.Entry{
		ID:      uuid.New(),
		UserID:  uuid.New(),
		Title:   "Updated Title",
		Content: "Updated Content",
	}

	suite.mockEntryRepo.On("Update", entry).Return(nil)

	// Act
	err := suite.repo.Update(entry)

	// Assert
	assert.NoError(suite.T(), err)
}

func (suite *RepositoryTestSuite) TestUpdateError() {
	// Arrange
	entry := &models.Entry{
		ID:      uuid.New(),
		UserID:  uuid.New(),
		Title:   "Updated Title",
		Content: "Updated Content",
	}

	expectedError := errors.New("update failed")
	suite.mockEntryRepo.On("Update", entry).Return(expectedError)

	// Act
	err := suite.repo.Update(entry)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedError, err)
}

func (suite *RepositoryTestSuite) TestDelete() {
	// Arrange
	entryID := uuid.New().String()

	suite.mockEntryRepo.On("Delete", entryID).Return(nil)

	// Act
	err := suite.repo.Delete(entryID)

	// Assert
	assert.NoError(suite.T(), err)
}

func (suite *RepositoryTestSuite) TestDeleteError() {
	// Arrange
	entryID := uuid.New().String()
	expectedError := errors.New("delete failed")

	suite.mockEntryRepo.On("Delete", entryID).Return(expectedError)

	// Act
	err := suite.repo.Delete(entryID)

	// Assert
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), expectedError, err)
}

func (suite *RepositoryTestSuite) TestList() {
	// Arrange
	expectedEntries := []*models.Entry{
		{
			ID:      uuid.New(),
			UserID:  uuid.New(),
			Title:   "Entry 1",
			Content: "Content 1",
		},
		{
			ID:      uuid.New(),
			UserID:  uuid.New(),
			Title:   "Entry 2",
			Content: "Content 2",
		},
	}

	suite.mockEntryRepo.On("List").Return(expectedEntries, nil)

	// Act
	entries, err := suite.repo.List()

	// Assert
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), expectedEntries, entries)
	assert.Len(suite.T(), entries, 2)
}

func (suite *RepositoryTestSuite) TestListError() {
	// Arrange
	expectedError := errors.New("list failed")

	suite.mockEntryRepo.On("List").Return(nil, expectedError)

	// Act
	entries, err := suite.repo.List()

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), entries)
	assert.Equal(suite.T(), expectedError, err)
}

func (suite *RepositoryTestSuite) TestListEmpty() {
	// Arrange
	expectedEntries := []*models.Entry{}

	suite.mockEntryRepo.On("List").Return(expectedEntries, nil)

	// Act
	entries, err := suite.repo.List()

	// Assert
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), entries)
}

func TestRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryTestSuite))
}

// --- Integration Tests ---

type RepositoryIntegrationTestSuite struct {
	suite.Suite
	db   *gorm.DB
	repo Repository
}

func (suite *RepositoryIntegrationTestSuite) SetupSuite() {
	// Создаем in-memory SQLite базу данных для интеграционных тестов
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	suite.Require().NoError(err)

	// Автомиграция
	err = db.AutoMigrate(&models.Entry{})
	suite.Require().NoError(err)

	suite.db = db
	suite.repo = NewRepository(db)
}

func (suite *RepositoryIntegrationTestSuite) TearDownTest() {
	// Очищаем таблицу после каждого теста
	suite.db.Exec("DELETE FROM entries")
}

func (suite *RepositoryIntegrationTestSuite) TestCompleteWorkflow() {
	// Arrange
	entry := &models.Entry{
		ID:      uuid.New(),
		UserID:  uuid.New(),
		Title:   "Integration Test Entry",
		Content: "Integration Test Content",
	}

	// Act & Assert - Create
	err := suite.repo.Create(entry)
	assert.NoError(suite.T(), err)

	// Act & Assert - Read
	readEntry, err := suite.repo.Read(entry.ID.String())
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), entry.Title, readEntry.Title)
	assert.Equal(suite.T(), entry.Content, readEntry.Content)

	// Act & Assert - Update
	readEntry.Title = "Updated Integration Test Entry"
	readEntry.Content = "Updated Integration Test Content"
	err = suite.repo.Update(readEntry)
	assert.NoError(suite.T(), err)

	// Verify update
	updatedEntry, err := suite.repo.Read(entry.ID.String())
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Updated Integration Test Entry", updatedEntry.Title)
	assert.Equal(suite.T(), "Updated Integration Test Content", updatedEntry.Content)

	// Act & Assert - List
	entries, err := suite.repo.List()
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), entries, 1)
	assert.Equal(suite.T(), updatedEntry.ID, entries[0].ID)

	// Act & Assert - Delete
	err = suite.repo.Delete(entry.ID.String())
	assert.NoError(suite.T(), err)

	// Verify deletion
	_, err = suite.repo.Read(entry.ID.String())
	assert.Error(suite.T(), err)
	assert.Equal(suite.T(), gorm.ErrRecordNotFound, err)

	// Verify list is empty
	entries, err = suite.repo.List()
	assert.NoError(suite.T(), err)
	assert.Empty(suite.T(), entries)
}

func (suite *RepositoryIntegrationTestSuite) TestMultipleEntries() {
	// Arrange
	entries := []*models.Entry{
		{
			ID:      uuid.New(),
			UserID:  uuid.New(),
			Title:   "Entry 1",
			Content: "Content 1",
		},
		{
			ID:      uuid.New(),
			UserID:  uuid.New(),
			Title:   "Entry 2",
			Content: "Content 2",
		},
		{
			ID:      uuid.New(),
			UserID:  uuid.New(),
			Title:   "Entry 3",
			Content: "Content 3",
		},
	}

	// Act - Create multiple entries
	for _, entry := range entries {
		err := suite.repo.Create(entry)
		assert.NoError(suite.T(), err)
	}

	// Assert - List all entries
	allEntries, err := suite.repo.List()
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), allEntries, 3)

	// Act & Assert - Read each entry
	for _, originalEntry := range entries {
		readEntry, err := suite.repo.Read(originalEntry.ID.String())
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), originalEntry.Title, readEntry.Title)
		assert.Equal(suite.T(), originalEntry.Content, readEntry.Content)
	}

	// Act & Assert - Delete one entry
	err = suite.repo.Delete(entries[1].ID.String())
	assert.NoError(suite.T(), err)

	// Verify only 2 entries remain
	remainingEntries, err := suite.repo.List()
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), remainingEntries, 2)
}

func TestRepositoryIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(RepositoryIntegrationTestSuite))
}

// --- Constructor Tests ---

func TestNewRepository(t *testing.T) {
	// Arrange
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Act
	repo := NewRepository(db)

	// Assert
	assert.NotNil(t, repo)
	assert.Implements(t, (*Repository)(nil), repo)
	assert.Implements(t, (*EntryRepository)(nil), repo)
}

func TestNewRepositoryWithNilDB(t *testing.T) {
	// Act
	repo := NewRepository(nil)

	// Assert
	assert.NotNil(t, repo)

	// Проверяем, что внутренний entryRepo тоже создается
	// Попытка использовать репозиторий с nil DB должна вызвать панику
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

func TestRepositoryImplementsInterface(t *testing.T) {
	// Arrange
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	// Act
	repo := NewRepository(db)

	// Assert - проверяем, что репозиторий реализует все необходимые интерфейсы
	assert.Implements(t, (*Repository)(nil), repo)
	assert.Implements(t, (*EntryRepository)(nil), repo)

	// Проверяем наличие всех методов
	_, ok := repo.(interface{ Create(*models.Entry) error })
	assert.True(t, ok)

	_, ok = repo.(interface {
		Read(string) (*models.Entry, error)
	})
	assert.True(t, ok)

	_, ok = repo.(interface{ Update(*models.Entry) error })
	assert.True(t, ok)

	_, ok = repo.(interface{ Delete(string) error })
	assert.True(t, ok)

	_, ok = repo.(interface {
		List() ([]*models.Entry, error)
	})
	assert.True(t, ok)
}
