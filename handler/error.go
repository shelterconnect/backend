package handler

import "net/http"

func ErrCreatingModel(err error) *AppError {
	return &AppError{err, err.Error(), http.StatusBadRequest}
}
