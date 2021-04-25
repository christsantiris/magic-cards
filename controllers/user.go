package controllers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/christsantiris/magic-cards/utils"

	"github.com/christsantiris/magic-cards/models"
	userRepository "github.com/christsantiris/magic-cards/repository/user"

	"github.com/badoux/checkmail"
	"golang.org/x/crypto/bcrypt"
)

var users []models.User

func (c Controller) Login(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		var jwt models.JWT
		var error models.Error

		json.NewDecoder(r.Body).Decode(&user)

		if user.Email == "" {
			error.Message = "Email is missing."
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}

		if user.Password == "" {
			error.Message = "Password is missing."
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}

		password := user.Password

		userRepo := userRepository.UserRepository{}
		user, err := userRepo.Login(db, user)

		log.Println(err)

		if err != nil {
			if err == sql.ErrNoRows {
				error.Message = "The user does not exist"
				utils.SendError(w, http.StatusBadRequest, error)
				return
			} else {
				log.Fatal(err)
			}
		}

		hashedPassword := user.Password

		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))

		if err != nil {
			error.Message = "Invalid Password"
			utils.SendError(w, http.StatusUnauthorized, error)
			return
		}

		token, err := utils.GenerateToken(user)

		if err != nil {
			log.Fatal(err)
		}

		w.WriteHeader(http.StatusOK)
		jwt.Token = token

		utils.SendSuccess(w, jwt)
	}
}

func (c Controller) Signup(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		var error models.Error

		json.NewDecoder(r.Body).Decode(&user)

		if user.Email == "" {
			error.Message = "Email is missing."
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}

		if user.Email != "" {
			if err := checkmail.ValidateFormat(user.Email); err != nil {
				error.Message = "Email is invalid."
				utils.SendError(w, http.StatusInternalServerError, error)
				return
			}
		}

		if user.Password == "" {
			error.Message = "Password is missing."
			utils.SendError(w, http.StatusBadRequest, error)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)

		if err != nil {
			log.Fatal(err)
		}

		user.Password = string(hash)

		userRepo := userRepository.UserRepository{}
		user, err = userRepo.Signup(db, user)

		if err != nil {
			error.Message = "Server error."
			utils.SendError(w, http.StatusInternalServerError, error)
			return
		}

		user.Password = ""

		w.Header().Set("Content-Type", "application/json")
		utils.SendSuccess(w, user)

		json.NewEncoder(w).Encode(user)
	}
}
