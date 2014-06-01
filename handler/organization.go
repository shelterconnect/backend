package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"code.google.com/p/go.crypto/bcrypt"

	"github.com/gorilla/mux"
	"github.com/zachlatta/shelterconnect/database"
	"github.com/zachlatta/shelterconnect/model"
)

func CreateOrganization(w http.ResponseWriter, r *http.Request,
	o *model.Organization) *AppError {
	defer r.Body.Close()
	org, err := model.NewOrganization(r.Body)
	if err != nil {
		return ErrCreatingModel(err)
	}

	err = database.SaveOrganization(org)
	if err != nil {
		if err == model.ErrInvalidOrganizationEmail {
			return ErrCreatingModel(err)
		}
		return ErrDatabase(err)
	}

	return renderJSON(w, org, http.StatusOK)
}

func GetOrganization(w http.ResponseWriter, r *http.Request,
	o *model.Organization) *AppError {
	vars := mux.Vars(r)
	stringID := vars["id"]

	var (
		id      int64
		isEmail bool
	)
	if stringID == "me" {
		if o == nil {
			return ErrUnauthorized()
		}

		id = o.ID
	} else if model.RegexpEmail.MatchString(stringID) {
		isEmail = true
	} else {
		var err error
		id, err = strconv.ParseInt(stringID, 10, 64)
		if err != nil {
			return ErrInvalidID(err)
		}
	}

	var (
		org *model.Organization
		err error
	)
	if isEmail {
		org, err = database.GetOrganizationByEmail(stringID)
	} else {
		org, err = database.GetOrganizationByID(id)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound(err)
		}
		return ErrDatabase(err)
	}

	return renderJSON(w, org, http.StatusOK)
}

func GetAllOrganizations(w http.ResponseWriter, r *http.Request,
	_ *model.Organization) *AppError {
	orgs, err := database.GetAllOrganizations()
	if err != nil {
		return ErrDatabase(err)
	}

	return renderJSON(w, orgs, http.StatusOK)
}

func AuthenticateOrganization(w http.ResponseWriter, r *http.Request,
	_ *model.Organization) *AppError {
	var requestOrg *model.RequestOrganization
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&requestOrg)
	if err != nil {
		return ErrUnmarshalling(err)
	}

	orgFromDB, err := database.GetOrganizationByEmail(requestOrg.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound(err)
		}
		return ErrDatabase(err)
	}

	err = orgFromDB.ComparePassword(requestOrg.Password)
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return &AppError{err, "invalid password", http.StatusBadRequest}
	} else if err != nil {
		return &AppError{err, "error checking password",
			http.StatusInternalServerError}
	}

	token, err := model.NewToken(orgFromDB)
	if err != nil {
		return &AppError{err, "problem creating jwt token",
			http.StatusInternalServerError}
	}

	return renderJSON(w, token, http.StatusOK)
}
