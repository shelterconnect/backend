package database

import "github.com/zachlatta/shelterconnect/model"

const shelterGetByIDStmt = `SELECT id, created, updated, name, description,
email, type, address, location, password FROM organizations WHERE type = 1 AND
id = $1`

const shelterGetAllStmt = `SELECT id, created, updated, name, description,
email, type, address, location, password FROM organizations WHERE type = 1`

func GetShelterByID(id int64) (*model.Organization, error) {
	s := model.Organization{}
	row := db.QueryRow(shelterGetByIDStmt, id)
	if err := row.Scan(&s.ID, &s.Created, &s.Updated, &s.Name, &s.Description,
		&s.Email, &s.Type, &s.Address, &s.Location, &s.Password); err != nil {
		return nil, err
	}
	return &s, nil
}

func GetAllShelters() ([]*model.Organization, error) {
	shelters := []*model.Organization{}

	rows, err := db.Query(shelterGetAllStmt)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		o := &model.Organization{}
		if err := rows.Scan(&o.ID, &o.Created, &o.Updated, &o.Name,
			&o.Description, &o.Email, &o.Type, &o.Address, &o.Location,
			&o.Password); err != nil {
			return nil, err
		}

		shelters = append(shelters, o)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return shelters, nil
}
