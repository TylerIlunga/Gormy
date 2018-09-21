package models

import "github.com/jinzhu/gorm"

type Store struct {
  gorm.Model
  Name string
  Brands []Brand
}

type Brand struct {
  gorm.Model
  Name string
  Sneakers []Sneaker
  StoreID uint
}

type Sneaker struct {
  gorm.Model
  BrandID uint
  SneakerModel string
  Price int
  Supply int
}

func (brand *Brand) AfterDelete(scope *gorm.Scope) error {
  var sneaker Sneaker
  if err := scope.DB().Model(sneaker).Where("id=?", brand.ID).Delete(sneaker).Error; err != nil {
    panic(nil)
  }
  return nil
}
