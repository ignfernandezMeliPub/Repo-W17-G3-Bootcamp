package buyer_repository

import (
	"app/pkg/custom_errors"
	"app/pkg/models"
	"database/sql"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

func setupBuyerRepository(t *testing.T) (*BuyerRepositorySQL, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := NewBuyerSQL(db)

	cleanup := func() {
		db.Close()
	}

	return repo, mock, cleanup
}

func TestNewBuyerSQL(t *testing.T) {
	t.Run("should create repository successfully", func(t *testing.T) {
		db, _, err := sqlmock.New()
		require.NoError(t, err)
		defer db.Close()

		repo := NewBuyerSQL(db)

		require.NotNil(t, repo)
		require.Equal(t, db, repo.db)
	})

	t.Run("should return nil when db is nil", func(t *testing.T) {
		repo := NewBuyerSQL(nil)

		require.Nil(t, repo)
	})
}

func TestBuyerRepositorySQL_GetAllBuyers(t *testing.T) {
	t.Run("should get all buyers successfully", func(t *testing.T) {
		repo, mock, cleanup := setupBuyerRepository(t)
		defer cleanup()

		expectedBuyers := []models.Buyer{
			{Id: 1, CardNumberId: "12345", FirstName: "John", LastName: "Doe"},
			{Id: 2, CardNumberId: "67890", FirstName: "Jane", LastName: "Smith"},
		}

		rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name"}).
			AddRow(expectedBuyers[0].Id, expectedBuyers[0].CardNumberId, expectedBuyers[0].FirstName, expectedBuyers[0].LastName).
			AddRow(expectedBuyers[1].Id, expectedBuyers[1].CardNumberId, expectedBuyers[1].FirstName, expectedBuyers[1].LastName)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, card_number_id, first_name, last_name FROM buyers")).
			WillReturnRows(rows)

		result, err := repo.GetAllBuyers()

		require.NoError(t, err)
		require.Equal(t, expectedBuyers, result)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return not found error when no buyers found", func(t *testing.T) {
		repo, mock, cleanup := setupBuyerRepository(t)
		defer cleanup()

		rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name"})

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, card_number_id, first_name, last_name FROM buyers")).
			WillReturnRows(rows)

		_, err := repo.GetAllBuyers()

		require.Error(t, err)
		require.IsType(t, custom_errors.ErrNotFound, err)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when query fails", func(t *testing.T) {
		repo, mock, cleanup := setupBuyerRepository(t)
		defer cleanup()

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, card_number_id, first_name, last_name FROM buyers")).
			WillReturnError(sqlmock.ErrCancelled)

		_, err := repo.GetAllBuyers()

		require.Error(t, err)

		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestBuyerRepositorySQL_GetBuyerById(t *testing.T) {
	t.Run("should get buyer by id successfully", func(t *testing.T) {
		repo, mock, cleanup := setupBuyerRepository(t)
		defer cleanup()

		expectedBuyer := models.Buyer{
			Id:           1,
			CardNumberId: "12345",
			FirstName:    "John",
			LastName:     "Doe",
		}

		rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name"}).
			AddRow(expectedBuyer.Id, expectedBuyer.CardNumberId, expectedBuyer.FirstName, expectedBuyer.LastName)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, card_number_id, first_name, last_name FROM buyers WHERE id = ?")).
			WithArgs(1).
			WillReturnRows(rows)

		result, err := repo.GetBuyerById(1)

		require.NoError(t, err)
		require.Equal(t, expectedBuyer, result)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return not found error when buyer does not exist", func(t *testing.T) {
		repo, mock, cleanup := setupBuyerRepository(t)
		defer cleanup()

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, card_number_id, first_name, last_name FROM buyers WHERE id = ?")).
			WithArgs(999).
			WillReturnError(sql.ErrNoRows)

		_, err := repo.GetBuyerById(999)

		require.Error(t, err)
		require.IsType(t, custom_errors.ErrNotFound, err)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when query fails", func(t *testing.T) {
		repo, mock, cleanup := setupBuyerRepository(t)
		defer cleanup()

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, card_number_id, first_name, last_name FROM buyers WHERE id = ?")).
			WithArgs(1).
			WillReturnError(sqlmock.ErrCancelled)

		_, err := repo.GetBuyerById(1)

		require.Error(t, err)

		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestBuyerRepositorySQL_CreateBuyer(t *testing.T) {

	buyer := models.Buyer{
		CardNumberId: "12345",
		FirstName:    "John",
		LastName:     "Doe",
	}

	t.Run("should create buyer successfully", func(t *testing.T) {
		repo, mock, cleanup := setupBuyerRepository(t)
		defer cleanup()

		expectedId := int64(1)

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO buyers (card_number_id, first_name, last_name) VALUES (?, ?, ?)")).
			WithArgs(buyer.CardNumberId, buyer.FirstName, buyer.LastName).
			WillReturnResult(sqlmock.NewResult(expectedId, 1))

		result, err := repo.CreateBuyer(buyer)

		require.NoError(t, err)
		buyer.Id = int(expectedId)
		require.Equal(t, buyer, result)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when card_number_id is duplicated (unique constraint violation)", func(t *testing.T) {
		repo, mock, cleanup := setupBuyerRepository(t)
		defer cleanup()

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO buyers (card_number_id, first_name, last_name) VALUES (?, ?, ?)")).
			WithArgs(buyer.CardNumberId, buyer.FirstName, buyer.LastName).
			WillReturnError(&mysql.MySQLError{Number: 1062, Message: "Duplicate entry '12345' for key 'card_number_id'"})

		_, err := repo.CreateBuyer(buyer)

		require.Error(t, err)
		require.IsType(t, custom_errors.ErrUniqueAttributeViolationError, err)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when insert fails", func(t *testing.T) {
		repo, mock, cleanup := setupBuyerRepository(t)
		defer cleanup()

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO buyers (card_number_id, first_name, last_name) VALUES (?, ?, ?)")).
			WithArgs(buyer.CardNumberId, buyer.FirstName, buyer.LastName).
			WillReturnError(sqlmock.ErrCancelled)

		_, err := repo.CreateBuyer(buyer)

		require.Error(t, err)

		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestBuyerRepositorySQL_UpdateBuyer(t *testing.T) {
	buyer := models.Buyer{
		Id:           1,
		CardNumberId: "12345",
		FirstName:    "John",
		LastName:     "Doe",
	}

	t.Run("should update buyer successfully", func(t *testing.T) {
		repo, mock, cleanup := setupBuyerRepository(t)
		defer cleanup()

		mock.ExpectExec(regexp.QuoteMeta("UPDATE buyers SET card_number_id = ?, first_name = ?, last_name = ? WHERE id = ?")).
			WithArgs(buyer.CardNumberId, buyer.FirstName, buyer.LastName, buyer.Id).
			WillReturnResult(sqlmock.NewResult(0, 1))

		result, err := repo.UpdateBuyer(buyer)

		require.NoError(t, err)
		require.Equal(t, buyer, result)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when card_number_id is duplicated (unique constraint violation)", func(t *testing.T) {
		repo, mock, cleanup := setupBuyerRepository(t)
		defer cleanup()

		mock.ExpectExec(regexp.QuoteMeta("UPDATE buyers SET card_number_id = ?, first_name = ?, last_name = ? WHERE id = ?")).
			WithArgs(buyer.CardNumberId, buyer.FirstName, buyer.LastName, buyer.Id).
			WillReturnError(&mysql.MySQLError{Number: 1062, Message: "Duplicate entry '12345' for key 'card_number_id'"})

		_, err := repo.UpdateBuyer(buyer)

		require.Error(t, err)
		require.IsType(t, custom_errors.ErrUniqueAttributeViolationError, err)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when update fails", func(t *testing.T) {
		repo, mock, cleanup := setupBuyerRepository(t)
		defer cleanup()

		mock.ExpectExec(regexp.QuoteMeta("UPDATE buyers SET card_number_id = ?, first_name = ?, last_name = ? WHERE id = ?")).
			WithArgs(buyer.CardNumberId, buyer.FirstName, buyer.LastName, buyer.Id).
			WillReturnError(sqlmock.ErrCancelled)

		_, err := repo.UpdateBuyer(buyer)

		require.Error(t, err)

		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestBuyerRepositorySQL_DeleteBuyerById(t *testing.T) {
	t.Run("should delete buyer successfully", func(t *testing.T) {
		repo, mock, cleanup := setupBuyerRepository(t)
		defer cleanup()

		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM buyers WHERE id = ?")).
			WithArgs(1).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := repo.DeleteBuyerById(1)

		require.NoError(t, err)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return not found error when buyer does not exist", func(t *testing.T) {
		repo, mock, cleanup := setupBuyerRepository(t)
		defer cleanup()

		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM buyers WHERE id = ?")).
			WithArgs(999).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := repo.DeleteBuyerById(999)

		require.Error(t, err)
		require.IsType(t, custom_errors.ErrNotFound, err)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return foreign key violation error when buyer has dependencies", func(t *testing.T) {
		repo, mock, cleanup := setupBuyerRepository(t)
		defer cleanup()

		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM buyers WHERE id = ?")).
			WithArgs(1).
			WillReturnError(&mysql.MySQLError{Number: 1451, Message: "Cannot delete or update a parent row: a foreign key constraint fails (`test`.`purchase_orders`, CONSTRAINT `fk_purchase_orders_buyer_id` FOREIGN KEY (`buyer_id`) REFERENCES `buyers` (`id`))"})

		err := repo.DeleteBuyerById(1)

		require.Error(t, err)
		require.IsType(t, custom_errors.ErrForeignKeyViolation, err)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when delete fails", func(t *testing.T) {
		repo, mock, cleanup := setupBuyerRepository(t)
		defer cleanup()

		mock.ExpectExec(regexp.QuoteMeta("DELETE FROM buyers WHERE id = ?")).
			WithArgs(1).
			WillReturnError(sqlmock.ErrCancelled)

		err := repo.DeleteBuyerById(1)

		require.Error(t, err)

		require.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestBuyerRepositorySQL_GetBuyersPurchaseOrdersCount(t *testing.T) {
	t.Run("should get all buyers purchase orders count successfully", func(t *testing.T) {
		repo, mock, cleanup := setupBuyerRepository(t)
		defer cleanup()

		expectedBuyers := []models.BuyerPurchaseOrdersCount{
			{Id: 1, CardNumberId: "12345", FirstName: "John", LastName: "Doe", PurchaseOrdersCount: 3},
			{Id: 2, CardNumberId: "67890", FirstName: "Jane", LastName: "Smith", PurchaseOrdersCount: 1},
		}

		rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name", "purchase_orders_count"}).
			AddRow(expectedBuyers[0].Id, expectedBuyers[0].CardNumberId, expectedBuyers[0].FirstName, expectedBuyers[0].LastName, expectedBuyers[0].PurchaseOrdersCount).
			AddRow(expectedBuyers[1].Id, expectedBuyers[1].CardNumberId, expectedBuyers[1].FirstName, expectedBuyers[1].LastName, expectedBuyers[1].PurchaseOrdersCount)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT buyers.id as id, buyers.card_number_id, buyers.first_name, buyers.last_name, COUNT(purchase_orders.buyer_id) as purchase_orders_count FROM buyers LEFT JOIN purchase_orders ON buyers.id = purchase_orders.buyer_id GROUP BY buyers.id")).
			WillReturnRows(rows)

		result, err := repo.GetBuyersPurchaseOrdersCount(nil)

		require.NoError(t, err)
		require.Equal(t, expectedBuyers, result)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should get specific buyer purchase orders count successfully", func(t *testing.T) {
		repo, mock, cleanup := setupBuyerRepository(t)
		defer cleanup()

		buyerId := 1
		expectedBuyer := models.BuyerPurchaseOrdersCount{
			Id:                  1,
			CardNumberId:        "12345",
			FirstName:           "John",
			LastName:            "Doe",
			PurchaseOrdersCount: 3,
		}

		rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name", "purchase_orders_count"}).
			AddRow(expectedBuyer.Id, expectedBuyer.CardNumberId, expectedBuyer.FirstName, expectedBuyer.LastName, expectedBuyer.PurchaseOrdersCount)

		mock.ExpectQuery(regexp.QuoteMeta("SELECT buyers.id as id, buyers.card_number_id, buyers.first_name, buyers.last_name, COUNT(purchase_orders.buyer_id) as purchase_orders_count FROM buyers LEFT JOIN purchase_orders ON buyers.id = purchase_orders.buyer_id WHERE buyers.id = ? GROUP BY buyers.id")).
			WithArgs(buyerId).
			WillReturnRows(rows)

		result, err := repo.GetBuyersPurchaseOrdersCount(&buyerId)

		require.NoError(t, err)
		require.Len(t, result, 1)
		require.Equal(t, expectedBuyer, result[0])

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return not found error when no buyers found", func(t *testing.T) {
		repo, mock, cleanup := setupBuyerRepository(t)
		defer cleanup()

		rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name", "purchase_orders_count"})

		mock.ExpectQuery(regexp.QuoteMeta("SELECT buyers.id as id, buyers.card_number_id, buyers.first_name, buyers.last_name, COUNT(purchase_orders.buyer_id) as purchase_orders_count FROM buyers LEFT JOIN purchase_orders ON buyers.id = purchase_orders.buyer_id GROUP BY buyers.id")).
			WillReturnRows(rows)

		_, err := repo.GetBuyersPurchaseOrdersCount(nil)

		require.Error(t, err)
		require.IsType(t, custom_errors.ErrNotFound, err)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when query fails", func(t *testing.T) {
		repo, mock, cleanup := setupBuyerRepository(t)
		defer cleanup()

		mock.ExpectQuery(regexp.QuoteMeta("SELECT buyers.id as id, buyers.card_number_id, buyers.first_name, buyers.last_name, COUNT(purchase_orders.buyer_id) as purchase_orders_count FROM buyers LEFT JOIN purchase_orders ON buyers.id = purchase_orders.buyer_id GROUP BY buyers.id")).
			WillReturnError(sqlmock.ErrCancelled)

		_, err := repo.GetBuyersPurchaseOrdersCount(nil)

		require.Error(t, err)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return not found error when no buyers found with specific buyer id", func(t *testing.T) {
		repo, mock, cleanup := setupBuyerRepository(t)
		defer cleanup()

		buyerId := 1

		rows := sqlmock.NewRows([]string{"id", "card_number_id", "first_name", "last_name", "purchase_orders_count"})

		mock.ExpectQuery(regexp.QuoteMeta("SELECT buyers.id as id, buyers.card_number_id, buyers.first_name, buyers.last_name, COUNT(purchase_orders.buyer_id) as purchase_orders_count FROM buyers LEFT JOIN purchase_orders ON buyers.id = purchase_orders.buyer_id WHERE buyers.id = ? GROUP BY buyers.id")).
			WithArgs(1).
			WillReturnRows(rows)

		_, err := repo.GetBuyersPurchaseOrdersCount(&buyerId)

		require.Error(t, err)
		require.IsType(t, custom_errors.ErrNotFound, err)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when specific buyer query fails", func(t *testing.T) {
		repo, mock, cleanup := setupBuyerRepository(t)
		defer cleanup()

		buyerId := 1

		mock.ExpectQuery(regexp.QuoteMeta("SELECT buyers.id as id, buyers.card_number_id, buyers.first_name, buyers.last_name, COUNT(purchase_orders.buyer_id) as purchase_orders_count FROM buyers LEFT JOIN purchase_orders ON buyers.id = purchase_orders.buyer_id WHERE buyers.id = ? GROUP BY buyers.id")).
			WithArgs(buyerId).
			WillReturnError(sqlmock.ErrCancelled)

		_, err := repo.GetBuyersPurchaseOrdersCount(&buyerId)

		require.Error(t, err)

		require.NoError(t, mock.ExpectationsWereMet())
	})
}
