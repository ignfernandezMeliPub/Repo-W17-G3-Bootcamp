package models

import (
	"app/pkg/custom_errors"
	"testing"

	"github.com/stretchr/testify/require"
)

var inboundOrderDetails = InboundOrderDetails{
	OrderDate:      "2021-04-04",
	OrderNumber:    "77777",
	EmployeeId:     4,
	ProductBatchId: 1,
	WarehouseId:    1,
}

func TestInboundOrderPost_Verify(t *testing.T) {

	nilFields := []struct {
		name string
		post InboundOrderRequestBody
	}{
		{
			name: "data.order_date",
			post: InboundOrderRequestBody{
				Data: &InboundOrderData{

					OrderDate:      nil,
					OrderNumber:    &inboundOrderDetails.OrderNumber,
					EmployeeId:     &inboundOrderDetails.EmployeeId,
					ProductBatchId: &inboundOrderDetails.ProductBatchId,
					WarehouseId:    &inboundOrderDetails.WarehouseId,
				},
			},
		},
		{
			name: "data.order_number",
			post: InboundOrderRequestBody{
				Data: &InboundOrderData{

					OrderDate:      &inboundOrderDetails.OrderDate,
					OrderNumber:    nil,
					EmployeeId:     &inboundOrderDetails.EmployeeId,
					ProductBatchId: &inboundOrderDetails.ProductBatchId,
					WarehouseId:    &inboundOrderDetails.WarehouseId,
				},
			},
		},
		{
			name: "data.employee_id",
			post: InboundOrderRequestBody{
				Data: &InboundOrderData{

					OrderDate:      &inboundOrderDetails.OrderDate,
					OrderNumber:    &inboundOrderDetails.OrderNumber,
					EmployeeId:     nil,
					ProductBatchId: &inboundOrderDetails.ProductBatchId,
					WarehouseId:    &inboundOrderDetails.WarehouseId,
				},
			},
		},
		{
			name: "data.product_batch_id",
			post: InboundOrderRequestBody{
				Data: &InboundOrderData{

					OrderDate:      &inboundOrderDetails.OrderDate,
					OrderNumber:    &inboundOrderDetails.OrderNumber,
					EmployeeId:     &inboundOrderDetails.EmployeeId,
					ProductBatchId: nil,
					WarehouseId:    &inboundOrderDetails.WarehouseId,
				},
			},
		},
		{
			name: "data.warehouse_id",
			post: InboundOrderRequestBody{
				Data: &InboundOrderData{

					OrderDate:      &inboundOrderDetails.OrderDate,
					OrderNumber:    &inboundOrderDetails.OrderNumber,
					EmployeeId:     &inboundOrderDetails.EmployeeId,
					ProductBatchId: &inboundOrderDetails.ProductBatchId,
					WarehouseId:    nil,
				},
			},
		}, {
			name: "data",
			post: InboundOrderRequestBody{
				Data: nil,
			},
		},
	}

	for _, field := range nilFields {
		t.Run("should return error if "+field.name+" is nil", func(t *testing.T) {
			err := field.post.Verify()
			require.Error(t, err)
			require.Equal(t, &custom_errors.MandatoryArgMissingErr{
				Argument: field.name,
			}, err)
		})
	}

	t.Run("error, invalid date format", func(t *testing.T) {

		wrongDate := "2020/12/06"

		inboundOrderBody := InboundOrderRequestBody{
			Data: &InboundOrderData{

				OrderDate:      &wrongDate,
				OrderNumber:    &inboundOrderDetails.OrderNumber,
				EmployeeId:     &inboundOrderDetails.EmployeeId,
				ProductBatchId: &inboundOrderDetails.ProductBatchId,
				WarehouseId:    &inboundOrderDetails.WarehouseId,
			},
		}

		err := inboundOrderBody.Verify()
		require.Error(t, err)
		require.Error(t, err)
		require.Equal(t, &custom_errors.InvalidArgValueErr{
			Argument: "data.order_date", Value: wrongDate, ExtraInfo: "Invalid date format.",
		}, err)
	})

	t.Run("error, future date", func(t *testing.T) {

		futureDate := "2030-12-06"

		inboundOrderBody := InboundOrderRequestBody{
			Data: &InboundOrderData{

				OrderDate:      &futureDate,
				OrderNumber:    &inboundOrderDetails.OrderNumber,
				EmployeeId:     &inboundOrderDetails.EmployeeId,
				ProductBatchId: &inboundOrderDetails.ProductBatchId,
				WarehouseId:    &inboundOrderDetails.WarehouseId,
			},
		}

		err := inboundOrderBody.Verify()
		require.Error(t, err)
		require.Error(t, err)
		require.Equal(t, &custom_errors.InvalidArgValueErr{
			Argument: "data.order_date", Value: futureDate, ExtraInfo: "Future date.",
		}, err)
	})

	t.Run("should return nil if all fields are valid", func(t *testing.T) {
		inboundOrderBody := InboundOrderRequestBody{
			Data: &InboundOrderData{

				OrderDate:      &inboundOrderDetails.OrderDate,
				OrderNumber:    &inboundOrderDetails.OrderNumber,
				EmployeeId:     &inboundOrderDetails.EmployeeId,
				ProductBatchId: &inboundOrderDetails.ProductBatchId,
				WarehouseId:    &inboundOrderDetails.WarehouseId,
			},
		}

		err := inboundOrderBody.Verify()
		require.NoError(t, err)
	})

}
