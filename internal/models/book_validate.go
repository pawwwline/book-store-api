package models

import "fmt"

func validateBook(book BookParams) error {
	if book.ID == "" {
		return fmt.Errorf("%w: book id is required", ErrDomainValidation)
	}
	if book.Title == "" {
		return fmt.Errorf("%w: book title is required", ErrDomainValidation)
	}
	if book.Author == "" {
		return fmt.Errorf("%w: book author is required", ErrDomainValidation)
	}
	if book.ISBN == "" {
		return fmt.Errorf("%w: book isbn is required", ErrDomainValidation)
	}
	if book.Price < 0 {
		return fmt.Errorf("%w: book price is negative", ErrDomainValidation)
	}
	return nil
}
