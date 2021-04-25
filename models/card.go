package models

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