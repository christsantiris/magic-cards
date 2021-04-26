package models

type CardResponse struct {
	Card       Card
	Message    string
	StatusCode int
}

type CardsResponse struct {
	Cards      []Card
	Message    string
	StatusCode int
}
