package model

import "fmt"

type location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (l *location) Scan(src interface{}) error {
	bytes := src.([]byte)
	str := string(bytes)
	_, err := fmt.Sscanf(str, "(%f,%f)", &l.Latitude, &l.Longitude)
	return err
}
