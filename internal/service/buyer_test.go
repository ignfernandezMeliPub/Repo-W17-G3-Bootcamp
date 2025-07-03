package service

import (
	"app/internal/repository/buyer_repository"
	"app/pkg/custom_errors"
	"app/pkg/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getBuyerService() BuyerService {
	db := map[int]models.Buyer{
		1:  {Id: 1, CardNumberId: "1001", FirstName: "Juan", LastName: "Pérez"},
		2:  {Id: 2, CardNumberId: "1002", FirstName: "Ana", LastName: "García"},
		3:  {Id: 3, CardNumberId: "1003", FirstName: "Luis", LastName: "Rodríguez"},
		4:  {Id: 4, CardNumberId: "1004", FirstName: "Laura", LastName: "Martínez"},
		5:  {Id: 5, CardNumberId: "1005", FirstName: "Carlos", LastName: "Pérez"},
		6:  {Id: 6, CardNumberId: "1006", FirstName: "Ana", LastName: "López"},
		7:  {Id: 7, CardNumberId: "1007", FirstName: "Mario", LastName: "García"},
		8:  {Id: 8, CardNumberId: "1008", FirstName: "Elena", LastName: "Hernández"},
		9:  {Id: 9, CardNumberId: "1009", FirstName: "Carmen", LastName: "Martínez"},
		10: {Id: 10, CardNumberId: "1010", FirstName: "Luis", LastName: "Alonso"},
	}
	repo := buyer_repository.NewBuyerMap(db)
	return NewBuyerDefault(repo)
}

func TestGetAllBuyers(t *testing.T) {
	s := getBuyerService()

	b, err := s.GetAllBuyers()

	require.NoError(t, err)
	assert.Equal(t, 10, len(b))
}
func TestGetBuyerByID(t *testing.T) {
	s := getBuyerService()

	t.Run("Success", func(t *testing.T) {
		b, err := s.GetBuyerById(1)
		assert.NoError(t, err)
		assert.Equal(t, "Juan", b.FirstName)
		assert.Equal(t, "Pérez", b.LastName)
		assert.Equal(t, "1001", b.CardNumberId)
	})

	t.Run("Not found", func(t *testing.T) {
		_, err := s.GetBuyerById(999)
		assert.Error(t, err)
		assert.ErrorAs(t, err, &custom_errors.ErrNotFound)
	})
}

func TestGetBuyerByCardNumberID(t *testing.T) {
	s := getBuyerService()

	t.Run("Success", func(t *testing.T) {
		b, err := s.GetBuyerByCardNumberId("1002")
		assert.NoError(t, err)
		assert.Equal(t, "Ana", b.FirstName)
		assert.Equal(t, "García", b.LastName)
		assert.Equal(t, 2, b.Id)
	})

	t.Run("Not found", func(t *testing.T) {
		_, err := s.GetBuyerByCardNumberId("")
		assert.Error(t, err)
		assert.ErrorAs(t, err, &custom_errors.ErrNotFound)
	})
}

func TestCreateBuyer(t *testing.T) {
	s := getBuyerService()

	t.Run("Success", func(t *testing.T) {
		_b := models.Buyer{Id: -1, CardNumberId: "9876", FirstName: "Alfonso", LastName: "Gregorio"}

		b, err := s.CreateBuyer(_b)
		assert.NoError(t, err)
		assert.Equal(t, _b.FirstName, b.FirstName)
		assert.Equal(t, _b.LastName, b.LastName)
		assert.Equal(t, _b.CardNumberId, b.CardNumberId)
		assert.NotEqual(t, _b.Id, b.Id)
	})

	t.Run("Fail to create due Unique constraint", func(t *testing.T) {
		_b := models.Buyer{Id: -1, CardNumberId: "1001", FirstName: "Alfonso", LastName: "Gregorio"}

		_, err := s.CreateBuyer(_b)
		assert.Error(t, err)
		assert.ErrorAs(t, err, &custom_errors.ErrUniqueAttributeViolationError)
	})

	t.Run("Fail to create due invalid buyer", func(t *testing.T) {
		_b := models.Buyer{Id: -1, CardNumberId: "100132", FirstName: "", LastName: "Gregorio"}

		_, err := s.CreateBuyer(_b)
		assert.Error(t, err)
		assert.ErrorAs(t, err, &custom_errors.ErrInvalidArgs)
	})
}

func TestUpdateBuyerByID(t *testing.T) {
	s := getBuyerService()

	t.Run("Success updating all fields", func(t *testing.T) {
		id := 1

		CardNumberId := "9876"
		FirstName := "New name"
		LastName := "New last name"

		_b := models.BuyerPatch{CardNumberId: &CardNumberId, FirstName: &FirstName, LastName: &LastName}

		b, err := s.UpdateBuyerById(id, _b)
		assert.NoError(t, err)
		assert.Equal(t, *_b.FirstName, b.FirstName)
		assert.Equal(t, *_b.LastName, b.LastName)
		assert.Equal(t, *_b.CardNumberId, b.CardNumberId)
		assert.Equal(t, id, b.Id)
	})

	t.Run("Success updating only one field", func(t *testing.T) {
		id := 1

		LastName := "Re updated last name"

		buyer_pre_update, _ := s.GetBuyerById(id)

		_b := models.BuyerPatch{LastName: &LastName}

		b, err := s.UpdateBuyerById(id, _b)
		assert.NoError(t, err)
		assert.Equal(t, buyer_pre_update.FirstName, b.FirstName)
		assert.Equal(t, LastName, b.LastName)
		assert.Equal(t, buyer_pre_update.CardNumberId, b.CardNumberId)
		assert.Equal(t, id, b.Id)
	})

	t.Run("Success updating CardNumberID repeating self value", func(t *testing.T) {
		id := 1
		CardNumberID := "1001"

		_b := models.BuyerPatch{CardNumberId: &CardNumberID}

		b, err := s.UpdateBuyerById(id, _b)
		assert.NoError(t, err)
		assert.Equal(t, CardNumberID, b.CardNumberId)
	})

	t.Run("Fail to update due Unique constraint", func(t *testing.T) {
		id := 1
		CardNumberID := "1002"

		_b := models.BuyerPatch{CardNumberId: &CardNumberID}

		_, err := s.UpdateBuyerById(id, _b)
		assert.Error(t, err)
		assert.ErrorAs(t, err, &custom_errors.ErrUniqueAttributeViolationError)
	})

	t.Run("Fail to update with blank patch info", func(t *testing.T) {
		id := 1
		CardNumberId := ""
		_b := models.BuyerPatch{CardNumberId: &CardNumberId}

		_, err := s.UpdateBuyerById(id, _b)
		assert.Error(t, err)
		assert.ErrorAs(t, err, &custom_errors.ErrInvalidArgs)
	})
}

func TestDeleteBuyerByID(t *testing.T) {
	s := getBuyerService()

	t.Run("Success", func(t *testing.T) {
		id := 1

		err := s.DeleteBuyerById(id)
		assert.NoError(t, err)
	})

	t.Run("Not found", func(t *testing.T) {
		id := 1

		err := s.DeleteBuyerById(id)
		assert.Error(t, err)
		assert.ErrorAs(t, err, &custom_errors.ErrNotFound)
	})
}
