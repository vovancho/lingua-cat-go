package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/vovancho/lingua-cat-go/dictionary/domain"
)

// MockDictionaryRepository реализует интерфейс domain.DictionaryRepository для тестирования
type MockDictionaryRepository struct {
	mock.Mock
}

func (m *MockDictionaryRepository) GetByID(ctx context.Context, id domain.DictionaryID) (*domain.Dictionary, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*domain.Dictionary), args.Error(1)
}

func (m *MockDictionaryRepository) GetRandomDictionaries(ctx context.Context, lang domain.DictionaryLang, limit uint8) ([]domain.Dictionary, error) {
	args := m.Called(ctx, lang, limit)
	return args.Get(0).([]domain.Dictionary), args.Error(1)
}

func (m *MockDictionaryRepository) IsExistsByNameAndLang(ctx context.Context, name string, lang domain.DictionaryLang) (bool, error) {
	args := m.Called(ctx, name, lang)
	return args.Bool(0), args.Error(1)
}

func (m *MockDictionaryRepository) Store(ctx context.Context, d *domain.Dictionary) error {
	args := m.Called(ctx, d)
	return args.Error(0)
}

func (m *MockDictionaryRepository) ChangeName(ctx context.Context, id domain.DictionaryID, name string) error {
	args := m.Called(ctx, id, name)
	return args.Error(0)
}

func (m *MockDictionaryRepository) Delete(ctx context.Context, id domain.DictionaryID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestDictionaryUseCase_GetByID(t *testing.T) {
	ctx := context.Background()
	dictID := domain.DictionaryID(1)
	expectedDict := &domain.Dictionary{
		ID:   dictID,
		Lang: domain.EnDictionary,
		Name: "test",
		Type: domain.SimpleDictionary,
	}

	// Подтест Success
	t.Run("Success", func(t *testing.T) {
		uc, mockRepo := newUseCaseWithMock(5 * time.Second)

		mockRepo.On("GetByID", mock.Anything, dictID).Return(expectedDict, nil).Once()

		result, err := uc.GetByID(ctx, dictID)
		assert.NoError(t, err)
		assert.Equal(t, expectedDict, result)
		mockRepo.AssertExpectations(t)
	})

	// Подтест NotFound
	t.Run("NotFound", func(t *testing.T) {
		uc, mockRepo := newUseCaseWithMock(5 * time.Second)

		mockRepo.On("GetByID", mock.Anything, dictID).Return((*domain.Dictionary)(nil), errors.New("not found")).Once()

		result, err := uc.GetByID(ctx, dictID)
		assert.ErrorIs(t, err, domain.DictNotFoundError)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	// Подтест Timeout
	t.Run("Timeout", func(t *testing.T) {
		uc, mockRepo := newUseCaseWithMock(5 * time.Second)

		// Создаём контекст с уже истёкшим дедлайном
		ctx, cancel := context.WithTimeout(context.Background(), 0)
		defer cancel()

		mockRepo.On("GetByID", mock.Anything, dictID).Return(expectedDict, nil).Once()

		// Вызываем GetByID с истёкшим контекстом
		result, err := uc.GetByID(ctx, dictID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "context deadline exceeded")
		assert.Nil(t, result)
		mockRepo.AssertNotCalled(t, "GetByID")
	})

	t.Run("InternalTimeout", func(t *testing.T) {
		uc, mockRepo := newUseCaseWithMock(5 * time.Second)

		ctx := context.Background()
		mockRepo.On("GetByID", mock.Anything, dictID).Return((*domain.Dictionary)(nil), context.DeadlineExceeded).Once()

		result, err := uc.GetByID(ctx, dictID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "context deadline exceeded")
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestDictionaryUseCase_GetRandomDictionaries(t *testing.T) {
	ctx := context.Background()
	lang := domain.EnDictionary
	limit := uint8(4)
	dicts := []domain.Dictionary{
		{ID: 1, Lang: lang, Name: "test1", Type: domain.SimpleDictionary},
		{ID: 2, Lang: lang, Name: "test2", Type: domain.SimpleDictionary},
		{ID: 3, Lang: lang, Name: "test3", Type: domain.SimpleDictionary},
		{ID: 4, Lang: lang, Name: "test4", Type: domain.SimpleDictionary},
	}

	// Подтест Success
	t.Run("Success", func(t *testing.T) {
		uc, mockRepo := newUseCaseWithMock(5 * time.Second)

		mockRepo.On("GetRandomDictionaries", mock.Anything, lang, limit).Return(dicts, nil).Once()

		result, err := uc.GetRandomDictionaries(ctx, lang, limit)
		assert.NoError(t, err)
		assert.Equal(t, dicts, result)
		mockRepo.AssertExpectations(t)
	})

	// Подтест InvalidLimit_LessThan4
	t.Run("InvalidLimit_LessThan4", func(t *testing.T) {
		uc, mockRepo := newUseCaseWithMock(5 * time.Second)

		result, err := uc.GetRandomDictionaries(ctx, lang, 3)
		assert.ErrorIs(t, err, DictsRandomCountError)
		assert.Nil(t, result)
		mockRepo.AssertNotCalled(t, "GetRandomDictionaries")
	})

	// Подтест InvalidLimit_MoreThan8
	t.Run("InvalidLimit_MoreThan8", func(t *testing.T) {
		uc, mockRepo := newUseCaseWithMock(5 * time.Second)

		result, err := uc.GetRandomDictionaries(ctx, lang, 9)
		assert.ErrorIs(t, err, DictsRandomCountError)
		assert.Nil(t, result)
		mockRepo.AssertNotCalled(t, "GetRandomDictionaries")
	})

	// Подтест InvalidLang
	t.Run("InvalidLang", func(t *testing.T) {
		uc, mockRepo := newUseCaseWithMock(5 * time.Second)

		result, err := uc.GetRandomDictionaries(ctx, "invalid", limit)
		assert.ErrorIs(t, err, domain.DictLangInvalidError)
		assert.Nil(t, result)
		mockRepo.AssertNotCalled(t, "GetRandomDictionaries")
	})

	// Подтест RepositoryError
	t.Run("RepositoryError", func(t *testing.T) {
		uc, mockRepo := newUseCaseWithMock(5 * time.Second)

		mockRepo.On("GetRandomDictionaries", mock.Anything, lang, limit).Return([]domain.Dictionary{}, errors.New("repo error")).Once()

		result, err := uc.GetRandomDictionaries(ctx, lang, limit)
		assert.ErrorIs(t, err, domain.DictsNotFoundError)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})

	// Подтест NotEnoughDictionaries
	t.Run("NotEnoughDictionaries", func(t *testing.T) {
		uc, mockRepo := newUseCaseWithMock(5 * time.Second)

		insufficientDictionaries := dicts[:2] // Только 2 словаря вместо 4
		mockRepo.On("GetRandomDictionaries", mock.Anything, lang, limit).Return(insufficientDictionaries, nil).Once()

		result, err := uc.GetRandomDictionaries(ctx, lang, limit)
		assert.ErrorIs(t, err, domain.DictsNotFoundError)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}

func TestDictionaryUseCase_Store(t *testing.T) {
	ctx := context.Background()
	dict := &domain.Dictionary{
		Lang: domain.EnDictionary,
		Name: "test",
		Type: domain.SimpleDictionary,
		Translations: []domain.Translation{
			{
				Dictionary: domain.Dictionary{
					Lang: domain.RuDictionary,
					Name: "тест",
					Type: domain.SimpleDictionary,
				},
			},
		},
		Sentences: []domain.Sentence{
			{TextRU: "Это пример предложения.", TextEN: "This is an example sentence."},
		},
	}

	// Подтест Success
	t.Run("Success", func(t *testing.T) {
		uc, mockRepo := newUseCaseWithMock(5 * time.Second)

		mockRepo.On("IsExistsByNameAndLang", mock.Anything, "test", domain.EnDictionary).Return(false, nil).Once()
		mockRepo.On("IsExistsByNameAndLang", mock.Anything, "тест", domain.RuDictionary).Return(false, nil).Once()
		mockRepo.On("Store", mock.Anything, dict).Return(nil).Once()

		err := uc.Store(ctx, dict)
		assert.NoError(t, err)
		assert.Equal(t, "test", dict.Name) // Проверяем, что имя приведено к нижнему регистру
		mockRepo.AssertExpectations(t)
	})

	// Подтест NoTranslations
	t.Run("NoTranslations", func(t *testing.T) {
		uc, mockRepo := newUseCaseWithMock(5 * time.Second)

		invalidDict := &domain.Dictionary{
			Lang: domain.EnDictionary,
			Name: "test",
			Type: domain.SimpleDictionary,
		}

		err := uc.Store(ctx, invalidDict)
		assert.ErrorIs(t, err, domain.DictTranslationRequiredError)
		mockRepo.AssertNotCalled(t, "IsExistsByNameAndLang")
		mockRepo.AssertNotCalled(t, "Store")
	})

	// Подтест InvalidTranslationLang
	t.Run("InvalidTranslationLang", func(t *testing.T) {
		uc, mockRepo := newUseCaseWithMock(5 * time.Second)

		invalidDict := &domain.Dictionary{
			Lang: domain.EnDictionary,
			Name: "test",
			Type: domain.SimpleDictionary,
			Translations: []domain.Translation{
				{
					Dictionary: domain.Dictionary{
						Lang: domain.EnDictionary, // Тот же язык
						Name: "test2",
						Type: domain.SimpleDictionary,
					},
				},
			},
		}

		err := uc.Store(ctx, invalidDict)
		assert.ErrorIs(t, err, domain.DictTranslationLangInvalidError)
		mockRepo.AssertNotCalled(t, "IsExistsByNameAndLang")
		mockRepo.AssertNotCalled(t, "Store")
	})

	// Подтест InvalidStruct
	t.Run("InvalidStruct", func(t *testing.T) {
		uc, mockRepo := newUseCaseWithMock(5 * time.Second)

		invalidDict := &domain.Dictionary{
			Lang: domain.EnDictionary,
			Name: "t", // Слишком короткое имя
			Type: domain.SimpleDictionary,
			Translations: []domain.Translation{
				{
					Dictionary: domain.Dictionary{
						Lang: domain.RuDictionary,
						Name: "тест",
						Type: domain.SimpleDictionary,
					},
				},
			},
		}

		err := uc.Store(ctx, invalidDict)
		assert.Error(t, err)
		assert.IsType(t, validator.ValidationErrors{}, err)
		mockRepo.AssertNotCalled(t, "IsExistsByNameAndLang")
		mockRepo.AssertNotCalled(t, "Store")
	})

	// Подтест DictionaryExists
	t.Run("DictionaryExists", func(t *testing.T) {
		uc, mockRepo := newUseCaseWithMock(5 * time.Second)

		mockRepo.On("IsExistsByNameAndLang", mock.Anything, "test", domain.EnDictionary).Return(true, nil).Once()

		err := uc.Store(ctx, dict)
		assert.ErrorIs(t, err, domain.DictExistsError)
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "Store")
	})

	// Подтест TranslationExists
	t.Run("TranslationExists", func(t *testing.T) {
		uc, mockRepo := newUseCaseWithMock(5 * time.Second)

		mockRepo.On("IsExistsByNameAndLang", mock.Anything, "test", domain.EnDictionary).Return(false, nil).Once()
		mockRepo.On("IsExistsByNameAndLang", mock.Anything, "тест", domain.RuDictionary).Return(true, nil).Once()

		err := uc.Store(ctx, dict)
		assert.ErrorIs(t, err, domain.DictExistsError)
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "Store")
	})

	// Подтест RepositoryError
	t.Run("RepositoryError", func(t *testing.T) {
		uc, mockRepo := newUseCaseWithMock(5 * time.Second)

		mockRepo.On("IsExistsByNameAndLang", mock.Anything, "test", domain.EnDictionary).Return(false, errors.New("repo error")).Once()

		err := uc.Store(ctx, dict)
		assert.Error(t, err)
		assert.EqualError(t, err, "repo error")
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "Store")
	})

	// Подтест MissingLang
	t.Run("MissingLang", func(t *testing.T) {
		uc, mockRepo := newUseCaseWithMock(5 * time.Second)

		invalidDict := &domain.Dictionary{
			Name: "test",
			Type: domain.SimpleDictionary,
			Translations: []domain.Translation{
				{
					Dictionary: domain.Dictionary{
						Lang: domain.RuDictionary,
						Name: "тест",
						Type: domain.SimpleDictionary,
					},
				},
			},
		}

		err := uc.Store(ctx, invalidDict)
		assert.Error(t, err)
		assert.IsType(t, validator.ValidationErrors{}, err)
		validationErrors := err.(validator.ValidationErrors)
		assert.Equal(t, "Lang", validationErrors[0].Field())
		assert.Equal(t, "required", validationErrors[0].Tag())
		mockRepo.AssertNotCalled(t, "IsExistsByNameAndLang")
		mockRepo.AssertNotCalled(t, "Store")
	})

	// Подтест MissingName
	t.Run("MissingName", func(t *testing.T) {
		uc, mockRepo := newUseCaseWithMock(5 * time.Second)

		invalidDict := &domain.Dictionary{
			Lang: domain.EnDictionary,
			Type: domain.SimpleDictionary,
			Translations: []domain.Translation{
				{
					Dictionary: domain.Dictionary{
						Lang: domain.RuDictionary,
						Name: "тест",
						Type: domain.SimpleDictionary,
					},
				},
			},
		}

		err := uc.Store(ctx, invalidDict)
		assert.Error(t, err)
		assert.IsType(t, validator.ValidationErrors{}, err)
		validationErrors := err.(validator.ValidationErrors)
		assert.Equal(t, "Name", validationErrors[0].Field())
		assert.Equal(t, "required", validationErrors[0].Tag())
		mockRepo.AssertNotCalled(t, "IsExistsByNameAndLang")
		mockRepo.AssertNotCalled(t, "Store")
	})

	// Подтест MissingType
	t.Run("MissingType", func(t *testing.T) {
		uc, mockRepo := newUseCaseWithMock(5 * time.Second)

		invalidDict := &domain.Dictionary{
			Lang: domain.EnDictionary,
			Name: "test",
			Translations: []domain.Translation{
				{
					Dictionary: domain.Dictionary{
						Lang: domain.RuDictionary,
						Name: "тест",
						Type: domain.SimpleDictionary,
					},
				},
			},
		}

		err := uc.Store(ctx, invalidDict)
		assert.Error(t, err)
		assert.IsType(t, validator.ValidationErrors{}, err)
		validationErrors := err.(validator.ValidationErrors)
		assert.Equal(t, "Type", validationErrors[0].Field())
		assert.Equal(t, "required", validationErrors[0].Tag())
		mockRepo.AssertNotCalled(t, "IsExistsByNameAndLang")
		mockRepo.AssertNotCalled(t, "Store")
	})

	// Подтест InvalidLang
	t.Run("InvalidLang", func(t *testing.T) {
		uc, mockRepo := newUseCaseWithMock(5 * time.Second)

		invalidDict := &domain.Dictionary{
			Lang: "invalid",
			Name: "test",
			Type: domain.SimpleDictionary,
			Translations: []domain.Translation{
				{
					Dictionary: domain.Dictionary{
						Lang: domain.RuDictionary,
						Name: "тест",
						Type: domain.SimpleDictionary,
					},
				},
			},
		}

		err := uc.Store(ctx, invalidDict)
		assert.Error(t, err)
		assert.IsType(t, validator.ValidationErrors{}, err)
		validationErrors := err.(validator.ValidationErrors)
		assert.Equal(t, "Lang", validationErrors[0].Field())
		assert.Equal(t, "valid_dictionary_lang", validationErrors[0].Tag())
		mockRepo.AssertNotCalled(t, "IsExistsByNameAndLang")
		mockRepo.AssertNotCalled(t, "Store")
	})

	// Подтест InvalidType
	t.Run("InvalidType", func(t *testing.T) {
		uc, mockRepo := newUseCaseWithMock(5 * time.Second)

		invalidDict := &domain.Dictionary{
			Lang: domain.EnDictionary,
			Name: "test",
			Type: 999,
			Translations: []domain.Translation{
				{
					Dictionary: domain.Dictionary{
						Lang: domain.RuDictionary,
						Name: "тест",
						Type: domain.SimpleDictionary,
					},
				},
			},
		}

		err := uc.Store(ctx, invalidDict)
		assert.Error(t, err)
		assert.IsType(t, validator.ValidationErrors{}, err)
		validationErrors := err.(validator.ValidationErrors)
		assert.Equal(t, "Type", validationErrors[0].Field())
		assert.Equal(t, "valid_dictionary_type", validationErrors[0].Tag())
		mockRepo.AssertNotCalled(t, "IsExistsByNameAndLang")
		mockRepo.AssertNotCalled(t, "Store")
	})

	// Подтест LowercaseName
	t.Run("LowercaseName", func(t *testing.T) {
		uc, mockRepo := newUseCaseWithMock(5 * time.Second)

		dict := &domain.Dictionary{
			Lang: domain.EnDictionary,
			Name: "TestName",
			Type: domain.SimpleDictionary,
			Translations: []domain.Translation{
				{
					Dictionary: domain.Dictionary{
						Lang: domain.RuDictionary,
						Name: "ТестИмя",
						Type: domain.SimpleDictionary,
					},
				},
			},
		}

		mockRepo.On("IsExistsByNameAndLang", mock.Anything, "testname", domain.EnDictionary).Return(false, nil).Once()
		mockRepo.On("IsExistsByNameAndLang", mock.Anything, "тестимя", domain.RuDictionary).Return(false, nil).Once()
		mockRepo.On("Store", mock.Anything, dict).Return(nil).Once()

		err := uc.Store(ctx, dict)
		assert.NoError(t, err)
		assert.Equal(t, "testname", dict.Name)
		assert.Equal(t, "тестимя", dict.Translations[0].Dictionary.Name)
		mockRepo.AssertExpectations(t)
	})

	// Подтест TrimSpaces
	t.Run("TrimSpaces", func(t *testing.T) {
		uc, mockRepo := newUseCaseWithMock(5 * time.Second)

		dict := &domain.Dictionary{
			Lang: domain.EnDictionary,
			Name: "  test name  ",
			Type: domain.SimpleDictionary,
			Translations: []domain.Translation{
				{
					Dictionary: domain.Dictionary{
						Lang: domain.RuDictionary,
						Name: "  тест имя  ",
						Type: domain.SimpleDictionary,
					},
				},
			},
		}

		mockRepo.On("IsExistsByNameAndLang", mock.Anything, "test name", domain.EnDictionary).Return(false, nil).Once()
		mockRepo.On("IsExistsByNameAndLang", mock.Anything, "тест имя", domain.RuDictionary).Return(false, nil).Once()
		mockRepo.On("Store", mock.Anything, dict).Return(nil).Once()

		err := uc.Store(ctx, dict)
		assert.NoError(t, err)
		assert.Equal(t, "test name", dict.Name)
		assert.Equal(t, "тест имя", dict.Translations[0].Dictionary.Name)
		mockRepo.AssertExpectations(t)
	})
}

func TestDictionaryUseCase_ChangeName(t *testing.T) {
	ctx := context.Background()
	dictID := domain.DictionaryID(1)
	newName := "newname"

	// Подтест Success
	t.Run("Success", func(t *testing.T) {
		uc, mockRepo := newUseCaseWithMock(5 * time.Second)

		mockRepo.On("GetByID", mock.Anything, dictID).
			Return(&domain.Dictionary{
				ID:   dictID,
				Lang: domain.EnDictionary,
				Name: "oldname",
				Type: domain.SimpleDictionary,
			}, nil).
			Once()

		mockRepo.On("IsExistsByNameAndLang", mock.Anything, "newname", domain.EnDictionary).Return(false, nil).Once()
		mockRepo.On("ChangeName", mock.Anything, dictID, newName).Return(nil).Once()

		err := uc.ChangeName(ctx, dictID, newName)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	// Подтест NotFound
	t.Run("NotFound", func(t *testing.T) {
		uc, mockRepo := newUseCaseWithMock(5 * time.Second)

		mockRepo.On("GetByID", mock.Anything, dictID).Return((*domain.Dictionary)(nil), errors.New("not found")).Once()

		err := uc.ChangeName(ctx, dictID, newName)
		assert.ErrorIs(t, err, domain.DictNotFoundError)
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "IsExistsByNameAndLang")
		mockRepo.AssertNotCalled(t, "ChangeName")
	})

	// Подтест SameName
	t.Run("SameName", func(t *testing.T) {
		uc, mockRepo := newUseCaseWithMock(5 * time.Second)

		mockRepo.On("GetByID", mock.Anything, dictID).
			Return(&domain.Dictionary{
				ID:   dictID,
				Lang: domain.EnDictionary,
				Name: "oldname",
				Type: domain.SimpleDictionary,
			}, nil).
			Once()

		err := uc.ChangeName(ctx, dictID, "oldname")
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "IsExistsByNameAndLang")
		mockRepo.AssertNotCalled(t, "ChangeName")
	})

	// Подтест InvalidName
	t.Run("InvalidName", func(t *testing.T) {
		uc, mockRepo := newUseCaseWithMock(5 * time.Second)

		mockRepo.On("GetByID", mock.Anything, dictID).
			Return(&domain.Dictionary{
				ID:   dictID,
				Lang: domain.EnDictionary,
				Name: "oldname",
				Type: domain.SimpleDictionary,
			}, nil).
			Once()

		err := uc.ChangeName(ctx, dictID, "n") // Слишком короткое
		assert.Error(t, err)
		assert.IsType(t, validator.ValidationErrors{}, err)
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "IsExistsByNameAndLang")
		mockRepo.AssertNotCalled(t, "ChangeName")
	})

	// Подтест NameExists
	t.Run("NameExists", func(t *testing.T) {
		uc, mockRepo := newUseCaseWithMock(5 * time.Second)

		mockRepo.On("GetByID", mock.Anything, dictID).
			Return(&domain.Dictionary{
				ID:   dictID,
				Lang: domain.EnDictionary,
				Name: "oldname",
				Type: domain.SimpleDictionary,
			}, nil).
			Once()

		mockRepo.On("IsExistsByNameAndLang", mock.Anything, "newname", domain.EnDictionary).Return(true, nil).Once()

		err := uc.ChangeName(ctx, dictID, newName)
		assert.ErrorIs(t, err, domain.DictExistsError)
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "ChangeName")
	})

	// Подтест RepositoryError
	t.Run("RepositoryError", func(t *testing.T) {
		uc, mockRepo := newUseCaseWithMock(5 * time.Second)

		mockRepo.On("GetByID", mock.Anything, dictID).
			Return(&domain.Dictionary{
				ID:   dictID,
				Lang: domain.EnDictionary,
				Name: "oldname",
				Type: domain.SimpleDictionary,
			}, nil).
			Once()

		mockRepo.On("IsExistsByNameAndLang", mock.Anything, "newname", domain.EnDictionary).Return(false, nil).Once()
		mockRepo.On("ChangeName", mock.Anything, dictID, newName).Return(errors.New("repo error")).Once()

		err := uc.ChangeName(ctx, dictID, newName)
		assert.Error(t, err)
		assert.EqualError(t, err, "repo error")
		mockRepo.AssertExpectations(t)
	})

	// Подтест EmptyName
	t.Run("EmptyName", func(t *testing.T) {
		uc, mockRepo := newUseCaseWithMock(5 * time.Second)

		mockRepo.On("GetByID", mock.Anything, dictID).
			Return(&domain.Dictionary{
				ID:   dictID,
				Lang: domain.EnDictionary,
				Name: "oldname",
				Type: domain.SimpleDictionary,
			}, nil).
			Once()

		err := uc.ChangeName(ctx, dictID, "")
		assert.Error(t, err)
		assert.IsType(t, validator.ValidationErrors{}, err)
		validationErrors := err.(validator.ValidationErrors)
		assert.Equal(t, "Name", validationErrors[0].Field())
		assert.Equal(t, "required", validationErrors[0].Tag())
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "IsExistsByNameAndLang")
		mockRepo.AssertNotCalled(t, "ChangeName")
	})

	// Подтест LowercaseName
	t.Run("LowercaseName", func(t *testing.T) {
		uc, mockRepo := newUseCaseWithMock(5 * time.Second)

		mockRepo.On("GetByID", mock.Anything, dictID).
			Return(&domain.Dictionary{
				ID:   dictID,
				Lang: domain.EnDictionary,
				Name: "oldname",
				Type: domain.SimpleDictionary,
			}, nil).
			Once()

		mockRepo.On("IsExistsByNameAndLang", mock.Anything, "newname", domain.EnDictionary).Return(false, nil).Once()
		mockRepo.On("ChangeName", mock.Anything, dictID, "newname").Return(nil).Once()

		err := uc.ChangeName(ctx, dictID, "NewName")
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	// Подтест TrimSpaces
	t.Run("TrimSpaces", func(t *testing.T) {
		uc, mockRepo := newUseCaseWithMock(5 * time.Second)

		mockRepo.On("GetByID", mock.Anything, dictID).
			Return(&domain.Dictionary{
				ID:   dictID,
				Lang: domain.EnDictionary,
				Name: "oldname",
				Type: domain.SimpleDictionary,
			}, nil).
			Once()

		mockRepo.On("IsExistsByNameAndLang", mock.Anything, "new name", domain.EnDictionary).Return(false, nil).Once()
		mockRepo.On("ChangeName", mock.Anything, dictID, "new name").Return(nil).Once()

		err := uc.ChangeName(ctx, dictID, "  new name  ")
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestDictionaryUseCase_Delete(t *testing.T) {
	ctx := context.Background()
	dictID := domain.DictionaryID(1)

	// Подтест Success
	t.Run("Success", func(t *testing.T) {
		uc, mockRepo := newUseCaseWithMock(5 * time.Second)

		mockRepo.On("GetByID", mock.Anything, dictID).
			Return(&domain.Dictionary{
				ID:   dictID,
				Lang: domain.EnDictionary,
				Name: "test",
				Type: domain.SimpleDictionary,
			}, nil).
			Once()
		mockRepo.On("Delete", mock.Anything, dictID).Return(nil).Once()

		err := uc.Delete(ctx, dictID)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	// Подтест NotFound
	t.Run("NotFound", func(t *testing.T) {
		uc, mockRepo := newUseCaseWithMock(5 * time.Second)

		mockRepo.On("GetByID", mock.Anything, dictID).Return((*domain.Dictionary)(nil), errors.New("not found")).Once()

		err := uc.Delete(ctx, dictID)
		assert.ErrorIs(t, err, domain.DictNotFoundError)
		mockRepo.AssertExpectations(t)
		mockRepo.AssertNotCalled(t, "Delete")
	})

	// Подтест RepositoryError
	t.Run("RepositoryError", func(t *testing.T) {
		uc, mockRepo := newUseCaseWithMock(5 * time.Second)

		mockRepo.On("GetByID", mock.Anything, dictID).
			Return(&domain.Dictionary{
				ID:   dictID,
				Lang: domain.EnDictionary,
				Name: "test",
				Type: domain.SimpleDictionary,
			}, nil).
			Once()
		mockRepo.On("Delete", mock.Anything, dictID).Return(errors.New("repo error")).Once()

		err := uc.Delete(ctx, dictID)
		assert.Error(t, err)
		assert.EqualError(t, err, "repo error")
		mockRepo.AssertExpectations(t)
	})
}

func newValidator() *validator.Validate {
	v := validator.New()
	v.RegisterValidation("valid_dictionary_lang", func(fl validator.FieldLevel) bool {
		return domain.DictionaryLang(fl.Field().String()).IsValid()
	})
	v.RegisterValidation("valid_dictionary_type", func(fl validator.FieldLevel) bool {
		return domain.DictionaryType(fl.Field().Uint()).IsValid()
	})
	return v
}

func newUseCaseWithMock(timeout time.Duration) (domain.DictionaryUseCase, *MockDictionaryRepository) {
	repo := new(MockDictionaryRepository)
	uc := NewDictionaryUseCase(repo, newValidator(), timeout).(domain.DictionaryUseCase)
	return uc, repo
}
