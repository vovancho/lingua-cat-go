package domain

import "errors"

var (
	ExerciseNotFoundError        = errors.New("EXERCISE_NOT_FOUND")
	ExerciseCompletedError       = errors.New("EXERCISE_COMPLETED")
	TaskNotFoundError            = errors.New("TASK_NOT_FOUND")
	NewTaskNotAllowedError       = errors.New("NEW_TASK_NOT_ALLOWED")
	DictionaryNotFoundError      = errors.New("DICTIONARY_NOT_FOUND")
	DictionariesNotFoundError    = errors.New("DICTIONARIES_NOT_FOUND")
	DictionaryLangIncorrectError = errors.New("DICTIONARIES_LANG_INCORRECT")
	DictionariesLimitError       = errors.New("DICTIONARIES_LIMIT_INVALID")
)
