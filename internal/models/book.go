package models

import (
	"time"

	"github.com/google/uuid"
)

type Book struct {
	ID          uuid.UUID
	Title       string
	Description string
	Author      string
	ISBN        string
	Price       int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type BookParams struct {
	ID          uuid.UUID
	Title       string
	Description string
	Author      string
	ISBN        string
	Price       int
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewBook(book BookParams) (Book, error) {
	err := validateBook(book)
	if err != nil {
		return Book{}, err
	}

	return Book(book), nil
}
