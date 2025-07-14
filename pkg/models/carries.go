package models

type Carries struct {
	Id          int    `db:"id"`
	Cid         string `db:"cid"`
	CompanyName string `db:"company_name"`
	Address     string `db:"address"`
	Telephone   string `db:"telephone"`
	LocalityId  int    `db:"locality_id"`
}
