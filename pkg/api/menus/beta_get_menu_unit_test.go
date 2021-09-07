package menus

import (
	"net/http"
	"testing"

	"github.com/flypay/engineering-test/pkg/repository"
	"github.com/flypay/engineering-test/pkg/schema"
)

func TestBetaMenu_UnitTesting(t *testing.T) {
	repo := repository.NewRepositoryImpl()
	expectedBeta := schema.NewUnifiedMenuPopulatedBetaMenuMock()
	expectedBeta.POS = "beta"
	request := new(http.Request)
	betaMenuHandler := NewGetMockBetaMenu()
	request, _ = betaMenuHandler.ParseArgs(request)
	resp := betaMenuHandler.Process(request)
	betaMenu := new(schema.Menu)
	_ = repo.DecodeReqRespBody(resp.Body, betaMenu)
	t.Run("test beta menu", func(t *testing.T) {
		assertEqualMenus(t, betaMenu, expectedBeta)
	})
}
