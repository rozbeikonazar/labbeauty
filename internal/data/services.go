package data

import "cosmetcab.dp.ua/internal/validator"

type Service struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	PhotoURL    string `json:"photo_url"`
}

func ValidateService(service *Service, v *validator.Validator) {
	v.Check(service.Title != "", "title", "title must be provided")
	v.Check(len([]rune(service.Title)) <= 55, "title", "must not be more than 55 chars")
	v.Check(service.Description != "", "description", "description must be provided")
	v.Check(len([]rune(service.Description)) >= 20, "description", "description must have more than 20 chars")
	v.Check(service.URL != "", "url", "url must be provided")
	v.Check(service.PhotoURL != "", "photo_url", "photo must be provided")

}
