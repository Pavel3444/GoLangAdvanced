package db

import (
	"4-order-api/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func NewDb(conf config.Config) *gorm.DB {
	db, err := gorm.Open(postgres.Open(conf.Db.Dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}
