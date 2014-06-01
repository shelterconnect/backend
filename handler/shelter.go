package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/zachlatta/shelterconnect/database"
	"github.com/zachlatta/shelterconnect/model"
)

func GetAllShelters(w http.ResponseWriter, r *http.Request,
	_ *model.Organization) *AppError {
	shelters, err := database.GetAllShelters()
	if err != nil {
		return ErrDatabase(err)
	}

	return renderJSON(w, shelters, http.StatusOK)
}

func GetShelter(w http.ResponseWriter, r *http.Request,
	_ *model.Organization) *AppError {
	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		return ErrInvalidID(err)
	}

	shelter, err := database.GetShelterByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound(err)
		}

		return ErrDatabase(err)
	}

	return renderJSON(w, shelter, http.StatusOK)
}
