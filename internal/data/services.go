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
