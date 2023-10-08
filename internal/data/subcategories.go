package data

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"cosmetcab.dp.ua/internal/validator"
)

type SubCategory struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type SubCategoryModel struct {
	DB *sql.DB
}

func ValidateSubCategory(subCategory *SubCategory, v *validator.Validator) {
	v.Check(subCategory.Name != "", "name", "must be provided")
	v.Check(len([]rune(subCategory.Name)) >= 8, "name", "name must have more than 8 chars")
}

func (m SubCategoryModel) Insert(subCategory *SubCategory) error {
	query := `
	INSERT INTO subcategories (name)
	VALUES($1)
	RETURNING id`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	return m.DB.QueryRowContext(ctx, query, subCategory.Name).Scan(&subCategory.ID)

}

func (m SubCategoryModel) Get(id int64) (*SubCategory, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}
	query := `
	SELECT id, name
	FROM subcategories
	WHERE id=$1`
	var subCategory SubCategory
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := m.DB.QueryRowContext(ctx, query, id).Scan(&subCategory.ID, &subCategory.Name)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}

	}
	return &subCategory, nil
}

func (m SubCategoryModel) GetAll() ([]*SubCategory, error) {
	query := `
	SELECT id, name
	FROM subcategories`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	subCategories := []*SubCategory{}
	for rows.Next() {
		var subCategory SubCategory
		err := rows.Scan(
			&subCategory.ID,
			&subCategory.Name,
		)
		if err != nil {
			return nil, err
		}
		subCategories = append(subCategories, &subCategory)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return subCategories, nil

}

func (m SubCategoryModel) Update(subCategory *SubCategory) error {
	query := `
	UPDATE subcategories
	SET name=$1
	WHERE id=$2`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := m.DB.ExecContext(ctx, query, subCategory.Name, subCategory.ID)
	return err

}

func (m SubCategoryModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	query := `
	DELETE FROM subcategories
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
