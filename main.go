package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"encoding/json"
	"strconv"
	"fmt"
	"os"
	"database/sql"
	"github.com/lib/pq"
	"github.com/subosito/gotenv"
)

type Card struct {
	ID int `json:"id"`
	Name string `json:"name"`
	Color string `json:"color"`
	StandardLegal bool `json:",omitempty"`
	Type string `json:"type"`
	Rarity string `json:"rarity"`
	Set string `json:"set"`
	CastingCost int `json:",omitempty"`
}

var cards []Card
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

	pgUrl, err := pq.ParseURL(os.Getenv("ELEPHANTSQL_URL"))
	logFatal(err)

	db, err = sql.Open("postgres", pgUrl)
	logFatal(err)

	err = db.Ping()
	logFatal(err)

	router := mux.NewRouter()

	router.HandleFunc("/cards", getCards).Methods("GET")
	router.HandleFunc("/cards/{id}", getCard).Methods("GET")
	router.HandleFunc("/cards", addCard).Methods("POST")
	router.HandleFunc("/cards", updateCard).Methods("PUT")
	router.HandleFunc("/cards/{id}", removeCard).Methods("DELETE")

	fmt.Println("The App is running on port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func getCards(w http.ResponseWriter, r *http.Request) {
	var card Card
	cards = []Card{}

	rows, err := db.Query("select * from cards")
	logFatal(err)

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&card.ID, &card.Name, &card.Color, &card.StandardLegal, &card.Type, &card.Rarity, &card.Set, &card.CastingCost)
		logFatal(err)

		cards = append(cards, card)
	}

	json.NewEncoder(w).Encode(cards)
}
func getCard(w http.ResponseWriter, r *http.Request) {
	var card Card
	params := mux.Vars(r)

	row := db.QueryRow("select * from cards where id=$1", params["id"])

	err := row.Scan(&card.ID, &card.Name, &card.Color, &card.StandardLegal, &card.Type, &card.Rarity, &card.Set, &card.CastingCost)
	logFatal(err)

	json.NewEncoder(w).Encode(card)
}
func addCard(w http.ResponseWriter, r *http.Request) {
	var card Card
	var cardID int

	json.NewDecoder(r.Body).Decode(&card)

	err := db.QueryRow("insert into cards (name, color, standard_legal, type, rarity, set, casting_cost) values($1, $2, $3, $4, $5, $6, $7) RETURNING id;", 
		card.Name, card.Color, card.StandardLegal, card.Type, card.Rarity, card.Set, card.CastingCost).Scan(&cardID)
	logFatal(err)

	json.NewEncoder(w).Encode(cardID)
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