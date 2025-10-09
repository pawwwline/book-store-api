package models

type Book struct {
	ID          string
	Title       string
	Description string
	Author      string
	ISBN        string
	Price       int
}

type BookParams struct {
	ID          string
	Title       string
	Description string
	Author      string
	ISBN        string
	Price       int
}

func NewBook(book BookParams) (*Book, error) {
	err := validateBook(book)
	if err != nil {
		return nil, err
	}

	return &Book{
		ID:          book.ID,
		Title:       book.Title,
		Description: book.Description,
		Author:      book.Author,
		ISBN:        book.ISBN,
		Price:       book.Price,
	}, nil
}
