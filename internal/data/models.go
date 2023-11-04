package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Categories    CategoryModel
	SubCategories SubCategoryModel
	Services      ServiceModel
	Users         UserModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Categories:    CategoryModel{DB: db},
		SubCategories: SubCategoryModel{DB: db},
		Services:      ServiceModel{DB: db},
		Users:         UserModel{DB: db},
	}
}
