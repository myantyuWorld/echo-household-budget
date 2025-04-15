package repository

import (
	"gorm.io/gorm"
)

type ShopRepository struct {
	db *gorm.DB
}

// func NewShopRepository(db *gorm.DB) domainmodel.ShoppingRepository {
// 	return &ShopRepository{db: db}
// }
