package seller_repository

import (
	"app/internal/repository/sql_utils"
	"app/pkg/custom_errors"
	"app/pkg/models"
	"database/sql"
)

type SellerRepositorySql struct {
	db *sql.DB
}

func NewSellerRepositorySql(db *sql.DB) SellerRepositorySql {
	return SellerRepositorySql{db: db}
}

// CreateSeller adds or updates the Seller and returns it
func (r *SellerRepositorySql) CreateSeller(seller models.Seller) (models.Seller, error) {
	newId, err := sql_utils.Insert(r.db, "INSERT INTO seller (cid, company_name, address, telephone) VALUES (?, ?, ?, ?)", []any{seller.CompanyId, seller.CompanyName, seller.Address, seller.Telephone})
	if err != nil {
		return seller, err
	}

	seller.Id = int(newId)
	return seller, nil
}

// GetSellerById returns the Seller by id or an error if it does not exist
func (r *SellerRepositorySql) GetSellerById(id int) (models.Seller, error) {
	return sql_utils.QueryRow[models.Seller](r.db, "SELECT (id, cid, company_name, address, telephone) FROM seller WHERE id = ?", []any{id})
}

// GetAllSellers GetAllSeller returns all the Sellers currently stored
func (r *SellerRepositorySql) GetAllSellers() ([]models.Seller, error) {
	return sql_utils.Query[models.Seller](r.db, "SELECT (id, cid, company_name, address, telephone) FROM seller", []any{})
}

// DeleteSellerById DeleteSeller removes a Seller by id
func (r *SellerRepositorySql) DeleteSellerById(id int) error {
	affectedRows, err := sql_utils.Delete(r.db, "DELETE FROM seller WHERE id = ?", []any{id})
	if err != nil {
		return err
	}
	if affectedRows == 0 {
		return custom_errors.ErrNotFound
	}

	return nil
}

// UpdateSellerById UpdateSeller Updates the Seller and returns it
func (r *SellerRepositorySql) UpdateSellerById(seller models.Seller) (models.Seller, error) {
	panic("TODO") // TODO
}
