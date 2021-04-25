package cardRepository

import (
	"github.com/christsantiris/magic-cards/models"
	"database/sql"
)

type CardRepository struct{}

func (c CardRepository) GetCards(db *sql.DB, card models.Card, cards []models.Card) ([]models.Card, error) {
	rows, err := db.Query("select * from cards")

	if err != nil {
		return []models.Card{}, err
	}

	for rows.Next() {
		err := rows.Scan(&card.ID, &card.Name, &card.Color, &card.StandardLegal, &card.Type, &card.Rarity, &card.Set, &card.CastingCost)
		if err != nil {
			return []models.Card{}, err
		}
		cards = append(cards, card)
	}

	return cards, nil
}

func (c CardRepository) GetCard(db *sql.DB, card models.Card, id int) (models.Card, error) {
	rows := db.QueryRow("select * from cards where id=$1", id)
	err := rows.Scan(&card.ID, &card.Name, &card.Color, &card.StandardLegal, &card.Type, &card.Rarity, &card.Set, &card.CastingCost)

	return card, err
}

func (c CardRepository) AddCard(db *sql.DB, card models.Card) (int, error) {
	err := db.QueryRow("insert into cards (name, color, standard_legal, type, rarity, set, casting_cost) values($1, $2, $3, $4, $5, $6, $7) RETURNING id;", 
	card.Name, card.Color, card.StandardLegal, card.Type, card.Rarity, card.Set, card.CastingCost).Scan(&card.ID)

	if err != nil {
		return 0, err
	}

	return card.ID, nil
}

func (c CardRepository) UpdateCard(db *sql.DB, card models.Card) (int64, error) {
	result, err := db.Exec("update cards set name=$1, color=$2, standard_legal=$3, type=$4, rarity=$5, set=$6, casting_cost=$7 where id=$8 RETURNING id", 
	&card.Name, &card.Color, &card.StandardLegal, &card.Type, &card.Rarity, &card.Set, &card.CastingCost, &card.ID)

	if err != nil {
		return 0, err
	}

	rowsUpdated, err := result.RowsAffected()

	if err != nil {
		return 0, err
	}

	return rowsUpdated, nil
}

func (c CardRepository) RemoveCard(db *sql.DB, id int) (int64, error) {
	result, err := db.Exec("delete from cards where id = $1", id)

	if err != nil {
		return 0, err
	}

	rowsDeleted, err := result.RowsAffected()

	if err != nil {
		return 0, err
	}

	return rowsDeleted, nil
}