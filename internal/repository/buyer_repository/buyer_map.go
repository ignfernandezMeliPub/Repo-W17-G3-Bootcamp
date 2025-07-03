package buyer_repository

import (
	"app/pkg/custom_errors"
	"app/pkg/models"
)

type BuyerRepositoryMap struct {
	db   map[int]models.Buyer
	seed int
}

func NewBuyerMap(db map[int]models.Buyer) *BuyerRepositoryMap {
	defaultDb := make(map[int]models.Buyer)
	maxID := 0

	if db != nil {
		defaultDb = db
		for k := range db {
			if k > maxID {
				maxID = k
			}
		}
	}
	return &BuyerRepositoryMap{db: defaultDb, seed: maxID + 1}
}

func (r *BuyerRepositoryMap) GetAllBuyers() (b []models.Buyer, err error) {
	b = make([]models.Buyer, len(r.db))

	i := 0
	for _, value := range r.db {
		b[i] = value
		i++
	}

	if len(b) == 0 {
		err = custom_errors.ErrNotFound
		return
	}

	return
}

func (r *BuyerRepositoryMap) GetBuyerById(id int) (b models.Buyer, err error) {
	b, ok := r.db[id]

	if !ok {
		err = custom_errors.ErrNotFound
		return
	}

	return
}
func (r *BuyerRepositoryMap) GetBuyerByCardNumberId(cardNumberId string) (b models.Buyer, err error) {

	for _, value := range r.db {
		if value.CardNumberId == cardNumberId {
			b = value
			return
		}
	}

	err = custom_errors.ErrNotFound
	return
}
func (r *BuyerRepositoryMap) CreateBuyer(_b models.Buyer) (b models.Buyer, err error) {

	b = _b
	b.Id = r.seed
	r.seed++

	r.db[b.Id] = b
	return
}
func (r *BuyerRepositoryMap) UpdateBuyer(_b models.Buyer) (b models.Buyer, err error) {

	r.db[_b.Id] = _b
	b = _b

	return
}

func (r *BuyerRepositoryMap) DeleteBuyerById(id int) (err error) {
	delete(r.db, id)
	return
}
