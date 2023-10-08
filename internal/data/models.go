package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Categories    CategoryModel
	SubCategories SubCategoryModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Categories:    CategoryModel{DB: db},
		SubCategories: SubCategoryModel{DB: db},
	}
}
