package types

import "time"

// Book struct represents a book in the library
type Book struct {
	ID        int64     `json:"id"`
	Title     string    `json:"title"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateBookRequest struct {
	Title string  `json:"title"`
	Price float64 `json:"price"`
}

func NewBook(title string, price float64) (*Book, error) {
	return &Book{
		Title:     title,
		Price:     price,
		CreatedAt: time.Now().UTC(),
	}, nil
}
