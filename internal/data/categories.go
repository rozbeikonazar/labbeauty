package data

import "cosmetcab.dp.ua/internal/validator"

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
	v.Check(category.PhotoURL != "", "photo_url", "photo must be provided")

}
