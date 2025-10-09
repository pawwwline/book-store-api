package models

import (
	"testing"

	"github.com/google/uuid"
)

func TestValidBook(t *testing.T) {
	tests := []struct {
		name    string
		book    BookParams
		wantErr bool
	}{
		{"valid book", BookParams{
			ID:          uuid.New(),
			Title:       "B",
			Description: "C",
			Author:      "D",
			ISBN:        "kdslf1",
			Price:       10,
		}, false},
		{"empty book", BookParams{}, true},
		{"invalid book", BookParams{
			ID:     uuid.Nil,
			Title:  "B",
			Author: "D",
			ISBN:   "kdslf1",
			Price:  10,
		}, true},
		{"invalid price", BookParams{
			ID:     uuid.New(),
			Title:  "B",
			Author: "D",
			ISBN:   "kdslf1",
			Price:  -10,
		}, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateBook(tt.book); (err != nil) != tt.wantErr {
				t.Errorf("BookParams() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
