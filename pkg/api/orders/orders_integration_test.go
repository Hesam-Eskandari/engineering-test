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
	requestOrder := setUpAlphaOrder()
	repo := repository.NewRepositoryImpl()
	serv := service.NewServiceImpl()
	requestPayload, _ := repo.EncodeReqRespBody(requestOrder)
	alphaOrderResponse := new(schema.OrderResponse)
	alphaOrderResp := serv.RequestPOSClient(http.MethodPost,
		fmt.Sprintf("http://:8086"+internal.Orders), requestPayload)

	if err := repo.DecodeReqRespBody(alphaOrderResp.Body, alphaOrderResponse); err != nil {
		t.Fatalf("failed decoding body resp. err: %s", err.Error())
	}
	expectedResponse := setUpAlphaResponse()
	assertEqualOrders(t, expectedResponse, alphaOrderResponse)
}

// TestBetaOrderResponse calls unified service url to set an order to beta pos
func TestBetaOrderResponse(t *testing.T) {
	requestOrder := setUpBetaOrder()
	repo := repository.NewRepositoryImpl()
	serv := service.NewServiceImpl()
	requestPayload, _ := repo.EncodeReqRespBody(requestOrder)
	betaOrderResponse := new(schema.OrderResponse)
	betaOrderResp := serv.RequestPOSClient(http.MethodPost,
		fmt.Sprintf("http://:8086"+internal.Orders), requestPayload)

	if err := repo.DecodeReqRespBody(betaOrderResp.Body, betaOrderResponse); err != nil {
		t.Fatalf("failed decoding body resp. err: %s", err.Error())
	}
	expectedResponse := setUpBetaResponse()
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

// setUpAlphaOrder returns an alpha test order request body
func setUpAlphaOrder() (order *schema.OrderRequest) {
	order = new(schema.OrderRequest)
	order.POS = internal.POSAlpha
	order.ID = "12345"
	order.Items = []schema.OrderItem{
		{
			ID:          "6001",
			Quantity:    1,
			Size:        "8001",
			Ingredients: []string{"9004"},
			Extras:      []string{"9010"},
		},
		{
			ID:          "6002",
			Quantity:    1,
			Size:        "8001",
			Ingredients: []string{},
			Extras:      []string{},
		},
		{
			ID:          "6003",
			Quantity:    2,
			Size:        "8005",
			Ingredients: []string{},
			Extras:      []string{"9009"},
		},
		{
			ID:          "6004",
			Quantity:    2,
			Size:        "8006",
			Ingredients: []string{},
			Extras:      []string{},
		},
	}
	return
}

// setUpAlphaResponse returns expected alpha order response that should be returned
func setUpAlphaResponse() (resp *schema.OrderResponse) {
	resp = new(schema.OrderResponse)
	resp.POS = internal.POSAlpha
	resp.ID = "12345"
	resp.TotalPrice = 48.68
	resp.Items = []*schema.OrderResponseItem{
		{
			Name:        "Mushroom Burger",
			Quantity:    1,
			Size:        "Individual",
			Extras:      []string{"Extra Mayo"},
			Ingredients: []string{"Seeded Bun", "Mushroom Patty", "Lettuce"},
			Price:       15.94,
		},
		{
			Name:        "Buttermilk Chicken Burger",
			Quantity:    1,
			Size:        "Individual",
			Extras:      []string{},
			Ingredients: []string{"Brioche Bun", "Chicken Patty", "Pickles", "Slaw"},
			Price:       15.95,
		},
		{
			Name:        "Fries",
			Quantity:    2,
			Size:        "Regular",
			Extras:      []string{"Paprika Salt"},
			Ingredients: []string{},
			Price:       8.89,
		},
		{
			Name:        "Fizzy Cola",
			Quantity:    2,
			Size:        "Large",
			Extras:      []string{},
			Ingredients: []string{},
			Price:       7.9,
		},
	}
	return
}

// setUpBetaOrder returns an beta test order request body
func setUpBetaOrder() (order *schema.OrderRequest) {
	order = new(schema.OrderRequest)
	order.POS = internal.POSBeta
	order.ID = "54321"
	order.Items = []schema.OrderItem{
		{
			ID:          "hjrlho",
			Quantity:    1,
			Size:        "hjrlho",
			Ingredients: []string{},
			Extras:      []string{"cbnufj"},
		},
		{
			ID:       "ukhqnd",
			Quantity: 2,
			Size:     "ukhqnd",
		},
		{
			ID:       "ugqcbb",
			Quantity: 1,
			Size:     "ugqcbb",
		},
	}
	return
}

// setUpBetaResponse returns expected beta order response that should be returned
func setUpBetaResponse() (resp *schema.OrderResponse) {
	resp = new(schema.OrderResponse)
	resp.ID = "54321"
	resp.POS = internal.POSBeta
	resp.TotalPrice = 17.79
	resp.Items = []*schema.OrderResponseItem{
		{
			Name:        "Falafel Wrap",
			Quantity:    1,
			Size:        "Regular",
			Extras:      []string{"Extra Chilli Sauce"},
			Ingredients: []string{},
			Price:       11.09,
		},
		{
			Name:        "Jam Donut",
			Quantity:    2,
			Size:        "Regular",
			Extras:      []string{},
			Ingredients: []string{},
			Price:       2.7,
		},
		{
			Name:        "Bottled Water",
			Quantity:    1,
			Size:        "Regular",
			Extras:      []string{},
			Ingredients: []string{},
			Price:       4,
		},
	}
	return
}
