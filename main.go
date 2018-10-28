package main

import (
  "fmt"
  "log"
  "os"
  "net/http"
  "practice/gormy/gormy-backend/routes"
  "practice/gormy/gormy-backend/models"
  "github.com/gorilla/handlers"
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB

func listen() {
  var err error
  connectString := "host=127.0.0.1 port=5432 user=tilios dbname=gormy sslmode=disable"
  db, err = gorm.Open("postgres", connectString)
  if err != nil {
    panic(err)
  }

  defer db.Close()
  db.LogMode(true)
  db.AutoMigrate(&models.Store{}, &models.Brand{}, &models.Sneaker{})
  fmt.Println("Connected to DB!")

  port := ":8081"
  if os.Getenv("PORT") != "" {
    port = ":" + os.Getenv("PORT")
  }
  fmt.Println("Listening on port" + port)

  log.Fatal(http.ListenAndServe(
    port, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"}),
    handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"}),
    handlers.AllowedOrigins([]string{"*"}))(routes.GetRouter(db))))
}

func main() {
  fmt.Println("Running...")
  listen()
}
