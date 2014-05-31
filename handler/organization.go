package handler

import (
	"net/http"

	"github.com/zachlatta/shelterconnect/model"
)

func CreateOrganization(w http.ResponseWriter, r *http.Request,
	o *model.Organization) *AppError {
	defer r.Body.Close()
	org, err := model.NewOrganization(r.Body)
	if err != nil {
		return ErrCreatingModel(err)
	}

	return renderJSON(w, org, http.StatusOK)
}
