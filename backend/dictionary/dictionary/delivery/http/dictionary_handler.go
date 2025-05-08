package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/vovancho/lingua-cat-go/dictionary/domain"
	_internalError "github.com/vovancho/lingua-cat-go/dictionary/internal/error"
	"github.com/vovancho/lingua-cat-go/dictionary/internal/request"
	"github.com/vovancho/lingua-cat-go/dictionary/internal/response"
	"net/http"
)

type DictionaryStoreRequest struct {
	Lang         domain.DictionaryLang `json:"lang" validate:"required,valid_dictionary_lang"`
	Name         string                `json:"name" validate:"required,min=2"`
	Type         domain.DictionaryType `json:"type" validate:"required,valid_dictionary_type"`
	Translations []struct {
		Lang domain.DictionaryLang `json:"lang" validate:"required,valid_dictionary_lang,valid_dict_translation_lang"`
		Name string                `json:"name" validate:"required,min=2"`
		Type domain.DictionaryType `json:"type" validate:"required,valid_dictionary_type"`
	} `json:"translations" validate:"required,min=1,dive"`
	Sentences []struct {
		TextRU string `json:"text_ru" validate:"required,min=5"`
		TextEN string `json:"text_en" validate:"required,min=5"`
	} `json:"sentences" validate:"dive"`
}

type DictionaryChangeNameRequest struct {
	Name string `json:"name" validate:"required,min=2"`
}

type DictionaryHandler struct {
	DUseCase domain.DictionaryUseCase
	validate *validator.Validate
}

func NewDictionaryHandler(router *http.ServeMux, v *validator.Validate, d domain.DictionaryUseCase) {
	handler := &DictionaryHandler{
		DUseCase: d,
		validate: v,
	}

	router.HandleFunc("GET /v1/dictionary/{id}", request.WithID(handler.GetByID))
	router.HandleFunc("POST /v1/dictionary", handler.Store)
	router.HandleFunc("POST /v1/dictionary/{id}/name", request.WithID(handler.ChangeName))
	router.HandleFunc("DELETE /v1/dictionary/{id}", request.WithID(handler.Delete))
}

// GetByID godoc
// @Summary Получить словарь по ID
// @Description Получает словарь по указанному идентификатору
// @Security BearerAuth
// @Tags dictionary
// @Accept json
// @Produce json
// @Param id path uint64 true "ID словаря"
// @Success 200 {object} response.APIResponse{data=map[string]domain.Dictionary} "Словарь найден"
// @Failure 404 {object} response.APIResponse "Словарь не найден"
// @Router /v1/dictionary/{id} [get]
func (d *DictionaryHandler) GetByID(w http.ResponseWriter, r *http.Request, id uint64) {
	dictID := domain.DictionaryID(id)
	dictionaries, err := d.DUseCase.GetByIDs(r.Context(), []domain.DictionaryID{dictID})
	if err != nil {
		appError := _internalError.NewAppError(http.StatusNotFound, "Словарь не найден", err)
		response.Error(appError, r)
		return
	}
	if len(dictionaries) == 0 {
		appError := _internalError.NewAppError(http.StatusNotFound, "Словарь не найден", domain.DictNotFoundError)
		response.Error(appError, r)
		return
	}
	dictionary := dictionaries[0]

	response.JSON(w, http.StatusOK, response.APIResponse{
		Data: map[string]any{
			"dictionary": dictionary,
		},
	})
}

// Store godoc
// @Summary Создать новый словарь
// @Description Создает новый словарь с предоставленными данными
// @Security BearerAuth
// @Tags dictionary
// @Accept json
// @Produce json
// @Param dictionary body DictionaryStoreRequest true "Данные словаря"
// @Success 201 {object} response.APIResponse{data=map[string]domain.Dictionary} "Словарь создан"
// @Failure 400 {object} response.APIResponse "Некорректный запрос"
// @Router /v1/dictionary [post]
func (d *DictionaryHandler) Store(w http.ResponseWriter, r *http.Request) {
	var requestBody DictionaryStoreRequest
	if err := d.validateRequest(r, &requestBody); err != nil {
		response.Error(err, r)
		return
	}

	dictionary := newDictionaryByRequest(requestBody)

	for _, t := range dictionary.Translations {
		if t.Dictionary.Lang == dictionary.Lang {
			err := _internalError.NewAppError(http.StatusBadRequest, "Ошибка валидации", domain.DictTranslationLangInvalidError)
			response.Error(err, r)
			return
		}
	}

	if err := d.DUseCase.Store(r.Context(), &dictionary); err != nil {
		err = _internalError.NewAppError(http.StatusBadRequest, "Ошибка сохранения словаря", err)
		response.Error(err, r)
		return
	}

	response.JSON(w, http.StatusCreated, response.APIResponse{
		Data: map[string]any{
			"dictionary": dictionary,
		},
	})
}

// ChangeName godoc
// @Summary Изменить имя словаря
// @Description Изменяет имя словаря по указанному идентификатору
// @Security BearerAuth
// @Tags dictionary
// @Accept json
// @Produce json
// @Param id path uint64 true "ID словаря"
// @Param name body DictionaryChangeNameRequest true "Новое имя словаря"
// @Success 200 {object} response.APIResponse{data=map[string]domain.Dictionary} "Имя словаря обновлено"
// @Failure 400 {object} response.APIResponse "Некорректный запрос"
// @Failure 404 {object} response.APIResponse "Словарь не найден"
// @Router /v1/dictionary/{id}/name [post]
func (d *DictionaryHandler) ChangeName(w http.ResponseWriter, r *http.Request, id uint64) {
	var requestBody DictionaryChangeNameRequest
	if err := d.validateRequest(r, &requestBody); err != nil {
		response.Error(err, r)
		return
	}

	if err := d.DUseCase.ChangeName(r.Context(), domain.DictionaryID(id), requestBody.Name); err != nil {
		err = _internalError.NewAppError(http.StatusBadRequest, "Ошибка сохранения словаря", err)
		response.Error(err, r)
		return
	}

	dictID := domain.DictionaryID(id)
	dictionaries, err := d.DUseCase.GetByIDs(r.Context(), []domain.DictionaryID{dictID})
	if err != nil {
		response.Error(err, r)
		return
	}
	if len(dictionaries) == 0 {
		appError := _internalError.NewAppError(http.StatusNotFound, "Словарь не найден", domain.DictNotFoundError)
		response.Error(appError, r)
		return
	}
	dictionary := dictionaries[0]

	response.JSON(w, http.StatusOK, response.APIResponse{
		Data: map[string]any{
			"dictionary": dictionary,
		},
	})
}

// Delete godoc
// @Summary Удалить словарь
// @Description Удаляет словарь по указанному идентификатору
// @Security BearerAuth
// @Tags dictionary
// @Accept json
// @Produce json
// @Param id path uint64 true "ID словаря"
// @Success 204 {object} response.APIResponse "Словарь удален"
// @Failure 404 {object} response.APIResponse "Словарь не найден"
// @Router /v1/dictionary/{id} [delete]
func (d *DictionaryHandler) Delete(w http.ResponseWriter, r *http.Request, id uint64) {
	if err := d.DUseCase.Delete(r.Context(), domain.DictionaryID(id)); err != nil {
		err = _internalError.NewAppError(http.StatusNotFound, "Словарь не найден", err)
		response.Error(err, r)
		return
	}

	response.JSON(w, http.StatusNoContent, response.APIResponse{})
}

func (d *DictionaryHandler) validateRequest(r *http.Request, req any) error {
	if err := request.FromJSON(r, req); err != nil {
		return _internalError.NewAppError(http.StatusBadRequest, "Некорректный синтаксис JSON", _internalError.InvalidDecodeJsonError)
	}

	if err := d.validate.Struct(req); err != nil {
		return err
	}

	return nil
}

func newDictionaryByRequest(req DictionaryStoreRequest) domain.Dictionary {
	dictionary := domain.Dictionary{
		Name: req.Name,
		Type: req.Type,
		Lang: req.Lang,
	}

	for _, s := range req.Sentences {
		sentence := domain.Sentence{
			TextRU: s.TextRU, // Русский перевод
			TextEN: s.TextEN, // Английский текст
		}
		dictionary.Sentences = append(dictionary.Sentences, sentence)
	}

	for _, t := range req.Translations {
		transDict := domain.Dictionary{
			Name: t.Name,
			Type: t.Type,
			Lang: t.Lang,
		}
		translation := domain.Translation{
			Dictionary: transDict,
		}
		dictionary.Translations = append(dictionary.Translations, translation)
	}

	return dictionary
}
