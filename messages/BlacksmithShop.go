package messages

type BlacksmithShop struct {
	Link string `json:"link"`
	Name string `json:"name"`
	OwnerTag *string `json:"ownerTag,omitempty"`
	OwnerName string `json:"ownerName"`
	OwnerCastle string `json:"ownerCastle"`
	Kind string `json:"kind"`
	Mana int `json:"mana"`
	Offers *[]BlacksmithOffer `json:"offers,omitempty"`
	Specialization *map[string]int `json:"specialization,omitempty"`
	QualityCraftLevel int `json:"qualityCraftLevel"`
	MaintenanceEnabled bool `json:"maintenanceEnabled"`
	MaintenanceCost int `json:"maintenanceCost"`
	GuildDiscount int `json:"guildDiscount"`
	CastleDiscount int `json:"castleDiscount"`
}