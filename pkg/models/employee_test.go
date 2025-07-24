package models

import (
	"app/pkg/custom_errors"
	"testing"

	"github.com/stretchr/testify/require"
)

var employeeAttributes = EmployeeAttributes{
	CardNumberId: "12345",
	FirstName:    "John",
	LastName:     "Doe",
	WarehouseId:  1,
}

func TestEmployeePatch_Verify(t *testing.T) {

	fields := []struct {
		name  string
		patch EmployeePatchRequestBody
	}{
		{
			name: "card_number_id, first_name, last_name, warehouse_id",
			patch: EmployeePatchRequestBody{
				CardNumberId: &employeeAttributes.CardNumberId,
				FirstName:    &employeeAttributes.FirstName,
				LastName:     &employeeAttributes.LastName,
				WarehouseId:  &employeeAttributes.WarehouseId,
			},
		}, {
			name: "first_name, last_name, warehouse_id",
			patch: EmployeePatchRequestBody{
				CardNumberId: nil,
				FirstName:    &employeeAttributes.FirstName,
				LastName:     &employeeAttributes.LastName,
				WarehouseId:  &employeeAttributes.WarehouseId,
			},
		},
		{
			name: "card_number_id, last_name, warehouse_id",
			patch: EmployeePatchRequestBody{
				CardNumberId: &employeeAttributes.CardNumberId,
				FirstName:    nil,
				LastName:     &employeeAttributes.LastName,
				WarehouseId:  &employeeAttributes.WarehouseId,
			},
		},
		{
			name: "card_number_id, first_name, warehouse_id",
			patch: EmployeePatchRequestBody{
				CardNumberId: &employeeAttributes.CardNumberId,
				FirstName:    &employeeAttributes.FirstName,
				LastName:     nil,
				WarehouseId:  &employeeAttributes.WarehouseId,
			},
		},
		{
			name: "card_number_id, first_name, last_name",
			patch: EmployeePatchRequestBody{
				CardNumberId: &employeeAttributes.CardNumberId,
				FirstName:    &employeeAttributes.FirstName,
				LastName:     &employeeAttributes.LastName,
				WarehouseId:  nil,
			},
		}, {
			name: "last_name, warehouse_id",
			patch: EmployeePatchRequestBody{
				CardNumberId: nil,
				FirstName:    nil,
				LastName:     &employeeAttributes.LastName,
				WarehouseId:  &employeeAttributes.WarehouseId,
			},
		}, {
			name: "first_name, warehouse_id",
			patch: EmployeePatchRequestBody{
				CardNumberId: nil,
				FirstName:    &employeeAttributes.FirstName,
				LastName:     nil,
				WarehouseId:  &employeeAttributes.WarehouseId,
			},
		}, {
			name: "first_name, last_name",
			patch: EmployeePatchRequestBody{
				CardNumberId: nil,
				FirstName:    &employeeAttributes.FirstName,
				LastName:     &employeeAttributes.LastName,
				WarehouseId:  nil,
			},
		}, {
			name: "card_number_id, warehouse_id",
			patch: EmployeePatchRequestBody{
				CardNumberId: &employeeAttributes.CardNumberId,
				FirstName:    nil,
				LastName:     nil,
				WarehouseId:  &employeeAttributes.WarehouseId,
			},
		}, {
			name: "card_number_id, last_name",
			patch: EmployeePatchRequestBody{
				CardNumberId: &employeeAttributes.CardNumberId,
				FirstName:    nil,
				LastName:     &employeeAttributes.LastName,
				WarehouseId:  nil,
			},
		}, {
			name: "card_number_id, first_name",
			patch: EmployeePatchRequestBody{
				CardNumberId: &employeeAttributes.CardNumberId,
				FirstName:    &employeeAttributes.FirstName,
				LastName:     nil,
				WarehouseId:  nil,
			},
		}, {
			name: "card_number_id",
			patch: EmployeePatchRequestBody{
				CardNumberId: &employeeAttributes.CardNumberId,
				FirstName:    nil,
				LastName:     nil,
				WarehouseId:  nil,
			},
		}, {
			name: "first_name",
			patch: EmployeePatchRequestBody{
				CardNumberId: nil,
				FirstName:    &employeeAttributes.FirstName,
				LastName:     nil,
				WarehouseId:  nil,
			},
		}, {
			name: "last_name",
			patch: EmployeePatchRequestBody{
				CardNumberId: nil,
				FirstName:    nil,
				LastName:     &employeeAttributes.LastName,
				WarehouseId:  nil,
			},
		}, {
			name: " warehouse_id",
			patch: EmployeePatchRequestBody{
				CardNumberId: nil,
				FirstName:    nil,
				LastName:     nil,
				WarehouseId:  &employeeAttributes.WarehouseId,
			},
		},
	}

	for _, field := range fields {
		t.Run("should return nil if "+field.name+" is/are filled", func(t *testing.T) {
			err := field.patch.Verify()
			require.NoError(t, err)
		})
	}

	t.Run("Return error if all fields are nil", func(t *testing.T) {
		employeePatchRequestBody := EmployeePatchRequestBody{
			CardNumberId: nil,
			FirstName:    nil,
			LastName:     nil,
			WarehouseId:  nil,
		}

		err := employeePatchRequestBody.Verify()
		require.Error(t, err)
		require.Equal(t, &custom_errors.MandatoryArgMissingErr{Argument: "card_number_id or first_name or last_name or warehouse_id"}, err)
	})

}

func TestEmployeePost_Verify(t *testing.T) {

	nilFields := []struct {
		name string
		post EmployeePostRequestBody
	}{
		{
			name: "card_number_id",
			post: EmployeePostRequestBody{
				CardNumberId: nil,
				FirstName:    &employeeAttributes.FirstName,
				LastName:     &employeeAttributes.LastName,
				WarehouseId:  &employeeAttributes.WarehouseId,
			},
		},
		{
			name: "first_name",
			post: EmployeePostRequestBody{
				CardNumberId: &employeeAttributes.CardNumberId,
				FirstName:    nil,
				LastName:     &employeeAttributes.LastName,
				WarehouseId:  &employeeAttributes.WarehouseId,
			},
		},
		{
			name: "last_name",
			post: EmployeePostRequestBody{
				CardNumberId: &employeeAttributes.CardNumberId,
				FirstName:    &employeeAttributes.FirstName,
				LastName:     nil,
				WarehouseId:  &employeeAttributes.WarehouseId,
			},
		},
		{
			name: "warehouse_id",
			post: EmployeePostRequestBody{
				CardNumberId: &employeeAttributes.CardNumberId,
				FirstName:    &employeeAttributes.FirstName,
				LastName:     &employeeAttributes.LastName,
				WarehouseId:  nil,
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

	t.Run("should return nil if all fields are valid", func(t *testing.T) {
		employeePost := EmployeePostRequestBody{
			CardNumberId: &employeeAttributes.CardNumberId,
			FirstName:    &employeeAttributes.FirstName,
			LastName:     &employeeAttributes.LastName,
			WarehouseId:  &employeeAttributes.WarehouseId,
		}

		err := employeePost.Verify()
		require.NoError(t, err)
	})

}
