package store

import (
	"context"

	"github.com/mingassing/app/backend/internal/domain"
)

// Function that contains the actual SQL query that creates a vehicle given
// a vehicle object
//
// note this takes a pointer to the vehicle and mutates it, the db generates
// id and created_at, RETURNING hands them back, and Scan copies them into v's
// fields through the pointer.
func (s *Store) CreateVehicle(ctx context.Context, v *domain.Vehicle) error {
	return s.pool.QueryRow(ctx, `
		INSERT INTO vehicles (name, make, model, year, mass_kg, drag_coefficient, frontal_area_m2)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at`,
		v.Name, v.Make, v.Model, v.Year, v.MassKg, v.DragCoefficient, v.FrontalAreaM2,
	).Scan(&v.ID, &v.CreatedAt)
}

// Function that contains the SQL query for returning a single vehicle object given a
// id of that vehicle
func (s *Store) GetVehicle(ctx context.Context, id int64) (*domain.Vehicle, error) {
	v := &domain.Vehicle{}
	err := s.pool.QueryRow(ctx, `
		SELECT id, name, make, model, year, mass_kg, drag_coefficient, frontal_area_m2, created_at
		FROM vehicles WHERE id = $1`, id,
	).Scan(&v.ID, &v.Name, &v.Make, &v.Model, &v.Year, &v.MassKg, &v.DragCoefficient, &v.FrontalAreaM2, &v.CreatedAt)
	if err != nil {
		return nil, err
	}
	return v, nil
}

// Function that queries the database to get the entire list of vehicles from the database
// TODO: Add pagination options + search / filter optimizations
func (s *Store) ListVehicles(ctx context.Context) ([]domain.Vehicle, error) {
	rows, err := s.pool.Query(ctx, `
		SELECT id, name, make, model, year, mass_kg, drag_coefficient, frontal_area_m2, created_at
		FROM vehicles ORDER BY id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	vehicles := []domain.Vehicle{}
	for rows.Next() {
		var v domain.Vehicle
		if err := rows.Scan(&v.ID, &v.Name, &v.Make, &v.Model, &v.Year, &v.MassKg, &v.DragCoefficient, &v.FrontalAreaM2, &v.CreatedAt); err != nil {
			return nil, err
		}
		vehicles = append(vehicles, v)
	}
	return vehicles, rows.Err()
}

// Function (that should never be exposed to a user) which deleted a vehicle from the database
// any func that deletes data from the database should be gated to users unless absolutely necessary
// TODO: Add deleted / archived flag to soft delete
func (s *Store) DeleteVehicle(ctx context.Context, id int64) (bool, error) {
	tag, err := s.pool.Exec(ctx, `DELETE FROM vehicles WHERE id = $1`, id)
	if err != nil {
		return false, err
	}
	return tag.RowsAffected() > 0, nil
}
