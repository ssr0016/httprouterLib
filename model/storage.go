package db

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/ssr0016/librarySystem/types"
)

type Storage interface {
	CreateBook(ctx context.Context, book types.Book) error
	GetBook(ctx context.Context, id int64) (types.Book, error)
	GetBooks(ctx context.Context) ([]types.Book, error)
	DeleteBook(ctx context.Context, id int64) error
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	connstr := "user=postgres dbname=postgres password=secret  sslmode=disable"
	db, err := sql.Open("postgres", connstr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil

}

func (db *PostgresStore) CreateBook(ctx context.Context, book types.Book) error {

	_, err := db.db.ExecContext(ctx, "INSERT INTO books (title, price, created_at) VALUES ($1, $2, $3)",
		book.Title, book.Price, book.CreatedAt)
	if err != nil {
		return fmt.Errorf("failed to execute query: %v", err)
	}

	return nil
}

func (db *PostgresStore) GetBook(ctx context.Context, id int64) (types.Book, error) {
	var book types.Book

	row := db.db.QueryRowContext(ctx, "SELECT * FROM books WHERE id = $1", id)

	err := row.Scan(
		&book.ID,
		&book.Title,
		&book.Price,
		&book.CreatedAt,
	)

	if err != nil {
		return types.Book{}, fmt.Errorf("failed to execute query: %v", err)
	}

	return book, nil
}

func (db *PostgresStore) GetBooks(ctx context.Context) ([]types.Book, error) {
	rows, err := db.db.QueryContext(ctx, "SELECT * FROM books")
	if err != nil {
		return nil, fmt.Errorf("failed to execute query: %v", err)
	}
	defer rows.Close()

	books := []types.Book{}
	for rows.Next() {
		var book types.Book
		err := rows.Scan(
			&book.ID,
			&book.Title,
			&book.Price,
			&book.CreatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to execute query: %v", err)
		}

		books = append(books, book)
	}

	return books, nil
}

func (db *PostgresStore) DeleteBook(ctx context.Context, id int64) error {
	_, err := db.db.ExecContext(ctx, "DELETE FROM books WHERE id = $1", id)
	if err != nil {
		return fmt.Errorf("failed to execute query: %v", err)
	}

	return nil
}
