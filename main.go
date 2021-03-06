package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/christsantiris/magic-cards/controllers"
	"github.com/christsantiris/magic-cards/driver"
	"github.com/christsantiris/magic-cards/utils"
	"github.com/gorilla/mux"
	"github.com/subosito/gotenv"
)

var db *sql.DB

func init() {
	gotenv.Load()
}

func main() {
	db = driver.ConnectDB()

	controller := controllers.Controller{}

	router := mux.NewRouter()

	// Card routes
	router.HandleFunc("/cards", utils.TokenVerifyMiddleWare(utils.Logging(controller.GetCards(db)))).Methods("GET")
	router.HandleFunc("/cards/{id}", utils.TokenVerifyMiddleWare(utils.Logging(controller.GetCard(db)))).Methods("GET")
	router.HandleFunc("/cards", utils.TokenVerifyMiddleWare(utils.Logging(controller.AddCard(db)))).Methods("POST")
	router.HandleFunc("/cards", utils.TokenVerifyMiddleWare(utils.Logging(controller.UpdateCard(db)))).Methods("PUT")
	router.HandleFunc("/cards/{id}", utils.TokenVerifyMiddleWare(utils.Logging(controller.RemoveCard(db)))).Methods("DELETE")

	// User routes
	router.HandleFunc("/signup", controller.Signup(db)).Methods("POST")
	router.HandleFunc("/login", controller.Login(db)).Methods("POST")

	fmt.Println("The App is running on port 8000")
	utils.LogFatal(http.ListenAndServe(":8000", router))
}
