// @title Book API
// @version 1.0
// @description CRUD по книгам
// @host book-store-api:8080
// @BasePath /api/v1/
// @schemes http

package httpv1

import (
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"log/slog"
	"net/http"

	_ "book-store-api/docs"
	"book-store-api/internal/converter"
	"book-store-api/internal/delivery"
	"book-store-api/internal/dto"
	"book-store-api/internal/models"
	"book-store-api/internal/repository"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler struct {
	usecase delivery.Usecase
	logger  *slog.Logger
}

func NewBookHandler(u delivery.Usecase, logger *slog.Logger) *Handler {
	return &Handler{usecase: u, logger: logger}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/book", h.GetAllBooks).Methods("GET")
	router.HandleFunc("/book/{id}", h.GetBookByID).Methods("GET")
	router.HandleFunc("/book", h.CreateBook).Methods("POST")
	router.HandleFunc("/book/{id}", h.UpdateBook).Methods("PUT")
	router.HandleFunc("/book/{id}", h.DeleteBook).Methods("DELETE")
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
}

// @Summary Получить все книги
// @Description Возвращает список всех книг
// @Tags books
// @Produce json
// @Success 200 {array} dto.BookDTO
// @Failure 500 {string} string "internal server error"
// @Router /book [get]
func (h *Handler) GetAllBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	books, err := h.usecase.GetAll(ctx)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	responseDTO := converter.ToBookResponseList(books)

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(responseDTO)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// @Summary Получить книгу по ID
// @Description Возвращает книгу по идентификатору
// @Tags books
// @Produce json
// @Param id path string true "Book ID"
// @Success 200 {object} dto.BookDTO
// @Failure 404 {string} string "not found"
// @Failure 500 {string} string "internal server error"
// @Router /book/{id} [get]
func (h *Handler) GetBookByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idParam := mux.Vars(r)["id"]
	if _, err := uuid.Parse(idParam); err != nil {
		http.Error(w, "invalid uuid format", http.StatusBadRequest)
		return
	}
	ctx := r.Context()
	book, err := h.usecase.GetByID(ctx, idParam)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)

			return
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}
	bookDTO := converter.ToBookResponse(*book)
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(bookDTO)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

}

// @Summary Создать книгу
// @Description Создает новую книгу
// @Tags books
// @Accept json
// @Produce json
// @Param book body dto.BookRequest true "Book data"
// @Success 201 {string} string "created id"
// @Failure 400 {string} string "invalid request body"
// @Failure 422 {string} string "validation error"
// @Failure 500 {string} string "internal server error"
// @Router /book [post]
func (h *Handler) CreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	var bookDTO dto.BookRequest
	if err := json.NewDecoder(r.Body).Decode(&bookDTO); err != nil {
		h.logger.Error("invalid request body", "err", err)
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	book := converter.ToBookParams(bookDTO)

	id, err := h.usecase.Create(ctx, book)
	if err != nil {
		h.logger.Error("failed to create book", "err", err)

		if errors.Is(err, models.ErrDomainValidation) {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(map[string]string{"id": id})
	if err != nil {
		return
	}

}

// @Summary Обновить книгу
// @Description Обновляет данные существующей книги
// @Tags books
// @Accept json
// @Produce json
// @Param id path string true "Book ID"
// @Param book body dto.BookRequest true "Book data"
// @Success 200 {string} string "ok"
// @Failure 400 {string} string "invalid request body"
// @Failure 404 {string} string "not found"
// @Failure 422 {string} string "validation error"
// @Failure 500 {string} string "internal server error"
// @Router /book/{id} [put]
func (h *Handler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	var bookDTO dto.BookRequest

	idParam := mux.Vars(r)["id"]
	uid, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, "invalid uuid format", http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&bookDTO)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)

		return
	}

	book := converter.ToBookParams(bookDTO)
	book.ID = uid

	err = h.usecase.Update(ctx, book)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)

			return
		}
		if errors.Is(err, models.ErrDomainValidation) {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)

			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)

}

// @Summary Удалить книгу
// @Description Удаляет книгу по ID
// @Tags books
// @Produce json
// @Param id path string true "Book ID"
// @Success 204 {string} string "no content"
// @Failure 404 {string} string "not found"
// @Failure 500 {string} string "internal server error"
// @Router /book/{id} [delete]
func (h *Handler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	idParam := mux.Vars(r)["id"]
	if _, err := uuid.Parse(idParam); err != nil {
		http.Error(w, "invalid uuid format", http.StatusBadRequest)
		return
	}

	err := h.usecase.DeleteBook(ctx, idParam)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)

			return
		}

		http.Error(w, "internal server error", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusNoContent)
}
