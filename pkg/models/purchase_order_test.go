package models

import (
	"app/pkg/custom_errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var purchaseOrderEmptyStringWithSpaces = "  "

var purchaseOrder = PurchaseOrder{
	OrderNumber:          "1234567890",
	OrderDate:            time.Now(),
	TrackingCode:         "1234567890",
	BuyerId:              1,
	PurchaseOrderDetails: []PurchaseOrderDetail{{ProductRecordId: 1, Quantity: 1}},
}

var validOrderDate = "2024-01-15"
var futureOrderDate = "9999-01-15"
var invalidOrderDate = "invalid-date"
var validProductRecordId = 1
var validQuantity = 5
var invalidProductRecordId = 0
var invalidQuantity = 0

func TestPurchaseOrderCreateRequest_Verify(t *testing.T) {
	emptyFields := []struct {
		name  string
		patch PurchaseOrderCreateRequest
	}{
		{
			name: "order_number",
			patch: PurchaseOrderCreateRequest{
				Data: &PurchaseOrderData{
					OrderNumber:  &purchaseOrderEmptyStringWithSpaces,
					OrderDate:    &validOrderDate,
					TrackingCode: &purchaseOrder.TrackingCode,
					BuyerId:      &purchaseOrder.BuyerId,
					PurchaseOrderDetails: []PurchaseOrderDetailRequest{
						{ProductRecordId: &validProductRecordId, Quantity: &validQuantity},
					},
				},
			},
		},
		{
			name: "tracking_code",
			patch: PurchaseOrderCreateRequest{
				Data: &PurchaseOrderData{
					OrderNumber:  &purchaseOrder.OrderNumber,
					OrderDate:    &validOrderDate,
					TrackingCode: &purchaseOrderEmptyStringWithSpaces,
					BuyerId:      &purchaseOrder.BuyerId,
					PurchaseOrderDetails: []PurchaseOrderDetailRequest{
						{ProductRecordId: &validProductRecordId, Quantity: &validQuantity},
					},
				},
			},
		},
	}

	for _, field := range emptyFields {
		t.Run("should return error if "+field.name+" is empty", func(t *testing.T) {
			err := field.patch.Verify()
			require.Error(t, err)
			require.Equal(t, &custom_errors.InvalidArgValueErr{
				Argument:  field.name,
				Value:     "",
				ExtraInfo: "Value must be non-empty",
			}, err)
		})
	}

	nilFields := []struct {
		name  string
		patch PurchaseOrderCreateRequest
	}{
		{
			name: "data",
			patch: PurchaseOrderCreateRequest{
				Data: nil,
			},
		},
		{
			name: "order_number",
			patch: PurchaseOrderCreateRequest{
				Data: &PurchaseOrderData{
					OrderNumber:  nil,
					OrderDate:    &validOrderDate,
					TrackingCode: &purchaseOrder.TrackingCode,
					BuyerId:      &purchaseOrder.BuyerId,
					PurchaseOrderDetails: []PurchaseOrderDetailRequest{
						{ProductRecordId: &validProductRecordId, Quantity: &validQuantity},
					},
				},
			},
		},
		{
			name: "order_date",
			patch: PurchaseOrderCreateRequest{
				Data: &PurchaseOrderData{
					OrderNumber:  &purchaseOrder.OrderNumber,
					OrderDate:    nil,
					TrackingCode: &purchaseOrder.TrackingCode,
					BuyerId:      &purchaseOrder.BuyerId,
					PurchaseOrderDetails: []PurchaseOrderDetailRequest{
						{ProductRecordId: &validProductRecordId, Quantity: &validQuantity},
					},
				},
			},
		},
		{
			name: "tracking_code",
			patch: PurchaseOrderCreateRequest{
				Data: &PurchaseOrderData{
					OrderNumber:  &purchaseOrder.OrderNumber,
					OrderDate:    &validOrderDate,
					TrackingCode: nil,
					BuyerId:      &purchaseOrder.BuyerId,
					PurchaseOrderDetails: []PurchaseOrderDetailRequest{
						{ProductRecordId: &validProductRecordId, Quantity: &validQuantity},
					},
				},
			},
		},
		{
			name: "buyer_id",
			patch: PurchaseOrderCreateRequest{
				Data: &PurchaseOrderData{
					OrderNumber:  &purchaseOrder.OrderNumber,
					OrderDate:    &validOrderDate,
					TrackingCode: &purchaseOrder.TrackingCode,
					BuyerId:      nil,
					PurchaseOrderDetails: []PurchaseOrderDetailRequest{
						{ProductRecordId: &validProductRecordId, Quantity: &validQuantity},
					},
				},
			},
		},
		{
			name: "product_record_id",
			patch: PurchaseOrderCreateRequest{
				Data: &PurchaseOrderData{
					OrderNumber:  &purchaseOrder.OrderNumber,
					OrderDate:    &validOrderDate,
					TrackingCode: &purchaseOrder.TrackingCode,
					BuyerId:      &purchaseOrder.BuyerId,
					PurchaseOrderDetails: []PurchaseOrderDetailRequest{
						{ProductRecordId: nil, Quantity: &validQuantity},
					},
				},
			},
		},
		{
			name: "quantity",
			patch: PurchaseOrderCreateRequest{
				Data: &PurchaseOrderData{
					OrderNumber:  &purchaseOrder.OrderNumber,
					OrderDate:    &validOrderDate,
					TrackingCode: &purchaseOrder.TrackingCode,
					BuyerId:      &purchaseOrder.BuyerId,
					PurchaseOrderDetails: []PurchaseOrderDetailRequest{
						{ProductRecordId: &validProductRecordId, Quantity: nil},
					},
				},
			},
		},
	}

	for _, field := range nilFields {
		t.Run("should return error if "+field.name+" is nil", func(t *testing.T) {
			err := field.patch.Verify()
			require.Error(t, err)
			require.Equal(t, &custom_errors.MandatoryArgMissingErr{
				Argument: field.name,
			}, err)
		})
	}

	t.Run("should return error if order_date is in the future", func(t *testing.T) {
		patch := PurchaseOrderCreateRequest{
			Data: &PurchaseOrderData{
				OrderNumber:  &purchaseOrder.OrderNumber,
				OrderDate:    &futureOrderDate,
				TrackingCode: &purchaseOrder.TrackingCode,
				BuyerId:      &purchaseOrder.BuyerId,
				PurchaseOrderDetails: []PurchaseOrderDetailRequest{
					{ProductRecordId: &validProductRecordId, Quantity: &validQuantity},
				},
			},
		}

		err := patch.Verify()
		require.Error(t, err)
		require.Equal(t, &custom_errors.InvalidArgValueErr{
			Argument:  "order_date",
			Value:     futureOrderDate,
			ExtraInfo: "Future dates are not allowed.",
		}, err)
	})
	t.Run("should return error if order_date has invalid format", func(t *testing.T) {
		patch := PurchaseOrderCreateRequest{
			Data: &PurchaseOrderData{
				OrderNumber:  &purchaseOrder.OrderNumber,
				OrderDate:    &invalidOrderDate,
				TrackingCode: &purchaseOrder.TrackingCode,
				BuyerId:      &purchaseOrder.BuyerId,
				PurchaseOrderDetails: []PurchaseOrderDetailRequest{
					{ProductRecordId: &validProductRecordId, Quantity: &validQuantity},
				},
			},
		}

		err := patch.Verify()
		require.Error(t, err)
		require.Equal(t, &custom_errors.InvalidArgValueErr{
			Argument:  "order_date",
			Value:     invalidOrderDate,
			ExtraInfo: "Value must be in the format YYYY-MM-DD",
		}, err)
	})

	t.Run("should return error if purchase_order_details is empty", func(t *testing.T) {
		patch := PurchaseOrderCreateRequest{
			Data: &PurchaseOrderData{
				OrderNumber:          &purchaseOrder.OrderNumber,
				OrderDate:            &validOrderDate,
				TrackingCode:         &purchaseOrder.TrackingCode,
				BuyerId:              &purchaseOrder.BuyerId,
				PurchaseOrderDetails: []PurchaseOrderDetailRequest{},
			},
		}

		err := patch.Verify()
		require.Error(t, err)
		require.Equal(t, &custom_errors.MandatoryArgMissingErr{
			Argument: "purchase_order_details",
		}, err)
	})

	invalidDetailFields := []struct {
		name  string
		patch PurchaseOrderCreateRequest
		value interface{}
	}{
		{
			name: "product_record_id",
			patch: PurchaseOrderCreateRequest{
				Data: &PurchaseOrderData{
					OrderNumber:  &purchaseOrder.OrderNumber,
					OrderDate:    &validOrderDate,
					TrackingCode: &purchaseOrder.TrackingCode,
					BuyerId:      &purchaseOrder.BuyerId,
					PurchaseOrderDetails: []PurchaseOrderDetailRequest{
						{ProductRecordId: &invalidProductRecordId, Quantity: &validQuantity},
					},
				},
			},
			value: &invalidProductRecordId,
		},
		{
			name: "quantity",
			patch: PurchaseOrderCreateRequest{
				Data: &PurchaseOrderData{
					OrderNumber:  &purchaseOrder.OrderNumber,
					OrderDate:    &validOrderDate,
					TrackingCode: &purchaseOrder.TrackingCode,
					BuyerId:      &purchaseOrder.BuyerId,
					PurchaseOrderDetails: []PurchaseOrderDetailRequest{
						{ProductRecordId: &validProductRecordId, Quantity: &invalidQuantity},
					},
				},
			},
			value: &invalidQuantity,
		},
	}

	for _, field := range invalidDetailFields {
		t.Run("should return error if "+field.name+" is invalid", func(t *testing.T) {
			err := field.patch.Verify()
			require.Error(t, err)
			require.Equal(t, &custom_errors.InvalidArgValueErr{
				Argument:  field.name,
				Value:     field.value,
				ExtraInfo: "Value must be greater than 0",
			}, err)
		})
	}

	t.Run("should return nil if all fields are valid", func(t *testing.T) {
		purchaseOrderPatch := PurchaseOrderCreateRequest{
			Data: &PurchaseOrderData{
				OrderNumber:  &purchaseOrder.OrderNumber,
				OrderDate:    &validOrderDate,
				TrackingCode: &purchaseOrder.TrackingCode,
				BuyerId:      &purchaseOrder.BuyerId,
				PurchaseOrderDetails: []PurchaseOrderDetailRequest{
					{ProductRecordId: &validProductRecordId, Quantity: &validQuantity},
				},
			},
		}

		err := purchaseOrderPatch.Verify()
		require.NoError(t, err)
	})
}

func TestPurchaseOrderCreateRequest_ToPurchaseOrder(t *testing.T) {
	orderNumber := "ORD-12345"
	orderDate := "2024-01-15"
	trackingCode := "TRK-67890"
	buyerId := 10
	productRecordId1 := 1
	quantity1 := 5
	productRecordId2 := 2
	quantity2 := 3

	t.Run("should convert PurchaseOrderCreateRequest to PurchaseOrder correctly", func(t *testing.T) {
		request := PurchaseOrderCreateRequest{
			Data: &PurchaseOrderData{
				OrderNumber:  &orderNumber,
				OrderDate:    &orderDate,
				TrackingCode: &trackingCode,
				BuyerId:      &buyerId,
				PurchaseOrderDetails: []PurchaseOrderDetailRequest{
					{ProductRecordId: &productRecordId1, Quantity: &quantity1},
				},
			},
		}

		result := request.ToPurchaseOrder()

		require.Equal(t, orderNumber, result.OrderNumber)
		require.Equal(t, trackingCode, result.TrackingCode)
		require.Equal(t, buyerId, result.BuyerId)

		expectedDate, err := time.Parse("2006-01-02", orderDate)
		require.NoError(t, err)
		require.Equal(t, expectedDate, result.OrderDate)

		require.Len(t, result.PurchaseOrderDetails, 1)
		require.Equal(t, productRecordId1, result.PurchaseOrderDetails[0].ProductRecordId)
		require.Equal(t, quantity1, result.PurchaseOrderDetails[0].Quantity)
		require.Equal(t, 0, result.PurchaseOrderDetails[0].Id) // ID debería ser 0 ya que es nuevo
	})

	t.Run("should convert multiple purchase order details correctly", func(t *testing.T) {
		request := PurchaseOrderCreateRequest{
			Data: &PurchaseOrderData{
				OrderNumber:  &orderNumber,
				OrderDate:    &orderDate,
				TrackingCode: &trackingCode,
				BuyerId:      &buyerId,
				PurchaseOrderDetails: []PurchaseOrderDetailRequest{
					{ProductRecordId: &productRecordId1, Quantity: &quantity1},
					{ProductRecordId: &productRecordId2, Quantity: &quantity2},
				},
			},
		}

		result := request.ToPurchaseOrder()

		require.Len(t, result.PurchaseOrderDetails, 2)

		require.Equal(t, productRecordId1, result.PurchaseOrderDetails[0].ProductRecordId)
		require.Equal(t, quantity1, result.PurchaseOrderDetails[0].Quantity)
		require.Equal(t, 0, result.PurchaseOrderDetails[0].Id)

		require.Equal(t, productRecordId2, result.PurchaseOrderDetails[1].ProductRecordId)
		require.Equal(t, quantity2, result.PurchaseOrderDetails[1].Quantity)
		require.Equal(t, 0, result.PurchaseOrderDetails[1].Id)
	})

	t.Run("should parse date correctly", func(t *testing.T) {
		testDate := "2023-12-25"
		request := PurchaseOrderCreateRequest{
			Data: &PurchaseOrderData{
				OrderNumber:  &orderNumber,
				OrderDate:    &testDate,
				TrackingCode: &trackingCode,
				BuyerId:      &buyerId,
				PurchaseOrderDetails: []PurchaseOrderDetailRequest{
					{ProductRecordId: &productRecordId1, Quantity: &quantity1},
				},
			},
		}

		result := request.ToPurchaseOrder()

		expectedDate := time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC)
		require.Equal(t, expectedDate, result.OrderDate)
	})

	t.Run("should handle empty purchase order details", func(t *testing.T) {
		request := PurchaseOrderCreateRequest{
			Data: &PurchaseOrderData{
				OrderNumber:          &orderNumber,
				OrderDate:            &orderDate,
				TrackingCode:         &trackingCode,
				BuyerId:              &buyerId,
				PurchaseOrderDetails: []PurchaseOrderDetailRequest{},
			},
		}

		result := request.ToPurchaseOrder()

		require.Empty(t, result.PurchaseOrderDetails)
		require.Len(t, result.PurchaseOrderDetails, 0)
	})
}
