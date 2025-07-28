package purchase_order_repository

import (
	"app/pkg/custom_errors"
	"app/pkg/models"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
)

func setupPurchaseOrderRepository(t *testing.T) (*PurchaseOrderRepositorySQL, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	repo := NewPurchaseOrderRepositorySQL(db)

	cleanup := func() {
		db.Close()
	}

	return repo, mock, cleanup
}

func TestPurchaseOrderRepositorySQL_CreatePurchaseOrder(t *testing.T) {
	t.Run("should create purchase order successfully", func(t *testing.T) {
		repo, mock, cleanup := setupPurchaseOrderRepository(t)
		defer cleanup()

		orderDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		purchaseOrder := models.PurchaseOrder{
			OrderNumber:  "PO-002",
			OrderDate:    orderDate,
			TrackingCode: "TRACK-002",
			BuyerId:      2,
			PurchaseOrderDetails: []models.PurchaseOrderDetail{
				{ProductRecordId: 1, Quantity: 10},
				{ProductRecordId: 2, Quantity: 5},
			},
		}

		expectedOrderId := int64(2)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO purchase_orders (order_number, order_date, tracking_code, buyer_id) VALUES (?, ?, ?, ?)")).
			WithArgs(purchaseOrder.OrderNumber, purchaseOrder.OrderDate, purchaseOrder.TrackingCode, purchaseOrder.BuyerId).
			WillReturnResult(sqlmock.NewResult(expectedOrderId, 1))

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO purchase_order_details (order_id, product_record_id, quantity) VALUES (?, ?, ?), (?, ?, ?)")).
			WithArgs(expectedOrderId, purchaseOrder.PurchaseOrderDetails[0].ProductRecordId, purchaseOrder.PurchaseOrderDetails[0].Quantity, expectedOrderId, purchaseOrder.PurchaseOrderDetails[1].ProductRecordId, purchaseOrder.PurchaseOrderDetails[1].Quantity).
			WillReturnResult(sqlmock.NewResult(1, 2))

		rows := sqlmock.NewRows([]string{"id", "order_id", "product_record_id", "quantity"}).
			AddRow(1, expectedOrderId, 1, 10).
			AddRow(2, expectedOrderId, 2, 5)
		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, order_id, product_record_id, quantity FROM purchase_order_details WHERE order_id = ?")).
			WithArgs(expectedOrderId).
			WillReturnRows(rows)

		mock.ExpectCommit()

		result, err := repo.CreatePurchaseOrder(purchaseOrder)

		require.NoError(t, err)
		require.Equal(t, int(expectedOrderId), result.Id)
		require.Equal(t, purchaseOrder.OrderNumber, result.OrderNumber)
		require.True(t, result.OrderDate.Equal(purchaseOrder.OrderDate))
		require.Equal(t, purchaseOrder.TrackingCode, result.TrackingCode)
		require.Equal(t, purchaseOrder.BuyerId, result.BuyerId)
		require.Len(t, result.PurchaseOrderDetails, 2)

		require.Equal(t, 1, result.PurchaseOrderDetails[0].Id)
		require.Equal(t, purchaseOrder.PurchaseOrderDetails[0].ProductRecordId, result.PurchaseOrderDetails[0].ProductRecordId)
		require.Equal(t, purchaseOrder.PurchaseOrderDetails[0].Quantity, result.PurchaseOrderDetails[0].Quantity)

		require.Equal(t, 2, result.PurchaseOrderDetails[1].Id)
		require.Equal(t, purchaseOrder.PurchaseOrderDetails[1].ProductRecordId, result.PurchaseOrderDetails[1].ProductRecordId)
		require.Equal(t, purchaseOrder.PurchaseOrderDetails[1].Quantity, result.PurchaseOrderDetails[1].Quantity)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should rollback transaction on purchase order insert error", func(t *testing.T) {
		repo, mock, cleanup := setupPurchaseOrderRepository(t)
		defer cleanup()

		orderDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		purchaseOrder := models.PurchaseOrder{
			OrderNumber:          "PO-003",
			OrderDate:            orderDate,
			TrackingCode:         "TRACK-003",
			BuyerId:              3,
			PurchaseOrderDetails: []models.PurchaseOrderDetail{},
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO purchase_orders (order_number, order_date, tracking_code, buyer_id) VALUES (?, ?, ?, ?)")).
			WithArgs(purchaseOrder.OrderNumber, purchaseOrder.OrderDate, purchaseOrder.TrackingCode, purchaseOrder.BuyerId).
			WillReturnError(sqlmock.ErrCancelled)
		mock.ExpectRollback()

		_, err := repo.CreatePurchaseOrder(purchaseOrder)

		require.Error(t, err)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should rollback transaction on purchase order details insert error", func(t *testing.T) {
		repo, mock, cleanup := setupPurchaseOrderRepository(t)
		defer cleanup()

		orderDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		purchaseOrder := models.PurchaseOrder{
			OrderNumber:  "PO-004",
			OrderDate:    orderDate,
			TrackingCode: "TRACK-004",
			BuyerId:      4,
			PurchaseOrderDetails: []models.PurchaseOrderDetail{
				{ProductRecordId: 1, Quantity: 10},
			},
		}

		expectedOrderId := int64(4)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO purchase_orders (order_number, order_date, tracking_code, buyer_id) VALUES (?, ?, ?, ?)")).
			WithArgs(purchaseOrder.OrderNumber, purchaseOrder.OrderDate, purchaseOrder.TrackingCode, purchaseOrder.BuyerId).
			WillReturnResult(sqlmock.NewResult(expectedOrderId, 1))

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO purchase_order_details (order_id, product_record_id, quantity) VALUES (?, ?, ?)")).
			WithArgs(expectedOrderId, purchaseOrder.PurchaseOrderDetails[0].ProductRecordId, purchaseOrder.PurchaseOrderDetails[0].Quantity).
			WillReturnError(sqlmock.ErrCancelled)
		mock.ExpectRollback()

		_, err := repo.CreatePurchaseOrder(purchaseOrder)

		require.Error(t, err)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when begin transaction fails", func(t *testing.T) {
		repo, mock, cleanup := setupPurchaseOrderRepository(t)
		defer cleanup()

		orderDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		purchaseOrder := models.PurchaseOrder{
			OrderNumber:          "PO-005",
			OrderDate:            orderDate,
			TrackingCode:         "TRACK-005",
			BuyerId:              5,
			PurchaseOrderDetails: []models.PurchaseOrderDetail{},
		}

		mock.ExpectBegin().WillReturnError(sqlmock.ErrCancelled)

		_, err := repo.CreatePurchaseOrder(purchaseOrder)

		require.Error(t, err)

		require.NoError(t, mock.ExpectationsWereMet())
	})
	t.Run("should return error when buyerId does not exist (foreign key violation)", func(t *testing.T) {
		repo, mock, cleanup := setupPurchaseOrderRepository(t)
		defer cleanup()

		orderDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		purchaseOrder := models.PurchaseOrder{
			OrderNumber:  "PO-008",
			OrderDate:    orderDate,
			TrackingCode: "TRACK-008",
			BuyerId:      9999, // Non-existent buyerId
			PurchaseOrderDetails: []models.PurchaseOrderDetail{
				{ProductRecordId: 1, Quantity: 10},
			},
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO purchase_orders (order_number, order_date, tracking_code, buyer_id) VALUES (?, ?, ?, ?)")).
			WithArgs(purchaseOrder.OrderNumber, purchaseOrder.OrderDate, purchaseOrder.TrackingCode, purchaseOrder.BuyerId).
			WillReturnError(&mysql.MySQLError{Number: 1451, Message: "FOREIGN KEY \\(`buyer_id`\\)"})

		mock.ExpectRollback()

		_, err := repo.CreatePurchaseOrder(purchaseOrder)

		require.Error(t, err)
		errExpected := custom_errors.ErrForeignKeyViolation
		require.NotNil(t, err)
		require.IsType(t, errExpected, err)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should return error when order number is duplicated (unique constraint violation)", func(t *testing.T) {
		repo, mock, cleanup := setupPurchaseOrderRepository(t)
		defer cleanup()

		orderDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		purchaseOrder := models.PurchaseOrder{
			OrderNumber:  "PO-001", // Assume this order number already exists
			OrderDate:    orderDate,
			TrackingCode: "TRACK-009",
			BuyerId:      1,
			PurchaseOrderDetails: []models.PurchaseOrderDetail{
				{ProductRecordId: 1, Quantity: 10},
			},
		}

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO purchase_orders (order_number, order_date, tracking_code, buyer_id) VALUES (?, ?, ?, ?)")).
			WithArgs(purchaseOrder.OrderNumber, purchaseOrder.OrderDate, purchaseOrder.TrackingCode, purchaseOrder.BuyerId).
			WillReturnError(&mysql.MySQLError{Number: 1062, Message: "Duplicate entry 'PO-001' for key 'order_number'"})
		mock.ExpectRollback()

		_, err := repo.CreatePurchaseOrder(purchaseOrder)

		require.Error(t, err)
		errExpected := custom_errors.ErrUniqueAttributeViolationError
		require.NotNil(t, err)
		require.IsType(t, errExpected, err)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should rollback transaction on commit error", func(t *testing.T) {
		repo, mock, cleanup := setupPurchaseOrderRepository(t)
		defer cleanup()

		orderDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		purchaseOrder := models.PurchaseOrder{
			OrderNumber:          "PO-006",
			OrderDate:            orderDate,
			TrackingCode:         "TRACK-006",
			BuyerId:              6,
			PurchaseOrderDetails: []models.PurchaseOrderDetail{},
		}

		expectedOrderId := int64(6)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO purchase_orders (order_number, order_date, tracking_code, buyer_id) VALUES (?, ?, ?, ?)")).
			WithArgs(purchaseOrder.OrderNumber, purchaseOrder.OrderDate, purchaseOrder.TrackingCode, purchaseOrder.BuyerId).
			WillReturnResult(sqlmock.NewResult(expectedOrderId, 1))
		mock.ExpectCommit().WillReturnError(sqlmock.ErrCancelled)

		_, err := repo.CreatePurchaseOrder(purchaseOrder)

		require.Error(t, err)

		require.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("should rollback transaction on purchase order details query error", func(t *testing.T) {
		repo, mock, cleanup := setupPurchaseOrderRepository(t)
		defer cleanup()

		orderDate := time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC)
		purchaseOrder := models.PurchaseOrder{
			OrderNumber:  "PO-007",
			OrderDate:    orderDate,
			TrackingCode: "TRACK-007",
			BuyerId:      7,
			PurchaseOrderDetails: []models.PurchaseOrderDetail{
				{ProductRecordId: 1, Quantity: 10},
			},
		}

		expectedOrderId := int64(7)

		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO purchase_orders (order_number, order_date, tracking_code, buyer_id) VALUES (?, ?, ?, ?)")).
			WithArgs(purchaseOrder.OrderNumber, purchaseOrder.OrderDate, purchaseOrder.TrackingCode, purchaseOrder.BuyerId).
			WillReturnResult(sqlmock.NewResult(expectedOrderId, 1))

		mock.ExpectExec(regexp.QuoteMeta("INSERT INTO purchase_order_details (order_id, product_record_id, quantity) VALUES (?, ?, ?)")).
			WithArgs(expectedOrderId, purchaseOrder.PurchaseOrderDetails[0].ProductRecordId, purchaseOrder.PurchaseOrderDetails[0].Quantity).
			WillReturnResult(sqlmock.NewResult(1, 1))

		mock.ExpectQuery(regexp.QuoteMeta("SELECT id, order_id, product_record_id, quantity FROM purchase_order_details WHERE order_id = ?")).
			WithArgs(expectedOrderId).
			WillReturnError(sqlmock.ErrCancelled)
		mock.ExpectRollback()

		_, err := repo.CreatePurchaseOrder(purchaseOrder)

		require.Error(t, err)

		require.NoError(t, mock.ExpectationsWereMet())
	})
}
