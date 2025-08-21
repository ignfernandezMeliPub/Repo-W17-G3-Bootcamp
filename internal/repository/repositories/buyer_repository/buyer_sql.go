package buyer_repository

import (
	"app/internal/logger"
	"app/internal/repository/sql_utils"
	"app/pkg/custom_errors"
	"app/pkg/models"
	"database/sql"
	"strconv"
)

type BuyerRepositorySQL struct {
	db *sql.DB
}

func NewBuyerSQL(db *sql.DB) *BuyerRepositorySQL {
	if db == nil {
		return nil
	}

	return &BuyerRepositorySQL{db: db}
}

func (r *BuyerRepositorySQL) GetAllBuyers() (b []models.Buyer, err error) {
	sql_utils.Log("GetAllBuyers", logger.LogStatusInProgress, "Select all buyers")
	b, err = sql_utils.Query[models.Buyer](r.db, "SELECT id, card_number_id, first_name, last_name FROM buyers", []any{})

	if err != nil {
		err = sql_utils.HandleSqlError(err)
		sql_utils.LogError("GetAllBuyers", "Select all buyers", err)
		return b, err
	}

	sql_utils.Log("GetAllBuyers", logger.LogStatusSuccess, "Select all buyers")
	return b, nil
}

func (r *BuyerRepositorySQL) GetBuyerById(id int) (b models.Buyer, err error) {
	sql_utils.Log("GetBuyerById", logger.LogStatusInProgress, "Select buyer by id "+strconv.Itoa(id))

	b, err = sql_utils.QueryRow[models.Buyer](r.db, "SELECT id, card_number_id, first_name, last_name FROM buyers WHERE id = ?", []any{id})

	if err != nil {
		err = sql_utils.HandleSqlError(err)
		sql_utils.LogError("GetBuyerById", "Select buyer by id "+strconv.Itoa(id), err)
		return b, err
	}

	sql_utils.Log("GetBuyerById", logger.LogStatusSuccess, "Select buyer by id "+strconv.Itoa(id))
	return b, nil
}

func (r *BuyerRepositorySQL) CreateBuyer(buyer models.Buyer) (newBuyer models.Buyer, err error) {
	sql_utils.LogAudit("CreateBuyer", logger.LogStatusInProgress, "Insert buyer")

	newId, err := sql_utils.Insert(r.db, "INSERT INTO buyers (card_number_id, first_name, last_name) VALUES (?, ?, ?)", []any{buyer.CardNumberId, buyer.FirstName, buyer.LastName})

	if err != nil {
		err = sql_utils.HandleSqlError(err)
		sql_utils.LogAuditError("CreateBuyer", "Insert buyer", err)
		return newBuyer, err
	}

	newBuyer = buyer
	newBuyer.Id = int(newId)

	sql_utils.LogAudit("CreateBuyer", logger.LogStatusSuccess, "Insert buyer. Id: "+strconv.Itoa(newBuyer.Id))
	return newBuyer, nil
}
func (r *BuyerRepositorySQL) UpdateBuyer(buyer models.Buyer) (newBuyer models.Buyer, err error) {
	sql_utils.LogAudit("UpdateBuyer", logger.LogStatusInProgress, "Update buyer by id: "+strconv.Itoa(buyer.Id))

	_, err = sql_utils.Update(r.db, "UPDATE buyers SET card_number_id = ?, first_name = ?, last_name = ? WHERE id = ?", []any{buyer.CardNumberId, buyer.FirstName, buyer.LastName, buyer.Id})

	if err != nil {
		err = sql_utils.HandleSqlError(err)
		sql_utils.LogAuditError("UpdateBuyer", "Update buyer by id: "+strconv.Itoa(buyer.Id), err)
		return newBuyer, err
	}

	newBuyer = buyer
	sql_utils.LogAudit("UpdateBuyer", logger.LogStatusSuccess, "Update buyer by id: "+strconv.Itoa(newBuyer.Id))
	return newBuyer, nil
}

func (r *BuyerRepositorySQL) DeleteBuyerById(id int) (err error) {
	sql_utils.LogAudit("DeleteBuyerById", logger.LogStatusInProgress, "Delete buyer by id: "+strconv.Itoa(id))

	rowsAffected, err := sql_utils.Delete(r.db, "DELETE FROM buyers WHERE id = ?", []any{id})
	if err != nil {
		err = sql_utils.HandleSqlError(err)
		sql_utils.LogAuditError("DeleteBuyerById", "Delete buyer by id: "+strconv.Itoa(id), err)
		return err
	}

	if rowsAffected == 0 {
		sql_utils.LogAuditError("DeleteBuyerById", "Delete buyer by id: "+strconv.Itoa(id), custom_errors.ErrNotFound)
		return custom_errors.ErrNotFound
	}

	sql_utils.LogAudit("DeleteBuyerById", logger.LogStatusSuccess, "Delete buyer by id: "+strconv.Itoa(id))
	return nil
}

func (r *BuyerRepositorySQL) GetBuyersPurchaseOrdersCount(buyerId *int) (b []models.BuyerPurchaseOrdersCount, err error) {
	if buyerId != nil {
		sql_utils.LogAudit("GetBuyersPurchaseOrdersCount", logger.LogStatusInProgress, "Select buyers purchase orders count by id: "+strconv.Itoa(*buyerId))
	} else {
		sql_utils.LogAudit("GetBuyersPurchaseOrdersCount", logger.LogStatusInProgress, "Select buyers purchase orders count")
	}

	query := "SELECT buyers.id as id, buyers.card_number_id, buyers.first_name, buyers.last_name, COUNT(purchase_orders.buyer_id) as purchase_orders_count FROM buyers LEFT JOIN purchase_orders ON buyers.id = purchase_orders.buyer_id"

	if buyerId != nil {
		query += " WHERE buyers.id = ? GROUP BY buyers.id"
		b, err = sql_utils.Query[models.BuyerPurchaseOrdersCount](r.db, query, []any{*buyerId})
	} else {
		query += " GROUP BY buyers.id"
		b, err = sql_utils.Query[models.BuyerPurchaseOrdersCount](r.db, query, []any{})
	}

	if err != nil {
		err = sql_utils.HandleSqlError(err)
		sql_utils.LogAuditError("GetBuyersPurchaseOrdersCount", "Select buyers purchase orders count", err)
		return b, err
	}

	if buyerId != nil {
		sql_utils.LogAudit("GetBuyersPurchaseOrdersCount", logger.LogStatusSuccess, "Select buyers purchase orders count by id: "+strconv.Itoa(*buyerId))
	} else {
		sql_utils.LogAudit("GetBuyersPurchaseOrdersCount", logger.LogStatusSuccess, "Select buyers purchase orders count")
	}
	return b, nil
}
