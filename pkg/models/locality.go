package models

type Locality struct {
	Id           string `json:"id"            db:"id"`
	LocalityName string `json:"locality_name" db:"locality_name"`
	ProvinceName string `json:"province_name" db:"province_name"`
	CountryName  string `json:"country_name"  db:"country_name"`
}

type LocalitySellerCount struct {
	Id           string `json:"id"            db:"id"`
	LocalityName string `json:"locality_name" db:"locality_name"`
	SellersCount int    `json:"sellers_count" db:"sellers_count"`
}
