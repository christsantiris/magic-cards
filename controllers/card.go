package controllers

import (
	"github.com/christsantiris/magic-cards/models"
	"github.com/christsantiris/magic-cards/repository/card"
	"github.com/christsantiris/magic-cards/utils"
	"net/http"
	"encoding/json"
	"log"
	"database/sql"
	"github.com/gorilla/mux"
	"strconv"
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
		var error models.Error

		cards = []models.Card{}
		cardRepo := cardRepository.CardRepository{}
		cards, err := cardRepo.GetCards(db, card, cards)

		if err != nil {
			error.Message = "Server error"
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		utils.SendSuccess(w, cards)
	}
}

func (c Controller) GetCard(db *sql.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		var card models.Card
		var error models.Error

		params := mux.Vars(r)

		cards = []models.Card{}
		cardRepo := cardRepository.CardRepository{}

		id, _ := strconv.Atoi(params["id"])

		card, err := cardRepo.GetCard(db, card, id)

		if err != nil {
			if err == sql.ErrNoRows {
				error.Message = "Not Found"
				utils.SendError(w, http.StatusNotFound, error)
				return
			} else {
				error.Message = "Server error"
				utils.SendError(w, http.StatusInternalServerError, error)
				return
			}
		}

		w.Header().Set("Content-Type", "application/json")
		utils.SendSuccess(w, card)
	}
}

func (c Controller) AddCard(db *sql.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		var card models.Card
		var cardID int
		var error models.Error

		json.NewDecoder(r.Body).Decode(&card)

		if card.Name == "" {
			error.Message = "Card name is required."
			utils.SendError(w, http.StatusBadRequest, error) //400
			return
		}

		cardRepo := cardRepository.CardRepository{}
		cardID, err := cardRepo.AddCard(db, card)

		if err != nil {
			error.Message = "Server error"
			utils.SendError(w, http.StatusInternalServerError, error) //500
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		utils.SendSuccess(w, cardID)
	}
}

func (c Controller) UpdateCard(db *sql.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		var card models.Card
		var error models.Error

		json.NewDecoder(r.Body).Decode(&card)

		if card.ID == 0 {
			error.Message = "Card ID is required to update a card."
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}

		cardRepo := cardRepository.CardRepository{}
		rowsUpdated, err := cardRepo.UpdateCard(db, card)

		if err != nil {
			error.Message = "Server error"
			utils.SendError(w, http.StatusInternalServerError, error) //500
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		utils.SendSuccess(w, rowsUpdated)
	}
}

func (c Controller) RemoveCard(db *sql.DB) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		var error models.Error
		params := mux.Vars(r)
		cardRepo := cardRepository.CardRepository{}
		id, _ := strconv.Atoi(params["id"])

		rowsDeleted, err := cardRepo.RemoveCard(db, id)

		if err != nil {
			error.Message = "Server error."
			utils.SendError(w, http.StatusInternalServerError, error) //500
			return
		}

		if rowsDeleted == 0 {
			error.Message = "Card Not Found"
			utils.SendError(w, http.StatusNotFound, error) //404
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		utils.SendSuccess(w, rowsDeleted)
	}	
}