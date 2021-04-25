package controllers

import (
	"github.com/christsantiris/magic-cards/models"
	"net/http"
	"encoding/json"
	"log"
	"database/sql"
	"github.com/gorilla/mux"
)

type Controller struct {}

var cards []models.Card

func logFatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (c Controller) GetCards(db *sql.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		var card models.Card
		cards = []models.Card{}

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
}

func (c Controller) GetCard(db *sql.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		var card models.Card
		params := mux.Vars(r)
	
		row := db.QueryRow("select * from cards where id=$1", params["id"])
	
		err := row.Scan(&card.ID, &card.Name, &card.Color, &card.StandardLegal, &card.Type, &card.Rarity, &card.Set, &card.CastingCost)
		logFatal(err)
	
		json.NewEncoder(w).Encode(card)
	}
}

func (c Controller) AddCard(db *sql.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		var card models.Card
		var cardID int
	
		json.NewDecoder(r.Body).Decode(&card)
	
		err := db.QueryRow("insert into cards (name, color, standard_legal, type, rarity, set, casting_cost) values($1, $2, $3, $4, $5, $6, $7) RETURNING id;", 
			card.Name, card.Color, card.StandardLegal, card.Type, card.Rarity, card.Set, card.CastingCost).Scan(&cardID)
		logFatal(err)
	
		json.NewEncoder(w).Encode(cardID)
	}
}

func (c Controller) UpdateCard(db *sql.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		var card models.Card
		json.NewDecoder(r.Body).Decode(&card)
	
		result, err := db.Exec("update cards set name=$1, color=$2, standard_legal=$3, type=$4, rarity=$5, set=$6, casting_cost=$7 where id=$8 RETURNING id", 
			&card.Name, &card.Color, &card.StandardLegal, &card.Type, &card.Rarity, &card.Set, &card.CastingCost, &card.ID)
	
		rowsUpdated, err := result.RowsAffected()
		logFatal(err)
	
		json.NewEncoder(w).Encode(rowsUpdated)
	}
}

func (c Controller) RemoveCard(db *sql.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		params := mux.Vars(r)

		result, err := db.Exec("delete from cards where id = $1", params["id"])
		logFatal(err)
	
		rowsDeleted, err := result.RowsAffected()
		logFatal(err)
	
		json.NewEncoder(w).Encode(rowsDeleted)
	}	
}