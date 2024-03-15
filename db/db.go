package db

import (
	"context"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// InitDatabase opens DB connection and sets up the database
func InitDatabase(ctx context.Context, dbName string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dbName))
	if err != nil {
		return nil, err
	}

	err = db.WithContext(ctx).AutoMigrate(&PackSize{})
	if err != nil {
		return nil, err
	}

	return db, nil
}
