package model

import (
	"encoding/json"
	"errors"
	"io"
	"regexp"
	"time"

	"github.com/kellydunn/golang-geo"

	"code.google.com/p/go.crypto/bcrypt"
)

type organizationType int

const (
	OrganizationDefault organizationType = iota
	OrganizationShelter
	OrganizationRestaurant
	OrganizationChurch
)

var geocoder = geo.GoogleGeocoder{}

var (
	ErrInvalidOrganizationName        = errors.New("invalid name")
	ErrInvalidOrganizationDescription = errors.New("invalid description")
	ErrInvalidOrganizationEmail       = errors.New("invalid email")
	ErrInvalidOrganizationType        = errors.New("invalid type")
	ErrInvalidOrganizationAddress     = errors.New("invalid address")
	ErrInvalidOrganizationPassword    = errors.New("invalid password")
)

var RegexpEmail = regexp.MustCompile(`^[^@]+@[^@.]+\.[^@.]+`)

type Organization struct {
	ID          int64            `json:"id"`
	Created     time.Time        `json:"created"`
	Updated     time.Time        `json:"updated"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Email       string           `json:"email"`
	Type        organizationType `json:"type"`
	Address     string           `json:"address"`
	Location    location         `json:"location"`
	Password    string           `json:"-"`
}

type RequestOrganization struct {
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Email       string           `json:"email"`
	Type        organizationType `json:"type"`
	Address     string           `json:"address"`
	Password    string           `json:"password"`
}

func NewOrganization(jsonReader io.Reader) (*Organization, error) {
	var rO RequestOrganization
	err := json.NewDecoder(jsonReader).Decode(&rO)
	if err != nil {
		return nil, err
	}

	err = rO.validate()
	if err != nil {
		return nil, err
	}

	b, err := bcrypt.GenerateFromPassword([]byte(rO.Password),
		bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	point, err := geocoder.Geocode(rO.Address)
	if err != nil {
		return nil, ErrInvalidOrganizationAddress
	}

	org := Organization{
		Name:    rO.Name,
		Email:   rO.Email,
		Type:    rO.Type,
		Address: rO.Address,
		Location: location{
			Latitude:  point.Lat(),
			Longitude: point.Lng(),
		},
		Password: string(b),
	}

	return &org, nil
}

func (o *RequestOrganization) validate() error {
	switch {
	case len(o.Name) == 0 || len(o.Name) > 255:
		return ErrInvalidOrganizationName
	case len(o.Description) == 0 || len(o.Description) > 4096:
		return ErrInvalidOrganizationName
	case RegexpEmail.MatchString(o.Email) == false:
		return ErrInvalidOrganizationEmail
	case o.Type < OrganizationShelter || o.Type > OrganizationChurch:
		return ErrInvalidOrganizationType
	case len(o.Address) == 0 || len(o.Address) > 255:
		return ErrInvalidOrganizationAddress
	case len(o.Password) < 6 || len(o.Password) > 255:
		return ErrInvalidOrganizationPassword
	default:
		return nil
	}
}

// ComparePassword compares the supplied password to the user password stored
// in the database.
func (o *Organization) ComparePassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(o.Password), []byte(password))
}
