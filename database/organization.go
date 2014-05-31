package database

import (
	"time"

	"github.com/zachlatta/shelterconnect/model"
)

const orgCreateStmt = `INSERT INTO organizations (created, updated, name,
email, address, location, password) VALUES ($1, $2, $3, $4, $5, point($6, $7),
$8) RETURNING id`

func SaveOrganization(o *model.Organization) error {
	if o.ID == 0 {
		// TODO: Check email for uniqueness

		o.Created = time.Now()
	}
	o.Updated = time.Now()

	row := db.QueryRow(orgCreateStmt, &o.Created, &o.Updated, &o.Name,
		&o.Email, &o.Address, &o.Location.Latitude, &o.Location.Longitude,
		&o.Password)

	err := row.Scan(&o.ID)
	if err != nil {
		return err
	}

	return nil
}
