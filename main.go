package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/cards", getCards).Methods("GET")
	router.HandleFunc("/cards/{id}", getCard).Methods("GET")
	router.HandleFunc("/cards", addCard).Methods("POST")
	router.HandleFunc("/cards", updateCard).Methods("PUT")
	router.HandleFunc("/cards/{id}", removeCard).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8000", router))
}

func getCards(w http.ResponseWriter, r *http.Request) {
	log.Println("Get Cards")
}
func getCard(w http.ResponseWriter, r *http.Request) {
	log.Println("Get Card")
}
func addCard(w http.ResponseWriter, r *http.Request) {
	log.Println("Add Card")
}
func updateCard(w http.ResponseWriter, r *http.Request) {
	log.Println("Update Card")
}
func removeCard(w http.ResponseWriter, r *http.Request) {
	log.Println("Remove Card")
}