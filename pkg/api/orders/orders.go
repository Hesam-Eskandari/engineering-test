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
	"github.com/flypay/engineering-test/pkg/internal/schema"
	"github.com/flypay/engineering-test/pkg/service"
)

type setOrderHandler struct {
	alphaClientBaseURL string
	betaClientBaseURL  string
}

func NewSetOrder() apiHandler.Handler {
	return &setOrderHandler{
		alphaClientBaseURL: internal.AlphaClientBaseURL,
		betaClientBaseURL:  internal.BetaClientBaseURL,
	}
}

func (h *setOrderHandler) URL() string {
	return internal.Orders
}

func (h *setOrderHandler) Methods() []string {
	return []string{http.MethodPost}
}

func (h *setOrderHandler) ParseArgs(r *http.Request) (*http.Request, error) {
	requestBody := new(schema.OrderRequest)
	if err := service.DecodeReqRespBody(r.Body, requestBody); err != nil {
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
			fmt.Sprintf("http://%s"+"%s", h.alphaClientBaseURL, internal.Orders))
		ctx = context.WithValue(ctx, internal.BodyContextKey, getAlphaReqBody(requestBody))
	case "beta":
		ctx = context.WithValue(ctx, internal.POSAddressContextKey,
			fmt.Sprintf("http://%s"+"%s", h.betaClientBaseURL, internal.OrdersCreateBeta))
		body, menu, err := getBetaReqBody(requestBody)
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

func (h *setOrderHandler) Process(r *http.Request) *http.Response {
	ctx := r.Context()
	var err error
	method := ctx.Value(internal.MethodContextKey).(string)
	destination := ctx.Value(internal.POSAddressContextKey).(string)
	body := ctx.Value(internal.BodyContextKey).(io.ReadCloser)
	reqBody := ctx.Value(internal.RequestContextKey).(*schema.OrderRequest)
	resp := service.RequestPOSClient(method, destination, body)
	log.Printf("server %v responded with status: %v, order id: %v", reqBody.POS, resp.Status, reqBody.ID)
	if resp.StatusCode == http.StatusOK {
		unifiedBody := new(schema.OrderResponse) // desired response body to return
		switch reqBody.POS {
		case internal.POSAlpha:

			if err = service.BuildRespFromAlphaOrder(reqBody, unifiedBody); err != nil {
				resp.StatusCode = http.StatusInternalServerError
			}
		case internal.POSBeta:
			menu := ctx.Value(internal.MenuContextKey).(*schema.BetaMenu)
			if err = service.BuildRespFromBetaOrder(reqBody, menu, unifiedBody); err != nil {
				resp.StatusCode = http.StatusInternalServerError
			}

		}
		log.Printf("Order submitted. orderID: %s, pos: %s, "+
			"total price: $%.2f", reqBody.ID, reqBody.POS, unifiedBody.TotalPrice)
		resp.Body, err = service.EncodeReqRespBody(unifiedBody)
		if err != nil {
			log.Printf("failed encoding: error: %s", err.Error())
			resp.StatusCode = http.StatusInternalServerError
		}
	}
	return resp
}

func getAlphaReqBody(body *schema.OrderRequest) io.ReadCloser {
	var alphaReqProducts []schema.AlphaReqProduct
	for _, item := range body.Items {
		alphaReqProducts = append(alphaReqProducts,
			schema.NewAlphaReqProduct(item.ID, item.Size, item.Ingredients, item.Quantity))
	}

	alphaBody, _ := service.EncodeReqRespBody(&schema.AlphaReqBody{OrderId: body.ID, Products: alphaReqProducts})

	return alphaBody
}

func getBetaReqBody(body *schema.OrderRequest) (io.ReadCloser, *schema.BetaMenu, error) {
	// fetch beta menu to find corresponding category IDs for item IDs
	menu := new(schema.BetaMenu)
	err := service.GetBetaMenu(http.MethodGet,
		fmt.Sprintf("http://%s"+"%s", internal.BetaClientBaseURL, internal.Menu), menu)
	if err != nil {
		log.Printf("failed getting beta menu, err: %s", err.Error())
		return nil, menu, err
	}
	var betaReqItems []schema.BetaReqItem
	for _, item := range body.Items {
		var categoryId string
		for catId, cat := range menu.Categories {
			for itemId := range cat.Items {
				if itemId == item.ID {
					categoryId = catId
				}
			}
		}
		betaReqItems = append(betaReqItems, schema.NewBetaReqItem(categoryId, item.ID,
			item.Quantity, append(item.Extras, item.Ingredients...)))
	}
	betaBody, _ := service.EncodeReqRespBody(&schema.BetaReqBody{Id: body.ID, Items: betaReqItems})
	return betaBody, menu, nil
}

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
