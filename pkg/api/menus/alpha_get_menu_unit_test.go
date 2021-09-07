package menus

import (
	"net/http"
	"testing"

	"github.com/flypay/engineering-test/pkg/repository"
	"github.com/flypay/engineering-test/pkg/schema"
)

func TestAlphaMenu_UnitTesting(t *testing.T) {
	repo := repository.NewRepositoryImpl()
	expectedAlpha := schema.NewUnifiedMenuPopulatedAlphaMenuMock()
	expectedAlpha.POS = "alpha"
	request := new(http.Request)
	alphaMenuHandler := NewGetMockAlphaMenu()
	request, _ = alphaMenuHandler.ParseArgs(request)
	resp := alphaMenuHandler.Process(request)
	alphaMenu := new(schema.Menu)
	_ = repo.DecodeReqRespBody(resp.Body, alphaMenu)
	t.Run("test alpha menu", func(t *testing.T) {
		assertEqualMenus(t, alphaMenu, expectedAlpha)
	})
}
