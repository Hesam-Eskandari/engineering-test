package orders

import (
	"net/http"
	"testing"

	"github.com/flypay/engineering-test/pkg/repository"
	"github.com/flypay/engineering-test/pkg/schema"
)

func TestAlphaOrder_UnitTesting(t *testing.T) {
	repo := repository.NewRepositoryImpl()
	expectedOrderResp := schema.NewUnifiedRespPopulatedAlphaOrderRespMock()
	requestBody, _ := repo.EncodeReqRespBody(schema.NewUnifiedOrderRequestPopulatedOrderAlphaMock())
	expectedOrderResp.POS = "alpha"
	request := new(http.Request)
	request.Body = requestBody
	orderHandler := NewSetMockOrder()
	request, _ = orderHandler.ParseArgs(request)
	resp := orderHandler.Process(request)
	orderBody := new(schema.OrderResponse)
	_ = repo.DecodeReqRespBody(resp.Body, orderBody)
	t.Run("test alpha order", func(t *testing.T) {
		assertEqualOrders(t, orderBody, expectedOrderResp)
	})
}

func TestBetaOrder_UnitTesting(t *testing.T) {
	repo := repository.NewRepositoryImpl()
	expectedOrderResp := schema.NewUnifiedRespPopulatedBetaOrderRespMock()
	requestBody, _ := repo.EncodeReqRespBody(schema.NewUnifiedOrderRequestPopulatedOrderBetaMock())
	expectedOrderResp.POS = "beta"
	request := new(http.Request)
	request.Body = requestBody
	orderHandler := NewSetMockOrder()
	request, _ = orderHandler.ParseArgs(request)
	resp := orderHandler.Process(request)
	orderBody := new(schema.OrderResponse)
	_ = repo.DecodeReqRespBody(resp.Body, orderBody)
	t.Run("test beta order", func(t *testing.T) {
		assertEqualOrders(t, orderBody, expectedOrderResp)
	})
}
