package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vovancho/lingua-cat-go/exercise/domain"
)

// MockDictionaryRepository implements domain.DictionaryRepository for testing
type MockDictionaryRepository struct {
	mock.Mock
}

func (m *MockDictionaryRepository) GetRandomDictionaries(ctx context.Context, lang domain.DictionaryLang, limit uint8) ([]domain.Dictionary, error) {
	args := m.Called(ctx, lang, limit)
	return args.Get(0).([]domain.Dictionary), args.Error(1)
}

func (m *MockDictionaryRepository) GetDictionariesByIds(ctx context.Context, dictIds []domain.DictionaryID) ([]domain.Dictionary, error) {
	args := m.Called(ctx, dictIds)
	return args.Get(0).([]domain.Dictionary), args.Error(1)
}

func TestDictionaryUseCase_GetRandomDictionaries(t *testing.T) {
	ctx := context.Background()
	lang := domain.EnDictionary
	limit := uint8(4)
	dictionaries := []domain.Dictionary{
		{ID: 1, Lang: domain.EnDictionary, Name: "word1", Type: domain.SimpleDictionary},
		{ID: 2, Lang: domain.EnDictionary, Name: "word2", Type: domain.SimpleDictionary},
		{ID: 3, Lang: domain.EnDictionary, Name: "word3", Type: domain.SimpleDictionary},
		{ID: 4, Lang: domain.EnDictionary, Name: "word4", Type: domain.SimpleDictionary},
	}

	t.Run("Success", func(t *testing.T) {
		uc, mockRepo := newDictionaryUseCaseWithMock()

		mockRepo.On("GetRandomDictionaries", mock.Anything, lang, limit).Return(dictionaries, nil).Once()

		result, err := uc.GetRandomDictionaries(ctx, lang, limit)
		assert.NoError(t, err)
		assert.Equal(t, dictionaries, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("InvalidLimit_Low", func(t *testing.T) {
		uc, mockRepo := newDictionaryUseCaseWithMock()

		result, err := uc.GetRandomDictionaries(ctx, lang, 3)
		assert.ErrorIs(t, err, domain.DictionariesLimitError)
		assert.Nil(t, result)
		mockRepo.AssertNotCalled(t, "GetRandomDictionaries")
	})

	t.Run("InvalidLimit_High", func(t *testing.T) {
		uc, mockRepo := newDictionaryUseCaseWithMock()

		result, err := uc.GetRandomDictionaries(ctx, lang, 9)
		assert.ErrorIs(t, err, domain.DictionariesLimitError)
		assert.Nil(t, result)
		mockRepo.AssertNotCalled(t, "GetRandomDictionaries")
	})

	t.Run("NotFound", func(t *testing.T) {
		uc, mockRepo := newDictionaryUseCaseWithMock()

		mockRepo.On("GetRandomDictionaries", mock.Anything, lang, limit).Return([]domain.Dictionary{}, errors.New("not found")).Once()

		result, err := uc.GetRandomDictionaries(ctx, lang, limit)
		assert.ErrorIs(t, err, domain.DictionariesNotFoundError)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestDictionaryUseCase_GetDictionariesByIds(t *testing.T) {
	ctx := context.Background()
	dictIDs := []domain.DictionaryID{1, 2, 3, 4}
	dictionaries := []domain.Dictionary{
		{ID: 1, Lang: domain.EnDictionary, Name: "word1", Type: domain.SimpleDictionary},
		{ID: 2, Lang: domain.EnDictionary, Name: "word2", Type: domain.SimpleDictionary},
		{ID: 3, Lang: domain.EnDictionary, Name: "word3", Type: domain.SimpleDictionary},
		{ID: 4, Lang: domain.EnDictionary, Name: "word4", Type: domain.SimpleDictionary},
	}

	t.Run("Success", func(t *testing.T) {
		uc, mockRepo := newDictionaryUseCaseWithMock()

		mockRepo.On("GetDictionariesByIds", mock.Anything, dictIDs).Return(dictionaries, nil).Once()

		result, err := uc.GetDictionariesByIds(ctx, dictIDs)
		assert.NoError(t, err)
		assert.Equal(t, dictionaries, result)
		mockRepo.AssertExpectations(t)
	})

	t.Run("NotFound", func(t *testing.T) {
		uc, mockRepo := newDictionaryUseCaseWithMock()

		mockRepo.On("GetDictionariesByIds", mock.Anything, dictIDs).Return([]domain.Dictionary{}, errors.New("not found")).Once()

		result, err := uc.GetDictionariesByIds(ctx, dictIDs)
		assert.ErrorIs(t, err, domain.DictionariesNotFoundError)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func newDictionaryUseCaseWithMock() (domain.DictionaryUseCase, *MockDictionaryRepository) {
	repo := new(MockDictionaryRepository)
	uc := NewDictionaryUseCase(repo)
	return uc, repo
}
