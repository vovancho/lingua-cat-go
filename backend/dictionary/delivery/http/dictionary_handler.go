package http

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/vovancho/lingua-cat-go/dictionary/domain"
	_internalError "github.com/vovancho/lingua-cat-go/pkg/error"
	"github.com/vovancho/lingua-cat-go/pkg/request"
	"github.com/vovancho/lingua-cat-go/pkg/response"
)

type DictionaryStoreRequest struct {
	Lang         domain.DictionaryLang `json:"lang" validate:"required,valid_dictionary_lang"`
	Name         string                `json:"name" validate:"min=2"`
	Type         domain.DictionaryType `json:"type" validate:"required,valid_dictionary_type"`
	Translations []struct {
		Lang domain.DictionaryLang `json:"lang" validate:"required,valid_dictionary_lang,valid_dict_translation_lang"`
		Name string                `json:"name" validate:"min=2"`
		Type domain.DictionaryType `json:"type" validate:"required,valid_dictionary_type"`
	} `json:"translations" validate:"min=1,dive"`
	Sentences []struct {
		TextRU string `json:"text_ru" validate:"min=5"`
		TextEN string `json:"text_en" validate:"min=5"`
	} `json:"sentences" validate:"dive"`
}

type DictionaryChangeNameRequest struct {
	Name string `json:"name" validate:"min=2"`
}

type DictionaryData struct {
	Dictionary domain.Dictionary `json:"dictionary"`
}

type dictionaryHandler struct {
	responder         response.Responder
	dictionaryUseCase domain.DictionaryUseCase
	validator         *validator.Validate
}

func NewDictionaryHandler(
	router *http.ServeMux,
	responder response.Responder,
	dictionaryUseCase domain.DictionaryUseCase,
	validator *validator.Validate,
) {
	handler := &dictionaryHandler{
		responder:         responder,
		dictionaryUseCase: dictionaryUseCase,
		validator:         validator,
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
// @Tags Dictionary
// @Param id path uint64 true "ID словаря"
// @Success 200 {object} response.APIResponse{data=DictionaryData} "Словарь найден"
// @Failure 404 {object} response.APIResponse "Словарь не найден"
// @Router /v1/dictionary/{id} [get]
func (h *dictionaryHandler) GetByID(w http.ResponseWriter, r *http.Request, id uint64) {
	dictID := domain.DictionaryID(id)
	dictionaries, err := h.dictionaryUseCase.GetByIDs(r.Context(), []domain.DictionaryID{dictID})
	if err != nil {
		appError := _internalError.NewAppError(http.StatusNotFound, "Словарь не найден", err)
		h.responder.Error(w, appError)

		return
	}

	if len(dictionaries) == 0 {
		appError := _internalError.NewAppError(http.StatusNotFound, "Словарь не найден", domain.DictNotFoundError)
		h.responder.Error(w, appError)

		return
	}

	dictionary := dictionaries[0]

	h.responder.Success(w, http.StatusOK, map[string]any{
		"dictionary": dictionary,
	})
}

// Store godoc
// @Summary Создать новый словарь
// @Description Создает новый словарь с предоставленными данными
// @Security BearerAuth
// @Tags Dictionary
// @Param dictionary body DictionaryStoreRequest true "Данные словаря"
// @Success 201 {object} response.APIResponse{data=DictionaryData} "Словарь создан"
// @Failure 400 {object} response.APIResponse "Некорректный запрос"
// @Router /v1/dictionary [post]
func (h *dictionaryHandler) Store(w http.ResponseWriter, r *http.Request) {
	var requestBody DictionaryStoreRequest
	if err := h.validateRequest(r, &requestBody); err != nil {
		h.responder.Error(w, err)

		return
	}

	dictionary := newDictionaryByRequest(requestBody)

	for _, t := range dictionary.Translations {
		if t.Dictionary.Lang == dictionary.Lang {
			err := _internalError.NewAppError(http.StatusBadRequest, "Ошибка валидации", domain.DictTranslationLangInvalidError)
			h.responder.Error(w, err)

			return
		}
	}

	if err := h.dictionaryUseCase.Store(r.Context(), &dictionary); err != nil {
		err = _internalError.NewAppError(http.StatusBadRequest, "Ошибка сохранения словаря", err)
		h.responder.Error(w, err)

		return
	}

	h.responder.Success(w, http.StatusCreated, map[string]any{
		"dictionary": dictionary,
	})
}

// ChangeName godoc
// @Summary Изменить имя словаря
// @Description Изменяет имя словаря по указанному идентификатору
// @Security BearerAuth
// @Tags Dictionary
// @Param id path uint64 true "ID словаря"
// @Param name body DictionaryChangeNameRequest true "Новое имя словаря"
// @Success 200 {object} response.APIResponse{data=DictionaryData} "Имя словаря обновлено"
// @Failure 400 {object} response.APIResponse "Некорректный запрос"
// @Failure 404 {object} response.APIResponse "Словарь не найден"
// @Router /v1/dictionary/{id}/name [post]
func (h *dictionaryHandler) ChangeName(w http.ResponseWriter, r *http.Request, id uint64) {
	var requestBody DictionaryChangeNameRequest
	if err := h.validateRequest(r, &requestBody); err != nil {
		h.responder.Error(w, err)

		return
	}

	if err := h.dictionaryUseCase.ChangeName(r.Context(), domain.DictionaryID(id), requestBody.Name); err != nil {
		err = _internalError.NewAppError(http.StatusBadRequest, "Ошибка сохранения словаря", err)
		h.responder.Error(w, err)

		return
	}

	dictID := domain.DictionaryID(id)
	dictionaries, err := h.dictionaryUseCase.GetByIDs(r.Context(), []domain.DictionaryID{dictID})
	if err != nil {
		h.responder.Error(w, err)

		return
	}

	if len(dictionaries) == 0 {
		appError := _internalError.NewAppError(http.StatusNotFound, "Словарь не найден", domain.DictNotFoundError)
		h.responder.Error(w, appError)

		return
	}

	dictionary := dictionaries[0]

	h.responder.Success(w, http.StatusOK, map[string]any{
		"dictionary": dictionary,
	})
}

// Delete godoc
// @Summary Удалить словарь
// @Description Удаляет словарь по указанному идентификатору
// @Security BearerAuth
// @Tags Dictionary
// @Param id path uint64 true "ID словаря"
// @Success 204 {object} response.APIResponse "Словарь удален"
// @Failure 404 {object} response.APIResponse "Словарь не найден"
// @Router /v1/dictionary/{id} [delete]
func (h *dictionaryHandler) Delete(w http.ResponseWriter, r *http.Request, id uint64) {
	if err := h.dictionaryUseCase.Delete(r.Context(), domain.DictionaryID(id)); err != nil {
		err = _internalError.NewAppError(http.StatusNotFound, "Словарь не найден", err)
		h.responder.Error(w, err)

		return
	}

	h.responder.Success(w, http.StatusNoContent, nil)
}

func (h *dictionaryHandler) validateRequest(r *http.Request, req any) error {
	if err := request.FromJSON(r, req); err != nil {
		return _internalError.NewAppError(http.StatusBadRequest, "Некорректный синтаксис JSON", _internalError.InvalidDecodeJsonError)
	}

	if err := h.validator.Struct(req); err != nil {
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
