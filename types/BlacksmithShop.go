package types

import (
	"blacksmith-watcher/utils"
	"context"
	"time"
)

type BlacksmithShop struct {
	Link               string             `json:"link"`
	Name               string             `json:"name"`
	OwnerTag           *string            `json:"ownerTag,omitempty"`
	OwnerName          string             `json:"ownerName"`
	OwnerCastle        string             `json:"ownerCastle"`
	Kind               string             `json:"kind"`
	Mana               int                `json:"mana"`
	Offers             *[]BlacksmithOffer `json:"offers,omitempty"`
	Specialization     *map[string]int    `json:"specialization,omitempty"`
	QualityCraftLevel  int                `json:"qualityCraftLevel"`
	MaintenanceEnabled bool               `json:"maintenanceEnabled"`
	MaintenanceCost    int                `json:"maintenanceCost"`
	GuildDiscount      int                `json:"guildDiscount"`
	CastleDiscount     int                `json:"castleDiscount"`
	LastOpenTime       time.Time          `json:"last_open_time"`
}

func UpdateShop(shop BlacksmithShop) {
	var result BlacksmithShop
	filter := struct {
		Link string `json:"id"`
	}{Link: shop.Link}
	findUserResult := utils.DB.Collection("shops").FindOne(context.TODO(), filter)
	if findUserResult.Err() == nil {
		findUserResult.Decode(&result)
	} else {
		utils.DB.Collection("shops").InsertOne(context.TODO(), shop)
	}
}