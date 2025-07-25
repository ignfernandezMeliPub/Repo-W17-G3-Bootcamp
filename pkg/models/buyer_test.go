package models

import (
	"app/pkg/custom_errors"
	"testing"

	"github.com/stretchr/testify/require"
)

var buyerEmptyStringWithSpaces = "  "

var buyer = Buyer{
	Id:           1,
	CardNumberId: "12345",
	FirstName:    "John",
	LastName:     "Doe",
}

func TestBuyerPatch_Verify(t *testing.T) {
	fields := []struct {
		name  string
		patch BuyerPatch
	}{
		{
			name: "card_number_id",
			patch: BuyerPatch{
				CardNumberId: &buyerEmptyStringWithSpaces,
				FirstName:    &buyer.FirstName,
				LastName:     &buyer.LastName,
			},
		},
		{
			name: "first_name",
			patch: BuyerPatch{
				CardNumberId: &buyer.CardNumberId,
				FirstName:    &buyerEmptyStringWithSpaces,
				LastName:     &buyer.LastName,
			},
		},
		{
			name: "last_name",
			patch: BuyerPatch{
				CardNumberId: &buyer.CardNumberId,
				FirstName:    &buyer.FirstName,
				LastName:     &buyerEmptyStringWithSpaces,
			},
		},
	}

	for _, field := range fields {
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

	t.Run("should return error if all fields are nil", func(t *testing.T) {
		buyerPatch := BuyerPatch{
			CardNumberId: nil,
			FirstName:    nil,
			LastName:     nil,
		}

		err := buyerPatch.Verify()
		require.Error(t, err)
		require.Equal(t, &custom_errors.MandatoryArgMissingErr{Argument: "card_number_id or first_name or last_name"}, err)
	})

	t.Run("should return nil if all fields are valid and not empty", func(t *testing.T) {
		buyerPatch := BuyerPatch{
			CardNumberId: &buyer.CardNumberId,
			FirstName:    &buyer.FirstName,
			LastName:     &buyer.LastName,
		}

		err := buyerPatch.Verify()
		require.NoError(t, err)
	})
}

func TestBuyer_Patch(t *testing.T) {
	t.Run("should patch the buyer fields one by one and all at once", func(t *testing.T) {
		patchData := struct {
			CardNumberId string
			FirstName    string
			LastName     string
		}{
			CardNumberId: "PatchCNID",
			FirstName:    "PatchFirst",
			LastName:     "PatchLast",
		}

		type patchCase struct {
			name  string
			patch BuyerPatch
			want  Buyer
		}

		cases := []patchCase{
			{
				name: "patch CardNumberId only",
				patch: BuyerPatch{
					CardNumberId: &patchData.CardNumberId,
				},
				want: Buyer{
					Id:           buyer.Id,
					CardNumberId: patchData.CardNumberId,
					FirstName:    buyer.FirstName,
					LastName:     buyer.LastName,
				},
			},
			{
				name: "patch FirstName only",
				patch: BuyerPatch{
					FirstName: &patchData.FirstName,
				},
				want: Buyer{
					Id:           buyer.Id,
					CardNumberId: buyer.CardNumberId,
					FirstName:    patchData.FirstName,
					LastName:     buyer.LastName,
				},
			},
			{
				name: "patch LastName only",
				patch: BuyerPatch{
					LastName: &patchData.LastName,
				},
				want: Buyer{
					Id:           buyer.Id,
					CardNumberId: buyer.CardNumberId,
					FirstName:    buyer.FirstName,
					LastName:     patchData.LastName,
				},
			},
			{
				name: "patch all fields",
				patch: BuyerPatch{
					CardNumberId: &patchData.CardNumberId,
					FirstName:    &patchData.FirstName,
					LastName:     &patchData.LastName,
				},
				want: Buyer{
					Id:           buyer.Id,
					CardNumberId: patchData.CardNumberId,
					FirstName:    patchData.FirstName,
					LastName:     patchData.LastName,
				},
			},
		}

		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				b := buyer
				b.Patch(tc.patch)
				require.Equal(t, tc.want, b)
			})
		}
	})
}

func TestBuyerCreateRequest_Verify(t *testing.T) {
	emptyFields := []struct {
		name  string
		patch BuyerCreateRequest
	}{
		{
			name: "card_number_id",
			patch: BuyerCreateRequest{
				CardNumberId: &buyerEmptyStringWithSpaces,
				FirstName:    &buyer.FirstName,
				LastName:     &buyer.LastName,
			},
		},
		{
			name: "first_name",
			patch: BuyerCreateRequest{
				CardNumberId: &buyer.CardNumberId,
				FirstName:    &buyerEmptyStringWithSpaces,
				LastName:     &buyer.LastName,
			},
		},
		{
			name: "last_name",
			patch: BuyerCreateRequest{
				CardNumberId: &buyer.CardNumberId,
				FirstName:    &buyer.FirstName,
				LastName:     &buyerEmptyStringWithSpaces,
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
		patch BuyerCreateRequest
	}{
		{
			name: "card_number_id",
			patch: BuyerCreateRequest{
				CardNumberId: nil,
				FirstName:    &buyer.FirstName,
				LastName:     &buyer.LastName,
			},
		},
		{
			name: "first_name",
			patch: BuyerCreateRequest{
				CardNumberId: &buyer.CardNumberId,
				FirstName:    nil,
				LastName:     &buyer.LastName,
			},
		},
		{
			name: "last_name",
			patch: BuyerCreateRequest{
				CardNumberId: &buyer.CardNumberId,
				FirstName:    &buyer.FirstName,
				LastName:     nil,
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

	t.Run("should return nil if all fields are valid", func(t *testing.T) {
		buyerPatch := BuyerCreateRequest{
			CardNumberId: &buyer.CardNumberId,
			FirstName:    &buyer.FirstName,
			LastName:     &buyer.LastName,
		}

		err := buyerPatch.Verify()
		require.NoError(t, err)
	})
}

func TestBuyerCreateRequest_ToBuyer(t *testing.T) {
	buyerCreateRequest := BuyerCreateRequest{
		CardNumberId: &buyer.CardNumberId,
		FirstName:    &buyer.FirstName,
		LastName:     &buyer.LastName,
	}

	buyerCreated := buyerCreateRequest.ToBuyer()
	buyerCreated.Id = buyer.Id
	require.Equal(t, buyer, buyerCreated)
}
