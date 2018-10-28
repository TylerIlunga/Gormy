package routes

import (
  "fmt"
  "encoding/json"
  "net/http"
  "github.com/gorilla/mux"
  "github.com/jinzhu/gorm"
  "practice/gormy/gormy-backend/models"
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
  var brands []models.Brand
  if err := database.Preload("Stores").Find(&brands).Error; err != nil {
    panic(err)
  }
  json.NewEncoder(w).Encode(&brands)
}

func GetAllBrandSneakers(w http.ResponseWriter, r *http.Request) {
  var sneakers []models.Sneaker
  if err := database.Preload("Brand").Find(&sneakers).Error; err != nil {
    panic(err)
  }
  json.NewEncoder(w).Encode(&sneakers)
}

func GetSpecificSneakerFromStore(w http.ResponseWriter, r *http.Request) {
  var sneakers models.Sneaker
  sneakerId := mux.Vars(r)["sneakerId"]

  err := database.Preload("Store").First(&sneakers, "sneaker_id = ?", sneakerId).Error;
  if err != nil {
    panic(err)
  }

  json.NewEncoder(w).Encode(&sneakers)
}

// Get brand that owns a certain sneaker
func GetSneakerBrand(w http.ResponseWriter, r *http.Request) {
  var brand models.Brand
  brandId := mux.Vars(r)["brandId"]

  err := database.Preload("Sneakers").First(&brand, "brand_id = ?", brandId).Error;
  if err != nil {
    panic(err)
  }

  json.NewEncoder(w).Encode(&brand)
}

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
  router.HandleFunc("/store/brands/all", GetAllStoreBrands).Methods("GET")
  router.HandleFunc("/brand/sneakers/all", GetAllBrandSneakers).Methods("GET")
  router.HandleFunc("/store/sneakers/{sneakerId}", GetSpecificSneakerFromStore).Methods("GET")
  router.HandleFunc("/sneakers/brand/{brandId}", GetSneakerBrand).Methods("GET")
  router.HandleFunc("/brand/delete/{brandId}", DeleteBrand).Methods("GET")
  return router
}
