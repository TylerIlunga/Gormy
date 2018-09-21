package routes

import (
  "fmt"
  "encoding/json"
  "net/http"
  "github.com/gorilla/mux"
  "github.com/jinzhu/gorm"
  "../../gormy-backend/models"
)

var database *gorm.DB

func CreateStore(w http.ResponseWriter, r *http.Request) {
  var store models.Store
  json.NewDecoder(r.Body).Decode(&store)
  database.Create(&store)
  json.NewEncoder(w).Encode(&store)
}

func CreateBrand(w http.ResponseWriter, r *http.Request) {
  var brand models.Brand
  json.NewDecoder(r.Body).Decode(&brand)
  if err := database.Create(&brand).Error; err != nil {
    fmt.Println(err)
    json.NewEncoder(w).Encode(err)
  }
  json.NewEncoder(w).Encode(&brand)
}

func CreateSneaker(w http.ResponseWriter, r *http.Request) {
  var sneaker models.Sneaker
  json.NewDecoder(r.Body).Decode(&sneaker)
  if err := database.Create(&sneaker).Error; err != nil {
     panic(err)
     return
  }
  json.NewEncoder(w).Encode(&sneaker)
}

func GetAllStoreBrands(w http.ResponseWriter, r *http.Request) {
  var store models.Store
  var storeBrands []models.Brand
  storeId := mux.Vars(r)["storeId"]

  if err := database.Find(&store, storeId).
                      Model(&store).
                      Related(&storeBrands).Error; err != nil {
    panic(err)
  }

  json.NewEncoder(w).Encode(&storeBrands)
}

// Get all sneakers a brand has
func GetAllBrandSneakers(w http.ResponseWriter, r *http.Request) {
  var brand models.Brand
  var brandSneakers []models.Sneaker
  brandId := mux.Vars(r)["brandId"]

  if err := database.Find(&brand, brandId).
                      Model(&brand).
                      Related(&brandSneakers).Error; err != nil {
    panic(err)
  }

  json.NewEncoder(w).Encode(&brandSneakers)
}

// Get stores carrying a specific sneaker
func GetSpecificSneakerFromStore(w http.ResponseWriter, r *http.Request) {
  var store models.Store
  var brand models.Brand
  var sneakers models.Sneaker
  storeId := mux.Vars(r)["storeId"]
  brandName := mux.Vars(r)["brandName"]
  sneakerModel := mux.Vars(r)["sneakerModel"]

  tx := database.Begin()
  if err := tx.Select("id").
                Where("ID = ?", storeId).
                Find(&store).Error; err != nil {
    panic(err)
    tx.Rollback()
  }
  if err := tx.Select("id").
                Where("store_id = ? AND name = ?", &store.ID, brandName).
                Find(&brand).Error; err != nil {
    panic(err)
    tx.Rollback()
  }
  if err := tx.Select("sneaker_model, price, supply").
                Where("brand_id = ? AND sneaker_model = ?", &brand.ID, sneakerModel).
                Find(&sneakers).Error; err != nil {
    panic(err)
    tx.Rollback()
  }

  tx.Commit()
  json.NewEncoder(w).Encode(&sneakers)
}

// Get brand that owns a certain sneaker
func GetSneakerBrand(w http.ResponseWriter, r *http.Request) {
  var sneaker models.Sneaker
  var brand models.Brand
  sneakerModel := mux.Vars(r)["sneakerModel"]

  tx := database.Begin()
  if err := tx.Where("sneaker_model = ?", sneakerModel).Find(&sneaker).Error; err != nil {
    panic(err)
    tx.Rollback()
  }
  if err := tx.Where("ID = ?", &sneaker.BrandID).
                Find(&brand).Error; err != nil {
    panic(err)
    tx.Rollback()
  }

  tx.Commit()
  json.NewEncoder(w).Encode(&brand)
}

// Delete brand (and the sneakers connect to it) from a certain store
func DeleteBrand(w http.ResponseWriter, r *http.Request) {
    var brand models.Brand
    brandId := mux.Vars(r)["brandId"]
    if err := database.Where("ID = ?", brandId).Find(&brand).Error; err != nil {
      panic(err)
    }
    database.Delete(&brand)
    json.NewEncoder(w).Encode(&brand)
}


func HomeHandler(w http.ResponseWriter, r *http.Request) {
  w.Header().Set("Content-Type", "text/html")
  http.ServeFile(w, r, "./static/views/index.html")
}

func GetRouter(db *gorm.DB) *mux.Router {
  database = db
  router := mux.NewRouter().StrictSlash(true)
  router.PathPrefix("/static/").
          Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("/static/"))))
  router.HandleFunc("/", HomeHandler).Methods("GET")
  router.HandleFunc("/create/store", CreateStore).Methods("POST")
  router.HandleFunc("/create/brand", CreateBrand).Methods("POST")
  router.HandleFunc("/create/sneaker", CreateSneaker).Methods("POST")
  router.HandleFunc("/store/brands/all/{storeId}", GetAllStoreBrands).Methods("GET")
  router.HandleFunc("/brand/sneakers/all/{brandId}", GetAllBrandSneakers).Methods("GET")
  router.HandleFunc("/store/sneakers/{storeId}", GetSpecificSneakerFromStore).
          Queries("brand", "{brandName}", "model", "{sneakerModel}").
          Methods("GET")
  router.HandleFunc("/sneakers/brand", GetSneakerBrand).
          Queries("model", "{sneakerModel}").
          Methods("GET")
  router.HandleFunc("/brand/delete/{brandId}", DeleteBrand).Methods("GET")
  return router
}
