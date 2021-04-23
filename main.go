package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"encoding/json"
	"strconv"
	"fmt"
)

type Card struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Color string `json:"color"`
	StandardLegal bool `json:"standard_legal"`
	Type string `json:"type"`
	Rarity string `json:"rarity"`
	Set string `json:"set"`
	CastingCost int `json:"casting_cost"`
}

var cards []Card

func main() {
	router := mux.NewRouter()

	cards = append(cards, Card{ID: 1, Name: "Bonecrusher Giant", Color: "Red", Set: "Throne of Eldraine", Type: "Creature", Rarity: "Rare", StandardLegal: true, CastingCost: 3},
		Card{ID: 2, Name: "Embercleave", Color: "Red", Set: "Throne of Eldraine", Type: "Artifact", Rarity: "Mythic Rare", StandardLegal: true, CastingCost: 6})

	router.HandleFunc("/cards", getCards).Methods("GET")
	router.HandleFunc("/cards/{id}", getCard).Methods("GET")
	router.HandleFunc("/cards", addCard).Methods("POST")
	router.HandleFunc("/cards", updateCard).Methods("PUT")
	router.HandleFunc("/cards/{id}", removeCard).Methods("DELETE")

	fmt.Println("The App is running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func getCards(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(cards)
}
func getCard(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	i, _ := strconv.Atoi(params["id"])

	for _, card := range cards {
		if card.ID == i {
			json.NewEncoder(w).Encode(&card)
		}
	}
}
func addCard(w http.ResponseWriter, r *http.Request) {
	var card Card

	json.NewDecoder(r.Body).Decode(&card)

	cards = append(cards, card)

	json.NewEncoder(w).Encode(cards)
}
func updateCard(w http.ResponseWriter, r *http.Request) {
	var card Card

	json.NewDecoder(r.Body).Decode(&card)

	for i, item := range cards {
		if item.ID == card.ID {
			cards[i] = card
		}
	}
	json.NewEncoder(w).Encode(cards)
}
func removeCard(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, _ := strconv.Atoi(params["id"])

	for i, item := range cards {
		if item.ID == id {
			cards = append(cards[:i], cards[i+1:]...)
		}
	}

	json.NewEncoder(w).Encode(cards)
}