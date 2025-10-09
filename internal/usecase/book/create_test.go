package book

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"book-store-api/internal/models"
)

func TestCreate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		book models.BookParams

		wantErr bool
	}{
		{"valid creation", models.BookParams{
			Title:       "sd",
			Description: "sad",
			Author:      "asdsa",
			ISBN:        "asdad",
			Price:       0,
		}, false},
		{"invalid сreation", models.BookParams{
			Title:       "",
			Description: "",
			Author:      "",
			ISBN:        "",
			Price:       0,
		}, true},
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	repoMock := &RepositoryMock{
		CreateFunc: func(ctx context.Context, b models.Book) error {
			return nil
		},
	}

	cacheMock := &CacheMock{
		SetFunc: func(ctx context.Context, key string, value interface{}) error {
			return nil
		},
	}
	service := NewService(*logger, repoMock, cacheMock)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ctx := context.Background()
			id, err := service.Create(ctx, tt.book)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && id == "" {
				t.Errorf("expected UUID to be set, got empty string")
			}
		})
	}
}
