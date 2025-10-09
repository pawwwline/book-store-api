package http

import (
	"book-store-api/internal/converter"
	"book-store-api/internal/delivery"
	"book-store-api/internal/dto"
	"book-store-api/internal/models"
	"book-store-api/internal/repository"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	usecase delivery.Usecase
}

func NewBookHandler(u delivery.Usecase) *Handler {
	return &Handler{usecase: u}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/book", h.GetAllBooks).Methods("GET")
	router.HandleFunc("/book/{id}", h.GetBookByID).Methods("GET")
	router.HandleFunc("/book", h.CreateBook).Methods("POST")
	router.HandleFunc("/book/{id}", h.UpdateBook).Methods("PUT")
	router.HandleFunc("/book/{id}", h.DeleteBook).Methods("DELETE")
}

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

func (h *Handler) GetBookByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	id := vars["id"]

	ctx := r.Context()
	book, err := h.usecase.GetByID(ctx, id)
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

func (h *Handler) CreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()

	var bookDTO dto.BookDTO
	err := json.NewDecoder(r.Body).Decode(&bookDTO)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)

		return
	}

	book := converter.ToBookParams(bookDTO)

	id, err := h.usecase.Create(ctx, book)
	if err != nil {
		if errors.Is(err, models.ErrDomainValidation) {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
		}
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	_, err = w.Write([]byte(id))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)

		return
	}

}

func (h *Handler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ctx := r.Context()
	var bookDTO dto.BookDTO
	err := json.NewDecoder(r.Body).Decode(&bookDTO)
	if err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)

		return
	}

	book := converter.ToBookParams(bookDTO)

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

func (h *Handler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	ctx := r.Context()

	vars := mux.Vars(r)
	id := vars["id"]

	err := h.usecase.DeleteBook(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			http.Error(w, err.Error(), http.StatusNotFound)

			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}
