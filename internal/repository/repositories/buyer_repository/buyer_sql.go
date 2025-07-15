package buyer_repository

import (
	"app/internal/repository/sql_utils"
	"app/pkg/custom_errors"
	"app/pkg/models"
	"database/sql"
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
	b, err = sql_utils.Query[models.Buyer](r.db, "SELECT id, card_number_id, first_name, last_name FROM buyers", []any{})
	return b, sql_utils.HandleSqlError(err)
}

func (r *BuyerRepositorySQL) GetBuyerById(id int) (b models.Buyer, err error) {
	b, err = sql_utils.QueryRow[models.Buyer](r.db, "SELECT id, card_number_id, first_name, last_name FROM buyers WHERE id = ?", []any{id})
	return b, sql_utils.HandleSqlError(err)
}

func (r *BuyerRepositorySQL) GetBuyerByCardNumberId(cardNumberId string) (b models.Buyer, err error) {
	b, err = sql_utils.QueryRow[models.Buyer](r.db, "SELECT id, card_number_id, first_name, last_name FROM buyers WHERE card_number_id = ?", []any{cardNumberId})
	return b, sql_utils.HandleSqlError(err)
}

func (r *BuyerRepositorySQL) CreateBuyer(buyer models.Buyer) (newBuyer models.Buyer, err error) {
	newId, err := sql_utils.Insert(r.db, "INSERT INTO buyers (card_number_id, first_name, last_name) VALUES (?, ?, ?)", []any{buyer.CardNumberId, buyer.FirstName, buyer.LastName})
	if err != nil {
		return newBuyer, sql_utils.HandleSqlError(err)
	}

	newBuyer = buyer
	newBuyer.Id = int(newId)
	return newBuyer, nil
}
func (r *BuyerRepositorySQL) UpdateBuyer(buyer models.Buyer) (newBuyer models.Buyer, err error) {
	_, err = sql_utils.Update(r.db, "UPDATE buyers SET card_number_id = ?, first_name = ?, last_name = ? WHERE id = ?", []any{buyer.CardNumberId, buyer.FirstName, buyer.LastName, buyer.Id})

	if err != nil {
		return newBuyer, sql_utils.HandleSqlError(err)
	}

	newBuyer = buyer
	return newBuyer, nil
}

func (r *BuyerRepositorySQL) DeleteBuyerById(id int) (err error) {
	rowsAffected, err := sql_utils.Delete(r.db, "DELETE FROM buyers WHERE id = ?", []any{id})
	if err != nil {
		return sql_utils.HandleSqlError(err)
	}

	if rowsAffected == 0 {
		return custom_errors.ErrNotFound
	}

	return
}

func (r *BuyerRepositorySQL) GetBuyerPurchaseOrdersCount(buyerId int) (b []models.BuyerPurchaseOrdersCount, err error) {
	b, err = sql_utils.Query[models.BuyerPurchaseOrdersCount](r.db, "SELECT buyers.id as id, buyers.card_number_id, buyers.first_name, buyers.last_name, COUNT(purchase_orders.buyer_id) as purchase_orders_count FROM buyers LEFT JOIN purchase_orders ON buyers.id = purchase_orders.buyer_id WHERE buyers.id = ? GROUP BY buyers.id", []any{buyerId})
	return b, sql_utils.HandleSqlError(err)
}

func (r *BuyerRepositorySQL) GetBuyersPurchaseOrdersCount() (b []models.BuyerPurchaseOrdersCount, err error) {
	b, err = sql_utils.Query[models.BuyerPurchaseOrdersCount](r.db, "SELECT buyers.id as id, buyers.card_number_id, buyers.first_name, buyers.last_name, COUNT(purchase_orders.buyer_id) as purchase_orders_count FROM buyers LEFT JOIN purchase_orders ON buyers.id = purchase_orders.buyer_id GROUP BY buyers.id", []any{})
	return b, sql_utils.HandleSqlError(err)
}
