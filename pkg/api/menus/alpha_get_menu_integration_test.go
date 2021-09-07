package menus

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/flypay/engineering-test/pkg/internal"
	"github.com/flypay/engineering-test/pkg/repository"
	"github.com/flypay/engineering-test/pkg/schema"
	"github.com/flypay/engineering-test/pkg/service"
)

// TestAlphaMenu calls unified service url to get alpha menu
func TestAlphaMenu_Integration(t *testing.T) {
	repo := repository.NewRepositoryImpl()
	serv := service.NewServiceImpl()
	alphaMenuExpected := schema.NewUnifiedMenuPopulatedAlphaMenuMock()
	t.Run("test alpha menu", func(t *testing.T) {
		alphaMenu := new(schema.Menu)
		alphaMenuResp := serv.RequestPOSClient(http.MethodGet,
			fmt.Sprintf("http://:8086"+internal.BasePathAlpha+internal.Menu), nil)
		if err := repo.DecodeReqRespBody(alphaMenuResp.Body, alphaMenu); err != nil {
			t.Fatalf("failed decoding body resp. err: %s", err.Error())
		}
		assertEqualMenus(t, alphaMenuExpected, alphaMenu)
	})
}

// assertEqualMenus runs deep equal on two menus and compares them
func assertEqualMenus(t *testing.T, expectedModel, retrievedModel *schema.Menu) {
	repo := repository.NewRepositoryImpl()
	repo.SortUnifiedMenu(expectedModel)
	repo.SortUnifiedMenu(retrievedModel)
	if !reflect.DeepEqual(expectedModel, retrievedModel) {
		t.Fatal("The two struct are not equal")
	}
}

/*
// assertEqualMenusAlternative is  an Alternative function to check equality of two structs
func assertEqualMenusAlternative(t *testing.T, expectedMenu, retrievedMenu *schema.Menu) {
	service.SortUnifiedMenu(expectedMenu)
	service.SortUnifiedMenu(retrievedMenu)
	jsonExpectedMenu, _ := json.Marshal(expectedMenu)
	jsonRetrievedMenu, _ := json.Marshal(retrievedMenu)
	if string(jsonExpectedMenu) != string(jsonRetrievedMenu) {
		t.Fatalf("the two menus are not equal")
	}
}
*/
