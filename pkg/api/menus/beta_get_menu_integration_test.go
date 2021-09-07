package menus

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/flypay/engineering-test/pkg/internal"
	"github.com/flypay/engineering-test/pkg/repository"
	"github.com/flypay/engineering-test/pkg/schema"
	"github.com/flypay/engineering-test/pkg/service"
)

// TestBetaMenu calls unified service url to get beta menu
func TestBetaMenu(t *testing.T) {
	betaMenuExpected := schema.NewUnifiedMenuPopulatedBetaMenuMock()
	repo := repository.NewRepositoryImpl()
	serv := service.NewServiceImpl()
	t.Run("get beta menu", func(t *testing.T) {
		betaMenu := new(schema.Menu)
		betaMenuResp := serv.RequestPOSClient(http.MethodGet,
			fmt.Sprintf("http://:8086"+internal.BasePathBeta+internal.Menu), nil)
		if err := repo.DecodeReqRespBody(betaMenuResp.Body, betaMenu); err != nil {
			t.Fatalf("failed decoding body resp. err: %s", err.Error())
		}
		assertEqualMenus(t, betaMenuExpected, betaMenu)
	})
}
