package database

import (
	"time"

	"github.com/zachlatta/shelterconnect/model"
)

const orgCreateStmt = `INSERT INTO organizations (created, updated, name,
email, type, address, location, password) VALUES ($1, $2, $3, $4, $5, $6,
point($7, $8), $9) RETURNING id`

const orgGetAllStmt = `SELECT id, created, updated, name, email, type,
address, location, password FROM organizations`

func SaveOrganization(o *model.Organization) error {
	if o.ID == 0 {
		// TODO: Check email for uniqueness

		o.Created = time.Now()
	}
	o.Updated = time.Now()

	row := db.QueryRow(orgCreateStmt, &o.Created, &o.Updated, &o.Name,
		&o.Email, &o.Type, &o.Address, &o.Location.Latitude,
		&o.Location.Longitude, &o.Password)

	err := row.Scan(&o.ID)
	if err != nil {
		return err
	}

	return nil
}

func GetAllOrganizations() ([]*model.Organization, error) {
	orgs := []*model.Organization{}

	rows, err := db.Query(orgGetAllStmt)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		o := &model.Organization{}
		if err := rows.Scan(&o.ID, &o.Created, &o.Updated, &o.Name, &o.Email,
			&o.Type, &o.Address, &o.Location, &o.Password); err != nil {
			return nil, err
		}

		orgs = append(orgs, o)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return orgs, nil
}
