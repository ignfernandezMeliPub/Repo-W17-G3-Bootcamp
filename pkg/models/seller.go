package models

type Seller struct {
	Id          int    `json:"id"           db:"id"`
	CompanyId   int    `json:"cid"          db:"cid"`
	CompanyName string `json:"company_name" db:"company_name"`
	Address     string `json:"address"      db:"address"`
	Telephone   string `json:"telephone"    db:"telephone"`
	LocalityId  string `json:"locality_id"  db:"locality_id"`
}
