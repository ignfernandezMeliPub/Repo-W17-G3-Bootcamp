package models

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCarries_Verify(t *testing.T) {

	carrie := Carries{
		Cid:         "1",
		CompanyName: "Company 1",
		Address:     "Street 123",
		Telephone:   "1234567890",
		LocalityId:  "L001",
	}
	tests := []struct {
		name          string
		fieldToModify func(c Carries) Carries
		errorExpected string
	}{
		{
			name: "should return error if the cid is empty",
			fieldToModify: func(c Carries) Carries {
				c.Cid = ""
				return c
			},
			errorExpected: "cid",
		},
		{
			name: "should return error if the company name is empty",
			fieldToModify: func(c Carries) Carries {
				c.CompanyName = ""
				return c
			},
			errorExpected: "company_name",
		},
		{
			name: "should return error if the address is empty",
			fieldToModify: func(c Carries) Carries {
				c.Address = ""
				return c
			},
			errorExpected: "address",
		},
		{
			name: "should return error if the telephone is empty",
			fieldToModify: func(c Carries) Carries {
				c.Telephone = ""
				return c
			},
			errorExpected: "telephone",
		},
		{
			name: "should return error if the locality id is empty",
			fieldToModify: func(c Carries) Carries {
				c.LocalityId = ""
				return c
			},
			errorExpected: "locality_id",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			inputCarrie := test.fieldToModify(carrie)
			err := inputCarrie.Verify()
			require.Error(t, err)
			require.Contains(t, err.Error(), test.errorExpected)
		})
	}
}
