package models

type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	ImageID     string  `json:"image_id"`
	Price       float64 `json:"price"`
}

// In-memory storage for demo
var Products = make(map[string]Product)
var ProductIDCounter = 1
