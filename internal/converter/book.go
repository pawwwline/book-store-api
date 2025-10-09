package converter

import (
	"book-store-api/internal/dto"
	"book-store-api/internal/models"
)

func ToBookResponse(b models.Book) dto.BookDTO {
	return dto.BookDTO{
		ID:          b.ID.String(),
		Title:       b.Title,
		Description: b.Description,
		Author:      b.Author,
		ISBN:        b.ISBN,
		Price:       b.Price,
	}
}

func ToBookResponseList(books []models.Book) []dto.BookDTO {
	resp := make([]dto.BookDTO, 0, len(books))
	for _, b := range books {
		resp = append(resp, ToBookResponse(b))
	}
	return resp
}

func ToBookParams(book dto.BookDTO) models.BookParams {
	return models.BookParams{
		Title:       book.Title,
		Description: book.Description,
		ISBN:        book.ISBN,
		Price:       book.Price,
		Author:      book.Author,
	}
}
