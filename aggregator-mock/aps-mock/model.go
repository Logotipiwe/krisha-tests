package aps_mock

type MockApBean struct {
	Id     int64    `json:"id,omitempty"`
	Title  string   `json:"title,omitempty"`
	Price  int64    `json:"price,omitempty"`
	Images []string `json:"images,omitempty"`
}
