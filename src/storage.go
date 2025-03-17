package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

type Storage interface {
	CreateProduct(*ProductFull) error
	GetProduct() ([]ProductSimple, error)
	GetProductById(int) (*ProductFull, error)

	GetAllergen() ([]AllergenSimple, error)
	GetAllergenByID(id int) (*AllergenFull, error)
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

	if err := executeSQLFile(db, "./sql/createDB.sql"); err != nil {
		return err
	}

	if err := executeSQLFile(db, "./sql/insertData.sql"); err != nil {
		return err
	}

	fmt.Println("New database created and initialized")
	return nil
}

func executeSQLFile(db *sql.DB, file string) error {
	sqlBytes, err := os.ReadFile(file)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", file, err)
	}
	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		return fmt.Errorf("failed to execute %s: %w", file, err)
	}

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

func (s *SQLiteStorage) GetAllergen() ([]AllergenSimple, error) {
	rows, err := s.db.Query("SELECT id, name, image FROM allergens")
	if err != nil {
		return nil, err
	}

	allergens := make([]AllergenSimple, 0)
	for rows.Next() {
		allergen := new(AllergenSimple)
		err := rows.Scan(
			&allergen.Id, &allergen.Name, &allergen.Image)
		if err != nil {
			return nil, err
		}
		allergens = append(allergens, *allergen)
	}
	return allergens, nil
}

func (s *SQLiteStorage) GetAllergenByID(id int) (*AllergenFull, error) {
	row := s.db.QueryRow("SELECT id, name, image, info FROM allergens WHERE id = ?", id)

	allergen := new(AllergenFull)

	// Create a variable to store the JSON string
	// The json string should be formatted as "section:text","section:text","section:text"
	var infoJSON string

	err := row.Scan(&allergen.Id, &allergen.Name, &allergen.Image, &infoJSON)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no allergen with id %d", id)
		}
		return nil, err
	}

	// Parse the JSON string into a map first
	var infoMap map[string]string
	err = json.Unmarshal([]byte(infoJSON), &infoMap)
	if err != nil {
		return nil, err
	}

	// Convert the map to an array of AllergenInfo
	allergen.Info = make([]AllergenInfo, 0, len(infoMap))
	for section, text := range infoMap {
		allergen.Info = append(allergen.Info, AllergenInfo{
			Section: section,
			Text:    text,
		})
	}

	return allergen, nil
}
