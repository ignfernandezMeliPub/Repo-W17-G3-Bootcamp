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
	if err != nil {
		err = sql_utils.HandleSqlError(err)
	}
	return
}

func (r *BuyerRepositorySQL) GetBuyerById(id int) (b models.Buyer, err error) {
	b, err = sql_utils.QueryRow[models.Buyer](r.db, "SELECT id, card_number_id, first_name, last_name FROM buyers WHERE id = ?", []any{id})
	if err != nil {
		err = sql_utils.HandleSqlError(err)
	}
	return
}

func (r *BuyerRepositorySQL) GetBuyerByCardNumberId(cardNumberId string) (b models.Buyer, err error) {
	b, err = sql_utils.QueryRow[models.Buyer](r.db, "SELECT id, card_number_id, first_name, last_name FROM buyers WHERE card_number_id = ?", []any{cardNumberId})
	if err != nil {
		err = sql_utils.HandleSqlError(err)
	}
	return
}

func (r *BuyerRepositorySQL) CreateBuyer(_b models.Buyer) (b models.Buyer, err error) {
	newId, err := sql_utils.Insert(r.db, "INSERT INTO buyers (card_number_id, first_name, last_name) VALUES (?, ?, ?)", []any{_b.CardNumberId, _b.FirstName, _b.LastName})
	if err != nil {
		err = sql_utils.HandleSqlError(err)
		return
	}

	b = _b
	b.Id = int(newId)
	return
}
func (r *BuyerRepositorySQL) UpdateBuyer(_b models.Buyer) (b models.Buyer, err error) {
	_, err = sql_utils.Update(r.db, "UPDATE buyers SET card_number_id = ?, first_name = ?, last_name = ? WHERE id = ?", []any{_b.CardNumberId, _b.FirstName, _b.LastName, _b.Id})
	if err != nil {
		err = sql_utils.HandleSqlError(err)
		return
	}

	b = _b
	return
}

func (r *BuyerRepositorySQL) DeleteBuyerById(id int) (err error) {
	affectedRows, err := sql_utils.Delete(r.db, "DELETE FROM buyers WHERE id = ?", []any{id})
	if err != nil {
		err = sql_utils.HandleSqlError(err)
		return
	}

	if affectedRows == 0 {
		return custom_errors.ErrNotFound
	}

	return
}

func (r *BuyerRepositorySQL) GetBuyerPurchaseOrdersCount(buyerId int) (p []models.BuyerPurchaseOrdersCount, err error) {
	p, err = sql_utils.Query[models.BuyerPurchaseOrdersCount](r.db, "SELECT buyers.id as id, buyers.card_number_id, buyers.first_name, buyers.last_name, COUNT(*) as purchase_orders_count FROM buyers LEFT JOIN purchase_orders ON buyers.id = purchase_orders.buyer_id WHERE buyers.id = ? GROUP BY buyers.id", []any{buyerId})
	if err != nil {
		err = sql_utils.HandleSqlError(err)
	}
	return
}

func (r *BuyerRepositorySQL) GetBuyersPurchaseOrdersCount() (p []models.BuyerPurchaseOrdersCount, err error) {
	p, err = sql_utils.Query[models.BuyerPurchaseOrdersCount](r.db, "SELECT buyers.id as id, buyers.card_number_id, buyers.first_name, buyers.last_name, COUNT(*) as purchase_orders_count FROM buyers LEFT JOIN purchase_orders ON buyers.id = purchase_orders.buyer_id GROUP BY buyers.id", []any{})
	if err != nil {
		err = sql_utils.HandleSqlError(err)
	}
	return
}
