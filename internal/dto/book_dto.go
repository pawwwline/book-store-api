package dto

import (
	"time"

	"github.com/google/uuid"
)

type BookDTO struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Description string    `json:"description"`
	Price       int       `json:"price"`
	ISBN        string    `json:"isbn"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type BookRequest struct {
	Title       string `json:"title"`
	Author      string `json:"author"`
	Description string `json:"description"`
	Price       int    `json:"price"`
	ISBN        string `json:"isbn"`
}
