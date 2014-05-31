package handler

import (
	"net/http"

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
