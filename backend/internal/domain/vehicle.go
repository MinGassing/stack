package domain

import "time"

// Structure of a vehicle in the database, converted to go struct
//
// the strings on each field are "struct tags", they tell the json
// encoder what to name each field when this struct is turned into json for
// the frontend
//
// Year is a pointer (*int) so it can be nil, thats how we represent a
// nullable column. "omitempty" means the key is left out of the json entirely
// when its nil
//
// This is placeholder, I am absolutely sure we will need more than this
// for isa's physics calcs, and im sure somethings that i have in here are useless
type Vehicle struct {
	ID              int64     `json:"id"`
	Name            string    `json:"name"`
	Make            string    `json:"make"`
	Model           string    `json:"model"`
	Year            *int      `json:"year,omitempty"`
	MassKg          float64   `json:"mass_kg"`
	DragCoefficient float64   `json:"drag_coefficient"`
	FrontalAreaM2   float64   `json:"frontal_area_m2"`
	CreatedAt       time.Time `json:"created_at"`
}
