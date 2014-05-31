package handler

import "net/http"

func ErrCreatingModel(err error) *AppError {
	return &AppError{err, err.Error(), http.StatusBadRequest}
}

func ErrDatabase(err error) *AppError {
	return &AppError{err, "database error", http.StatusInternalServerError}
}
