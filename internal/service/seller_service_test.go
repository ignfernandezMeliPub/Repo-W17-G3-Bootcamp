package service

import (
	"app/internal/repository/seller_repository"
	"app/pkg/custom_errors"
	"app/pkg/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func getSellerServiceImpl() SellerServiceImpl {
	db := map[int]models.Seller{
		1:  {Id: 1, CompanyId: 1001, CompanyName: "Acme Corp", Address: "Av. Siempre Viva 123", Telephone: "111-1111"},
		2:  {Id: 2, CompanyId: 1002, CompanyName: "Globex SA", Address: "Calle Falsa 456", Telephone: "222-2222"},
		3:  {Id: 3, CompanyId: 1003, CompanyName: "Soylent Ltd.", Address: "Boulevard Central 789", Telephone: "333-3333"},
		4:  {Id: 4, CompanyId: 1004, CompanyName: "Initech", Address: "Diagonal Norte 321", Telephone: "444-4444"},
		5:  {Id: 5, CompanyId: 1005, CompanyName: "Umbrella", Address: "Av. de Mayo 654", Telephone: "555-5555"},
		6:  {Id: 6, CompanyId: 1006, CompanyName: "Cyberdyne", Address: "Rivadavia 987", Telephone: "666-6666"},
		7:  {Id: 7, CompanyId: 1007, CompanyName: "Stark Industries", Address: "Mitre 147", Telephone: "777-7777"},
		8:  {Id: 8, CompanyId: 1008, CompanyName: "Wayne Enterprises", Address: "San Martín 258", Telephone: "888-8888"},
		9:  {Id: 9, CompanyId: 1009, CompanyName: "Wonka Inc.", Address: "Esmeralda 369", Telephone: "999-9999"},
		10: {Id: 10, CompanyId: 1010, CompanyName: "Oscorp", Address: "Callao 753", Telephone: "101-0101"},
	}

	repo := seller_repository.NewSellerRepositoryMap(db)
	return NewSellerService(&repo)
}

// TestCreateSeller_Success verifies that creating a new seller succeeds.
func TestCreateSeller_Success(t *testing.T) {
	service := getSellerServiceImpl()
	seller, err := service.Create(11, 2001, "NewCo", "Nueva 123", "5555-6666")

	require.NoError(t, err)
	assert.Equal(t, 11, seller.Id)
	assert.Equal(t, 2001, seller.CompanyId)
	assert.Equal(t, "NewCo", seller.CompanyName)

	got, err := service.GetById(11)
	require.NoError(t, err)
	assert.Equal(t, "NewCo", got.CompanyName)
}

// TestCreateSeller_DuplicateId verifies that creating a seller with a duplicate ID fails.
func TestCreateSeller_DuplicateId(t *testing.T) {
	service := getSellerServiceImpl()
	_, err := service.Create(1, 9999, "Otra", "Dir", "555")

	var unique *custom_errors.UniqueAttributeViolationErr
	require.ErrorAs(t, err, &unique)

	assert.Equal(t, "id", unique.AttributeName)
}

// TestCreateSeller_DuplicateCompanyId verifies that creating a seller with a duplicate CompanyID fails.
func TestCreateSeller_DuplicateCompanyId(t *testing.T) {
	service := getSellerServiceImpl()
	_, err := service.Create(99, 1001, "Otra", "Dir", "555")

	var unique *custom_errors.UniqueAttributeViolationErr
	require.ErrorAs(t, err, &unique)
	assert.Equal(t, "companyId", unique.AttributeName)
}

// TestGetById_NotFound verifies that GetById returns an error when the seller does not exist.
func TestGetById_NotFound(t *testing.T) {
	service := getSellerServiceImpl()
	_, err := service.GetById(123456)

	require.Error(t, err)

	notFoundErrorPoint := &custom_errors.ResourceNotFoundError{}
	assert.ErrorAs(t, err, &notFoundErrorPoint)
}

// TestGetAll verifies that GetAll returns all stored sellers.
func TestGetAll(t *testing.T) {
	service := getSellerServiceImpl()
	all, err := service.GetAll()

	require.NoError(t, err)
	assert.Len(t, all, 10)
}

// TestDeleteSeller verifies that deleting a seller works, and that deleting a non-existent seller returns an error.
func TestDeleteSeller(t *testing.T) {
	service := getSellerServiceImpl()
	err := service.Delete(2)

	require.NoError(t, err)

	err = service.Delete(2)
	var notFound *custom_errors.ResourceNotFoundError
	require.ErrorAs(t, err, &notFound)
}

// TestPatch_Success verifies that Patch correctly updates the fields of an existing seller.
func TestPatch_Success(t *testing.T) {
	service := getSellerServiceImpl()

	tel := "111-8888"
	name := "Soylent Inc."
	s, err := service.Patch(3, nil, &name, nil, &tel)

	require.NoError(t, err)
	assert.Equal(t, "111-8888", s.Telephone)
	assert.Equal(t, "Soylent Inc.", s.CompanyName)
}

// TestPatch_DuplicateCompanyId verifies that Patch fails if the new CompanyID already exists, and that no changes are applied.
func TestPatch_DuplicateCompanyId(t *testing.T) {
	service := getSellerServiceImpl()

	id := 1001
	_, err := service.Patch(4, &id, nil, nil, nil)

	var unique *custom_errors.UniqueAttributeViolationErr
	require.ErrorAs(t, err, &unique)
	assert.Equal(t, "companyId", unique.AttributeName)
}

// TestPatch_NotFound verifies that Patch returns an error when attempting to update a non-existent seller.
func TestPatch_NotFound(t *testing.T) {
	service := getSellerServiceImpl()

	cn := "Other"
	_, err := service.Patch(9999, nil, &cn, nil, nil)

	var notfound *custom_errors.ResourceNotFoundError
	require.ErrorAs(t, err, &notfound)
}
