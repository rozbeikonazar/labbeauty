package data

type Service struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	PhotoURL    string `json:"photo_url"`
}
