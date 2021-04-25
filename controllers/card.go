package controllers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/christsantiris/magic-cards/models"
	cardRepository "github.com/christsantiris/magic-cards/repository/card"
	"github.com/christsantiris/magic-cards/utils"
	"github.com/gorilla/mux"
)

type Controller struct{}

var cards []models.Card

func (c Controller) GetCards(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS, PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
		utils.SendSuccess(w, cards)
	}
}

func (c Controller) GetCard(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS, PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
		utils.SendSuccess(w, card)
	}
}

func (c Controller) AddCard(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS, PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
		utils.SendSuccess(w, cardID)
	}
}

func (c Controller) UpdateCard(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var card models.Card
		var prevCard models.Card
		var error models.Error
		var parameters map[string]interface{}

		json.NewDecoder(r.Body).Decode(&parameters)

		// Check if params with zero values exist in request body before marshalling struct
		var casting_cost_empty bool
		_, ok := parameters["CastingCost"]
		if !ok {
			casting_cost_empty = true
		}

		var standard_legal_empty bool
		_, ok = parameters["StandardLegal"]
		if !ok {
			standard_legal_empty = true
		}

		// Marshal struct from map
		jsonString, _ := json.Marshal(parameters)
		json.Unmarshal(jsonString, &card)

		cardRepo := cardRepository.CardRepository{}

		// Get previous card values
		prevCard, err := cardRepo.GetCard(db, prevCard, card.ID)
		if err != nil {
			error.Message = "Card not found."
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}

		// Card.ID is required.
		// If other fields not provided to update, use existing values
		if card.ID == 0 {
			error.Message = "Card ID is required to update a card."
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}

		if card.Name == "" {
			card.Name = prevCard.Name
		}

		if card.Color == "" {
			card.Color = prevCard.Color
		}

		if card.Type == "" {
			card.Type = prevCard.Type
		}

		if card.Rarity == "" {
			card.Rarity = prevCard.Rarity
		}

		if card.Set == "" {
			card.Set = prevCard.Set
		}

		if casting_cost_empty {
			card.CastingCost = prevCard.CastingCost
		}

		if standard_legal_empty {
			card.StandardLegal = prevCard.StandardLegal
		}

		rowsUpdated, err := cardRepo.UpdateCard(db, card)

		if err != nil {
			error.Message = "Server error"
			utils.SendError(w, http.StatusInternalServerError, error) //500
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS, PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
		utils.SendSuccess(w, rowsUpdated)
	}
}

func (c Controller) RemoveCard(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "DELETE, POST, GET, OPTIONS, PUT")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
		utils.SendSuccess(w, rowsDeleted)
	}
}
