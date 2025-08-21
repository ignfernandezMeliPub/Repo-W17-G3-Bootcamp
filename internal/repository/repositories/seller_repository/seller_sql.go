package seller_repository

import (
	"app/internal/logger"
	"app/internal/repository/sql_utils"
	"app/pkg/custom_errors"
	"app/pkg/models"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
)

/* TODOs:
2. Mejorar mapeo de errores
3. Mejorar organizacion de carpetas y directorios
4. Mejorar manejo de config
*/

type SellerRepositorySql struct {
	db *sql.DB
}

func NewSellerRepositorySql(db *sql.DB) SellerRepositorySql {
	return SellerRepositorySql{db: db}
}

// CreateSeller adds or updates the Seller and returns it
func (r *SellerRepositorySql) CreateSeller(seller models.Seller) (models.Seller, error) {
	sql_utils.LogAudit("CreateSeller", logger.LogStatusInProgress, "Insert seller")

	newId, err := sql_utils.Insert(r.db, "INSERT INTO sellers (cid, company_name, address, telephone, locality_id) VALUES (?, ?, ?, ?, ?)", []any{seller.CompanyId, seller.CompanyName, seller.Address, seller.Telephone, seller.LocalityId})
	if err != nil {
		err = sql_utils.HandleSqlError(err)
		sql_utils.LogAuditError("CreateSeller", "Insert seller", err)
		return seller, err
	}

	seller.Id = int(newId)

	sql_utils.LogAudit("CreateSeller", logger.LogStatusSuccess, "Insert seller. Id: "+strconv.Itoa(seller.Id))
	return seller, nil
}

// GetSellerById returns the Seller by id or an error if it does not exist
func (r *SellerRepositorySql) GetSellerById(id int) (models.Seller, error) {
	sql_utils.Log("GetSellerById", logger.LogStatusInProgress, "Select seller by id "+strconv.Itoa(id))

	s, err := sql_utils.QueryRow[models.Seller](r.db, "SELECT id, cid, company_name, address, telephone, locality_id FROM sellers WHERE id = ?", []any{id})
	if err != nil {
		err = sql_utils.HandleSqlError(err)
		sql_utils.LogError("GetSellerById", "Select seller by id "+strconv.Itoa(id), err)
		return models.Seller{}, err
	}

	sql_utils.Log("GetSellerById", logger.LogStatusSuccess, "Select seller by id "+strconv.Itoa(id))
	return s, nil
}

// GetAllSellers GetAllSeller returns all the Sellers currently stored
func (r *SellerRepositorySql) GetAllSellers() ([]models.Seller, error) {
	sql_utils.Log("GetAllSellers", logger.LogStatusInProgress, "Select all sellers")

	s, err := sql_utils.Query[models.Seller](r.db, "SELECT id, cid, company_name, address, telephone, locality_id FROM sellers", []any{})
	if err != nil {
		err = sql_utils.HandleSqlError(err)
		sql_utils.LogError("GetAllSellers", "Select all sellers", err)
		return nil, err
	}

	sql_utils.Log("GetAllSellers", logger.LogStatusSuccess, "Select all sellers")
	return s, nil
}

// DeleteSellerById DeleteSeller removes a Seller by id
func (r *SellerRepositorySql) DeleteSellerById(id int) error {
	sql_utils.LogAudit("DeleteSellerById", logger.LogStatusInProgress, "Delete seller by id: "+strconv.Itoa(id))

	affectedRows, err := sql_utils.Delete(r.db, "DELETE FROM sellers WHERE id = ?", []any{id})
	if err != nil {
		err = sql_utils.HandleSqlError(err)
		sql_utils.LogAuditError("DeleteSellerById", "Delete seller by id: "+strconv.Itoa(id), err)
		return err
	}
	if affectedRows == 0 {
		err = custom_errors.ErrNotFound
		sql_utils.LogAuditError("DeleteSellerById", "Delete seller by id: "+strconv.Itoa(id), err)
		return err
	}

	sql_utils.LogAudit("DeleteSellerById", logger.LogStatusSuccess, "Delete seller by id: "+strconv.Itoa(id))
	return nil
}

// UpdateSellerById performs a partial update (PATCH) on a seller record by ID.
//
// This method updates only the fields that are explicitly provided (non-nil pointers),
// allowing for flexible partial updates without affecting unchanged fields.
//
// Parameters:
//   - id: The unique identifier of the seller to update
//   - companyId: Optional pointer to new company ID. Pass nil to leave unchanged
//   - companyName: Optional pointer to new company name. Pass nil to leave unchanged
//   - address: Optional pointer to new address. Pass nil to leave unchanged
//   - telephone: Optional pointer to new telephone. Pass nil to leave unchanged
//
// At least one optional parameter must be provided (non-nil) for the update to proceed.
//
// Returns:
//   - models.Seller: The updated seller with all current field values
//   - error: Possible errors include:
//   - MandatoryArgMissingErr: When all optional parameters are nil
//   - ErrNotFound: When no seller exists with the given ID
//   - Database errors: Any SQL execution errors
func (r *SellerRepositorySql) UpdateSellerById(id int, companyId *int, companyName *string, address *string, telephone *string) (models.Seller, error) {
	sql_utils.LogAudit("UpdateSellerById", logger.LogStatusInProgress, "Update seller by id: "+strconv.Itoa(id))

	var columnsToSet []string
	var args []any

	if companyId != nil {
		columnsToSet = append(columnsToSet, "cid = ?")
		args = append(args, *companyId)
	}

	if companyName != nil {
		columnsToSet = append(columnsToSet, "company_name = ?")
		args = append(args, *companyName)
	}

	if address != nil {
		columnsToSet = append(columnsToSet, "address = ?")
		args = append(args, *address)
	}

	if telephone != nil {
		columnsToSet = append(columnsToSet, "telephone = ?")
		args = append(args, *telephone)
	}

	if len(columnsToSet) == 0 {
		err := &custom_errors.MandatoryArgMissingErr{Argument: "companyId or companyName or address or telephone"}
		sql_utils.LogAuditError("UpdateSellerById", "Update seller by id: "+strconv.Itoa(id), err)
		return models.Seller{}, err
	}

	query := fmt.Sprintf("UPDATE sellers SET %s WHERE id = ?", strings.Join(columnsToSet, ", "))
	args = append(args, id)

	_, err := sql_utils.Update(r.db, query, args)
	if err != nil {
		err = sql_utils.HandleSqlError(err)
		sql_utils.LogAuditError("UpdateSellerById", "Update seller by id: "+strconv.Itoa(id), err)
		return models.Seller{}, err
	}

	seller, err := r.GetSellerById(id) // Si bien es ineficiente porque hacemos 2 llamadas a la base de datos, devolvemos GetSellerById para cumplir con el requisito del sprint 1 de la respuesta del patch con la data del objeto patcheado.

	if err == nil {
		sql_utils.LogAudit("UpdateSellerById", logger.LogStatusSuccess, "Update seller by id: "+strconv.Itoa(id))
	}

	return seller, err
}
