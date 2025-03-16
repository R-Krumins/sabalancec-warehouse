package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Storage interface {
	CreateProduct(*ProductFull) error
	GetProduct() ([]ProductSimple, error)
	GetProductById(int) (*ProductFull, error)
}

type SQLiteStorage struct {
	db *sql.DB
}

func NewSQLiteStorage(dbPath string) (*SQLiteStorage, error) {
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		err := CreateNewSqliteDB(dbPath)
		if err != nil {
			return nil, err
		}
	}

	db, err := sql.Open("sqlite3", dbPath)

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	fmt.Println("DB Connection established")
	return &SQLiteStorage{db: db}, nil
}

func CreateNewSqliteDB(dbPath string) error {
	// Create empty db file
	file, err := os.Create(dbPath)
	if err != nil {
		return fmt.Errorf("failed to create database file: %w", err)
	}
	file.Close()

	// Open database connection
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer db.Close()

	// Read and execute createDB.sql
	sqlBytes, err := os.ReadFile("./createDB.sql")
	if err != nil {
		return fmt.Errorf("failed to read createDB.sql: %w", err)
	}

	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		return fmt.Errorf("failed to execute createDB.sql: %w", err)
	}

	fmt.Println("New database created and initialized")
	return nil
}

func (s *SQLiteStorage) CreateProduct(p *ProductFull) error {
	query := `INSERT INTO products
        (name, price, category, image, amount_sold, amount_in_stock, has_allergens, rating)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?) RETURNING id`

	var id int
	err := s.db.QueryRow(
		query, p.Name, p.Price, p.Category, p.Image,
		p.AmountSold, p.AmountInStock, p.HasAllergens, p.Rating).Scan(&id)

	if err != nil {
		return err
	}

	// Update the product with the new ID
	p.Id = id
	return nil
}

func (s *SQLiteStorage) GetProduct() ([]ProductSimple, error) {
	rows, err := s.db.Query("SELECT id, name, category, image FROM products")
	if err != nil {
		return nil, err
	}

	products := make([]ProductSimple, 0)
	for rows.Next() {
		product := new(ProductSimple)
		err := rows.Scan(
			&product.Id, &product.Name, &product.Category, &product.Image)
		if err != nil {
			return nil, err
		}
		products = append(products, *product)
	}
	return products, nil

}

func (s *SQLiteStorage) GetProductById(id int) (*ProductFull, error) {
	rows, err := s.db.Query("SELECT * FROM products WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	product := new(ProductFull)
	for rows.Next() {
		err = rows.Scan(
			&product.Id, &product.Name, &product.Price, &product.Category,
			&product.Image, &product.AmountSold, &product.AmountInStock,
			&product.HasAllergens, &product.Rating)
		if err != nil {
			return nil, err
		}
		return product, nil
	}

	return nil, fmt.Errorf("no product with id %d", id)

}
