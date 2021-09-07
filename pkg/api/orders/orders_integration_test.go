package orders

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

// TestAlphaOrderResponse calls unified service url to set an order to alpha pos
func TestAlphaOrderResponse(t *testing.T) {
	// Place an order for
	// 	1x Individual Mushroom Burger
	//	 	Extra Mayo
	//		No Tomato
	// 	1x Individual Buttermilk Chicken Burger
	// 	2x Regular Fries
	//		Paprika Salt
	//  2x Large Fizzy Cola
	requestOrder := schema.NewUnifiedOrderRequestPopulatedOrderAlphaMock()
	repo := repository.NewRepositoryImpl()
	serv := service.NewServiceImpl()
	requestPayload, _ := repo.EncodeReqRespBody(requestOrder)
	alphaOrderResponse := new(schema.OrderResponse)
	alphaOrderResp := serv.RequestPOSClient(http.MethodPost,
		fmt.Sprintf("http://:8086"+internal.Orders), requestPayload)

	if err := repo.DecodeReqRespBody(alphaOrderResp.Body, alphaOrderResponse); err != nil {
		t.Fatalf("failed decoding body resp. err: %s", err.Error())
	}
	expectedResponse := schema.NewUnifiedRespPopulatedAlphaOrderRespMock()
	assertEqualOrders(t, expectedResponse, alphaOrderResponse)
}

// TestBetaOrderResponse calls unified service url to set an order to beta pos
func TestBetaOrderResponse(t *testing.T) {
	requestOrder := schema.NewUnifiedOrderRequestPopulatedOrderBetaMock()
	repo := repository.NewRepositoryImpl()
	serv := service.NewServiceImpl()
	requestPayload, _ := repo.EncodeReqRespBody(requestOrder)
	betaOrderResponse := new(schema.OrderResponse)
	betaOrderResp := serv.RequestPOSClient(http.MethodPost,
		fmt.Sprintf("http://:8086"+internal.Orders), requestPayload)

	if err := repo.DecodeReqRespBody(betaOrderResp.Body, betaOrderResponse); err != nil {
		t.Fatalf("failed decoding body resp. err: %s", err.Error())
	}
	expectedResponse := schema.NewUnifiedRespPopulatedBetaOrderRespMock()
	assertEqualOrders(t, expectedResponse, betaOrderResponse)
}

// assertEqualOrders runs deep equal on two order response bodies and compares them
func assertEqualOrders(t *testing.T, expectedModel, retrievedModel *schema.OrderResponse) {
	repo := repository.NewRepositoryImpl()
	repo.SortUnifiedResponse(expectedModel)
	repo.SortUnifiedResponse(retrievedModel)
	if !reflect.DeepEqual(expectedModel, retrievedModel) {
		t.Fatal("The two struct are not equal")
	}
}
