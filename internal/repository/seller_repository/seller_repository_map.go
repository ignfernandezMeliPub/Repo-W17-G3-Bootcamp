package seller_repository

import (
	"app/pkg/custom_errors"
	"app/pkg/models"
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

// Save adds or updates the Seller and returns it
func (r *SellerRepositoryMap) Save(seller models.Seller) (models.Seller, error) {
	r.lock.Lock()
	defer r.lock.Unlock()

	r.database[seller.Id] = seller
	return seller, nil
}

// GetById returns the Seller by id or an error if it does not exist
func (r *SellerRepositoryMap) GetById(id int) (models.Seller, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	seller, ok := r.database[id]
	if !ok {
		return models.Seller{}, &custom_errors.ResourceNotFoundError{}
	}
	return seller, nil
}

// IdExists checks if a seller with the given ID exists in the repository.
// Returns true if it exists, otherwise false. Error is always nil
func (r *SellerRepositoryMap) IdIsUsed(id int) (bool, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	_, ok := r.database[id]
	return ok, nil
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

// GetAll returns all the Sellers currently stored
func (r *SellerRepositoryMap) GetAll() ([]models.Seller, error) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	sellers := make([]models.Seller, 0, len(r.database))
	for _, seller := range r.database {
		sellers = append(sellers, seller)
	}

	return sellers, nil
}

// Delete removes a Seller by id
func (r *SellerRepositoryMap) Delete(id int) error {
	r.lock.Lock()
	defer r.lock.Unlock()

	_, ok := r.database[id]
	if !ok {
		return &custom_errors.ResourceNotFoundError{}
	}

	delete(r.database, id)
	return nil
}

// Update Updates the Seller and returns it
func (r *SellerRepositoryMap) Update(seller models.Seller) (models.Seller, error) {
	return r.Save(seller)
}
