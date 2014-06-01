package handler

import (
	"errors"
	"net/http"
)

func ErrCreatingModel(err error) *AppError {
	return &AppError{err, err.Error(), http.StatusBadRequest}
}

func ErrDatabase(err error) *AppError {
	return &AppError{err, "database error", http.StatusInternalServerError}
}

func ErrInvalidID(err error) *AppError {
	return &AppError{err, "invalid id", http.StatusBadRequest}
}

func ErrNotFound(err error) *AppError {
	return &AppError{err, "not found", http.StatusNotFound}
}

func ErrUnauthorized() *AppError {
	err := errors.New("unauthorized")
	return &AppError{err, err.Error(), http.StatusUnauthorized}
}

func ErrUnmarshalling(err error) *AppError {
	return &AppError{err, err.Error(), http.StatusBadRequest}
}
