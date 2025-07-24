package service

import (
	"app/test/repository"
	"fmt"
	"testing"
)

func TestSectionsServiceImpl_UpdateSectionById(t *testing.T) {
	dbMock := repository.NewSectionsMock()
	svSc := NewSectionsService(dbMock)
	fmt.Println(svSc)

	t.Run("success", func(t *testing.T) {

	})
}
