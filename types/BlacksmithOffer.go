package types

type BlacksmithOffer struct {
	Item  string `json:"item"`
	Price int    `json:"price"`
	Mana  int    `json:"mana"`
}
