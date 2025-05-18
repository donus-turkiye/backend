package domain

type Category struct {
	CategoryId int    `json:"category_id"`
	WasteType  string `json:"waste_type"`
	UnitType   string `json:"unit_type"`
}
