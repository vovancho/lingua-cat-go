package domain

import "errors"

var (
	ExerciseCompleteListNotFoundError = errors.New("EXERCISE_COMPLETE_LIST_NOT_FOUND")
	UserNotFoundError                 = errors.New("USER_NOT_FOUND")
)
