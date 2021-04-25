package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"fmt"
	"database/sql"
	"github.com/subosito/gotenv"
	"github.com/christsantiris/magic-cards/models"
	"github.com/christsantiris/magic-cards/driver"
	"github.com/christsantiris/magic-cards/controllers"
)

var cards []models.Card
var db *sql.DB

func init() {
	gotenv.Load()
}

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	db = driver.ConnectDB()

	controller := controllers.Controller{}

	router := mux.NewRouter()

	router.HandleFunc("/cards", controller.GetCards(db)).Methods("GET")
	router.HandleFunc("/cards/{id}", controller.GetCard(db)).Methods("GET")
	router.HandleFunc("/cards", controller.AddCard(db)).Methods("POST")
	router.HandleFunc("/cards", controller.UpdateCard(db)).Methods("PUT")
	router.HandleFunc("/cards/{id}", controller.RemoveCard(db)).Methods("DELETE")

	fmt.Println("The App is running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}
