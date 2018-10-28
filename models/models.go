package models

import (
  "time"
  "github.com/jinzhu/gorm"
)

type Store struct {
  StoreID uint `gorm:"primary_key:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
  Name string
  Brands []Brand `gorm:"many2many:brand_stores;association_foreignkey:BrandID;foreignkey:StoreID"`
  BrandID uint
  Sneakers []Sneaker `gorm:"foreignkey:BrandID;association_foreignkey:SneakerID"`
  SneakerID uint
}

type Brand struct {
  BrandID uint `gorm:"primary_key:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
  Name string
  Stores []Store `gorm:"many2many:brand_stores;association_foreignkey:StoreID;foreignkey:BrandID"`
  StoreID uint
  Sneakers []Sneaker `gorm:"foreignkey:BrandID;association_foreignkey:SneakerID"`
  SneakerID uint
}

type Sneaker struct {
  SneakerID uint `gorm:"primary_key:true"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
  Price int
  Supply int
  Brand Brand `gorm:"foreignkey:SneakerID;association_foreignkey:BrandID"`
  BrandID uint
  Store Store `gorm:"foreignkey:SneakerID;association_foreignkey:StoreID"`
  StoreID uint
}

func (brand *Brand) AfterDelete(scope *gorm.Scope) error {
  var sneaker Sneaker
  if err := scope.DB().Model(sneaker).Where("id=?", brand.BrandID).Delete(sneaker).Error; err != nil {
    panic(err)
  }
  return nil
}
