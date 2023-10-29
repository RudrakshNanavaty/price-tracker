package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"price-tracker/entities"
)

type DB struct {
	products *gorm.DB
}

func NewDB() *DB {
	// connect to sqlite DB
	db, err := gorm.Open(sqlite.Open("urls.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	return &DB{
		products: db,
	}
}

// get everything from the DB (testing purposes)
func (db *DB) GetAll() (*[]entities.Product, error) {
	var products []entities.Product

	result := db.products.Find(&products)

	if result.Error != nil {
		return nil, result.Error
	} else if result.RowsAffected == 0 {
		return nil, result.Error
	}

	return &products, nil
}

func (db *DB) GetPrice(id int) (float64, error) {
	var product entities.Product

	// Use the First method with a condition to find the user by ID.
	result := db.products.Select("price").First(&product, id)

	if result.Error != nil {
		return 0, result.Error
	} else if result.RowsAffected == 0 {
		return 0, result.Error
	}

	return product.Price, nil
}

// add a product to track price in the DB
func (db *DB) AddProduct(product *entities.Product) (int, error) {

	// Use the First method with a condition to find the user by ID.
	result := db.products.Create(product)

	if result.Error != nil {
		return 0, result.Error
	} else if result.RowsAffected == 0 {
		return 0, result.Error
	}

	return product.ID, nil
}

// update the price of a product in the DB
func (db *DB) UpdatePrice(id int, price float64) error {

	db.products.Model(&entities.Product{ID: id}).Update("price", price)

	return nil
}
