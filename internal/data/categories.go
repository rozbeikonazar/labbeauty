package data

import (
	"database/sql"

	"cosmetcab.dp.ua/internal/validator"
)

type Category struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	PhotoURL    string `json:"photo_url"`
}

func ValidateCategory(category *Category, v *validator.Validator) {
	v.Check(category.Title != "", "title", "title must be provided")
	v.Check(len([]rune(category.Title)) <= 55, "title", "must not be more than 55 chars")
	v.Check(category.Description != "", "description", "description must be provided")
	v.Check(len([]rune(category.Description)) >= 20, "description", "description must have more than 20 chars")

}

type CategoryModel struct {
	DB *sql.DB
}

func (c CategoryModel) Insert(category *Category) error {
	query := `
	INSERT INTO categories (title, description, photo_url)
	VALUES ($1, $2, $3)
	RETURNING id`

	args := []any{category.Title, category.Description, category.PhotoURL}

	return c.DB.QueryRow(query, args...).Scan(&category.ID)
}

func (c CategoryModel) Get(id int64) (*Category, error) {
	return nil, nil
}

func (c CategoryModel) Update(category *Category) error {
	return nil

}

func (c CategoryModel) Delete(id int64) error {
	return nil
}
