package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"cosmetcab.dp.ua/internal/validator"
)

type Service struct {
	ID            int64         `json:"id"`
	Time          sql.NullInt16 `json:"time"`
	Description   string        `json:"description"`
	Price         int           `json:"price"`
	CategoryID    int64         `json:"category_id"`
	SubCategoryID int64         `json:"subcategory_id"`
}

type ServiceWithSubcategory struct {
	ID          int64         `json:"id"`
	Time        sql.NullInt16 `json:"time"`
	Description string        `json:"description"`
	Price       int           `json:"price"`
	Subcategory SubCategory   `json:"subcategory"`
}

type ServiceModel struct {
	DB *sql.DB
}

func ValidateService(service *Service, v *validator.Validator) {
	v.Check(service.Description != "", "description", "must be provided")
	v.Check(service.Price > 0, "price", "must be greater than zero")
	v.Check(service.Time.Int16 >= 0, "time", "must be greater or equal zero")

}
func (m ServiceModel) Insert(service *Service) error {
	query := `
	INSERT INTO services (time, description, price, category_id, subcategory_id)
	VALUES ($1, $2, $3, $4, $5)
	RETURNING id
	`
	args := []any{service.Time, service.Description, service.Price, service.CategoryID, service.SubCategoryID}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&service.ID)
}

func (m ServiceModel) Get(id int64) (*Service, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
	SELECT id, time, description, price, category_id, subcategory_id
	FROM services
	WHERE id=$1;
	`
	var service Service
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&service.ID,
		&service.Time,
		&service.Description,
		&service.Price,
		&service.CategoryID,
		&service.SubCategoryID,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &service, nil

}

func (m ServiceModel) GetAll() ([]*Service, error) {
	query := `
	SELECT id, time, description, price, category_id, subcategory_id
	FROM services
	ORDER BY id`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	services := []*Service{}
	for rows.Next() {
		var service Service
		err = rows.Scan(
			&service.ID,
			&service.Time,
			&service.Description,
			&service.Price,
			&service.CategoryID,
			&service.SubCategoryID,
		)
		if err != nil {
			return nil, err
		}
		services = append(services, &service)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return services, nil
}

func (m ServiceModel) GetAllServicesWithSubcategoriesByID(category_id int64) ([]*ServiceWithSubcategory, error) {
	if category_id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
	SELECT
		s.id AS service_id,
		s.time,
		s.description,
		s.price,
		sc.id AS subcategory_id,
		sc.name
	FROM 
		services s
	LEFT JOIN
		subcategories sc ON s.subcategory_id = sc.id
	WHERE s.category_id = $1;
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query, category_id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	servicesWithSubcategories := []*ServiceWithSubcategory{}
	for rows.Next() {
		var serviceWithSubcategory ServiceWithSubcategory
		err = rows.Scan(
			&serviceWithSubcategory.ID,
			&serviceWithSubcategory.Time,
			&serviceWithSubcategory.Description,
			&serviceWithSubcategory.Price,
			&serviceWithSubcategory.Subcategory.ID,
			&serviceWithSubcategory.Subcategory.Name,
		)
		if err != nil {
			return nil, err
		}
		servicesWithSubcategories = append(servicesWithSubcategories, &serviceWithSubcategory)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return servicesWithSubcategories, nil
}

func (m ServiceModel) Update(service *Service) error {
	query := `
	UPDATE services
	SET time=$1, description=$2, price=$3, category_id=$4, subcategory_id=$5
	WHERE id=$6`
	args := []any{
		service.Time,
		service.Description,
		service.Price,
		service.CategoryID,
		service.SubCategoryID,
		service.ID,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, args...)
	return err
}

func (m ServiceModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
	DELETE FROM services
	WHERE id=$1`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}
