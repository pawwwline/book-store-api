package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"book-store-api/internal/models"
)

type BookRepository struct {
	pool *pgxpool.Pool
}

func NewBookRepository(pool *pgxpool.Pool) *BookRepository {
	return &BookRepository{pool: pool}
}

func (r *BookRepository) Create(ctx context.Context, book models.Book) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO books (uuid, title, description, author, isbn, price, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())`,
		book.ID, book.Title, book.Description, book.Author, book.ISBN, book.Price,
	)
	return err
}

func (r *BookRepository) GetAll(ctx context.Context) ([]models.Book, error) {
	rows, err := r.pool.Query(ctx, `SELECT id, uuid, title, description, author, isbn, price, created_at, updated_at FROM books`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var b models.Book
		if err := rows.Scan(&b.ID, &b.ID, &b.Title, &b.Description, &b.Author, &b.ISBN, &b.Price, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, nil
}

func (r *BookRepository) GetById(ctx context.Context, id string) (models.Book, error) {
	var b models.Book
	err := r.pool.QueryRow(ctx,
		`SELECT id, uuid, title, description, author, isbn, price, created_at, updated_at FROM books WHERE uuid=$1`, id,
	).Scan(&b.ID, &b.ID, &b.Title, &b.Description, &b.Author, &b.ISBN, &b.Price, &b.CreatedAt, &b.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return models.Book{}, ErrNotFound
	}
	if err != nil {
		return models.Book{}, err
	}
	return b, nil
}

func (r *BookRepository) Update(ctx context.Context, book models.Book) error {
	commandTag, err := r.pool.Exec(ctx,
		`UPDATE books SET title=$1, description=$2, author=$3, isbn=$4, price=$5, updated_at=NOW() WHERE uuid=$6`,
		book.Title, book.Description, book.Author, book.ISBN, book.Price, book.ID,
	)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *BookRepository) Delete(ctx context.Context, id string) error {
	commandTag, err := r.pool.Exec(ctx, `DELETE FROM books WHERE uuid=$1`, id)
	if err != nil {
		return err
	}
	if commandTag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}

func (r *BookRepository) GetAllWithLimit(ctx context.Context, limit int) ([]models.Book, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, uuid, title, description, author, isbn, price, created_at, updated_at FROM books ORDER BY id LIMIT $1`, limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var books []models.Book
	for rows.Next() {
		var b models.Book
		if err := rows.Scan(&b.ID, &b.ID, &b.Title, &b.Description, &b.Author, &b.ISBN, &b.Price, &b.CreatedAt, &b.UpdatedAt); err != nil {
			return nil, err
		}
		books = append(books, b)
	}
	return books, nil
}
