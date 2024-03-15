package db

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PackSize struct {
	gorm.Model
	ID   string `gorm:"primaryKey"`
	Size int
}

type PackSizeRepository struct {
	db *gorm.DB
}

func NewPackSizeRepository(db *gorm.DB) *PackSizeRepository {
	return &PackSizeRepository{db: db}
}

func (r *PackSizeRepository) CreateSize(size int) error {
	result := r.db.Create(&PackSize{ID: uuid.New().String(), Size: size})
	return result.Error
}

func (r *PackSizeRepository) FindAll() []PackSize {
	var packSizes []PackSize
	r.db.Find(&packSizes)
	return packSizes
}

func (r *PackSizeRepository) DeleteAll() {
	r.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&PackSize{})
}
