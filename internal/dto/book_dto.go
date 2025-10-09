package dto

import "time"

type BookDTO struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	ISBN        string    `json:"isbn"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
