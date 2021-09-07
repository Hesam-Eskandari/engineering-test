package orders

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/flypay/engineering-test/pkg/api/apiHandler"
	"github.com/flypay/engineering-test/pkg/internal"
	"github.com/flypay/engineering-test/pkg/repository"
	"github.com/flypay/engineering-test/pkg/schema"
	"github.com/flypay/engineering-test/pkg/service"
)

// setOrderHandler handles the request of setting orders
type setOrderHandler struct {
	service service.Service
	repo    repository.Repository
}

// NewSetOrder returns an instance of setOrderHandler
func NewSetOrder() apiHandler.Handler {
	return &setOrderHandler{
		service: service.NewServiceImpl(),
		repo:    repository.NewRepositoryImpl(),
	}
}

// NewSetMockOrder returns an instance of setOrderHandler
func NewSetMockOrder() apiHandler.Handler {
	return &setOrderHandler{
		service: service.NewServiceMock(),
		repo:    repository.NewRepositoryImpl(),
	}
}

// URL returns the request URL to this handler
func (h *setOrderHandler) URL() string {
	return internal.Orders
}

// Methods returns the bounding HTTP methods
func (h *setOrderHandler) Methods() []string {
	return []string{http.MethodPost}
}

// ParseArgs parses request arguments in and generates context keys
func (h *setOrderHandler) ParseArgs(r *http.Request) (*http.Request, error) {
	requestBody := new(schema.OrderRequest)
	if err := h.repo.DecodeReqRespBody(r.Body, requestBody); err != nil {
		return r, fmt.Errorf("[orders] error decoding request body, error: %s", err.Error())
	}
	if err := validateOrder(requestBody); err != nil {
		return r, fmt.Errorf("[orders] failed validating order. err: %s", err.Error())
	}
	ctx := context.WithValue(r.Context(), internal.MethodContextKey, http.MethodPost)
	ctx = context.WithValue(ctx, internal.RequestContextKey, requestBody)
	switch requestBody.POS {
	case "alpha":
		ctx = context.WithValue(ctx, internal.POSAddressContextKey,
			fmt.Sprintf("http://%s"+"%s", internal.AlphaClientBaseURL, internal.Orders))
		body, menu, err := h.service.GetAlphaReqBody(requestBody)
		if err != nil {
			return r, fmt.Errorf("[orders] pos %s is not available. err: %s", requestBody.POS, err.Error())
		}
		ctx = context.WithValue(ctx, internal.BodyContextKey, body)
		ctx = context.WithValue(ctx, internal.MenuContextKey, menu)
	case "beta":
		ctx = context.WithValue(ctx, internal.POSAddressContextKey,
			fmt.Sprintf("http://%s"+"%s", internal.BetaClientBaseURL, internal.OrdersCreateBeta))
		body, menu, err := h.service.GetBetaReqBody(requestBody)
		if err != nil {
			return r, fmt.Errorf("[orders] pos %s is not available. err: %s", requestBody.POS, err.Error())
		}
		ctx = context.WithValue(ctx, internal.BodyContextKey, body)
		ctx = context.WithValue(ctx, internal.MenuContextKey, menu)
	default:
		return r, fmt.Errorf("[orders] pos %s is not valid", requestBody.POS)
	}
	r = r.WithContext(ctx)
	return r, nil
}

// Process processes the request
func (h *setOrderHandler) Process(r *http.Request) *http.Response {
	ctx := r.Context()
	var err error
	method := ctx.Value(internal.MethodContextKey).(string)
	destination := ctx.Value(internal.POSAddressContextKey).(string)
	body := ctx.Value(internal.BodyContextKey).(io.ReadCloser)
	reqBody := ctx.Value(internal.RequestContextKey).(*schema.OrderRequest)
	resp := h.service.RequestPOSClient(method, destination, body)
	log.Printf("server %v responded with status: %v, order id: %v", reqBody.POS, resp.Status, reqBody.ID)
	if resp.StatusCode == http.StatusOK {
		unifiedBody := new(schema.OrderResponse) // desired response body to return
		unifiedMenu := new(schema.Menu)          // unified menu structure
		switch reqBody.POS {
		case internal.POSAlpha:
			menu := ctx.Value(internal.MenuContextKey).(*schema.AlphaMenu)
			h.repo.PopulateUnifiedMenuFromAlphaMenu(menu, unifiedMenu)
			h.repo.PopulateUnifiedOrderRespBody(reqBody, unifiedMenu, unifiedBody)
		case internal.POSBeta:
			menu := ctx.Value(internal.MenuContextKey).(*schema.BetaMenu)
			h.repo.PopulateUnifiedMenuFromBetaMenu(menu, unifiedMenu)
			h.repo.PopulateUnifiedOrderRespBody(reqBody, unifiedMenu, unifiedBody)

		}
		log.Printf("Order submitted. orderID: %s, pos: %s, "+
			"total price: $%.2f", reqBody.ID, reqBody.POS, unifiedBody.TotalPrice)
		resp.Body, err = h.repo.EncodeReqRespBody(unifiedBody)
		if err != nil {
			log.Printf("failed encoding: error: %s", err.Error())
			resp.StatusCode = http.StatusInternalServerError
		}
	}
	return resp
}

// validateOrder validates request payload
func validateOrder(order *schema.OrderRequest) error {
	if !internal.ValidatePoses(order.POS) {
		return errors.New(fmt.Sprintf("POS %s is not valid", order.POS))
	}
	if order.ID == "" {
		return errors.New("order id cannot be empty")
	}
	if order.Items == nil || len(order.Items) == 0 {
		return errors.New("must choose at least one item")
	}
	for index, item := range order.Items {
		if item.Quantity == 0 {
			order.Items[index].Quantity = 1
		}
	}
	return nil
}
