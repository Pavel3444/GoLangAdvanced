package migrations

import (
	"4-order-api/internal/product"
	"gorm.io/gorm"
)

func InitialMigration(db *gorm.DB) error {
	return db.AutoMigrate(&product.Product{})
}
