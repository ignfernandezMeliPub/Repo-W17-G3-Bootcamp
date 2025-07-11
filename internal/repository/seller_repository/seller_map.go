package seller_repository

import (
	"app/pkg/custom_errors"
	"app/pkg/models"
	"math/rand"
	"sync"
)

type SellerRepositoryMap struct {
	lock     sync.RWMutex
	database map[int]models.Seller
}

func NewSellerRepositoryMap(database map[int]models.Seller) SellerRepositoryMap {
	if database == nil {
		database = map[int]models.Seller{}
	}

	return SellerRepositoryMap{lock: sync.RWMutex{}, database: database}
}

// CreateSeller adds or updates the Seller and returns it
func (r *SellerRepositoryMap) CreateSeller(seller models.Seller) (models.Seller, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	var newId int
	for {
		newId = rand.Intn(999999) + 1 // ID entre 1 y 999999
		_, ok := r.database[newId]
		if !ok {
			break
		}
	}

	seller.Id = newId
	r.database[seller.Id] = seller
	return seller, nil
}

// GetSellerById returns the Seller by id or an error if it does not exist
func (r *SellerRepositoryMap) GetSellerById(id int) (models.Seller, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	seller, ok := r.database[id]
	if !ok {
		return models.Seller{}, &custom_errors.ResourceNotFoundError{}
	}
	return seller, nil
}

// CompanyIdIsUsed checks if any seller in the repository is using the specified company ID.
// Returns true if the company ID is found, otherwise false. Error is always nil
func (r *SellerRepositoryMap) CompanyIdIsUsed(companyId int) (bool, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	for _, seller := range r.database {
		if seller.CompanyId == companyId {
			return true, nil
		}
	}

	return false, nil
}

// GetAllSellers GetAllSeller returns all the Sellers currently stored
func (r *SellerRepositoryMap) GetAllSellers() ([]models.Seller, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	sellers := make([]models.Seller, 0, len(r.database))
	for _, seller := range r.database {
		sellers = append(sellers, seller)
	}

	return sellers, nil
}

// DeleteSellerById DeleteSeller removes a Seller by id
func (r *SellerRepositoryMap) DeleteSellerById(id int) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	_, ok := r.database[id]
	if !ok {
		return &custom_errors.ResourceNotFoundError{}
	}

	delete(r.database, id)
	return nil
}

// UpdateSellerById UpdateSeller Updates the Seller and returns it
func (r *SellerRepositoryMap) UpdateSellerById(seller models.Seller) (models.Seller, error) {
	r.database[seller.Id] = seller
	return seller, nil
}
