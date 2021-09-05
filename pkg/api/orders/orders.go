package orders

import (
	"context"
	"fmt"
	"io"
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
	requestBody := &schema.OrderRequest{}
	if err := service.DecodeReqRespBody(r.Body, requestBody); err != nil {
		return r, fmt.Errorf("[orders] error decoding request body, error: %s", err.Error())
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
		body, menu := getBetaReqBody(requestBody)
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
	method := ctx.Value(internal.MethodContextKey).(string)
	destination := ctx.Value(internal.POSAddressContextKey).(string)
	body := ctx.Value(internal.BodyContextKey).(io.ReadCloser)
	reqBody := ctx.Value(internal.RequestContextKey).(*schema.OrderRequest)
	resp := service.RequestPOSClient(method, destination, body)

	if resp.StatusCode == http.StatusOK {
		switch reqBody.POS {
		case internal.POSAlpha:
			resp = service.BuildRespFromAlphaOrder(resp, reqBody)
		case internal.POSBeta:
			menu := ctx.Value(internal.MenuContextKey).(*schema.BetaMenu)
			resp = service.BuildRespFromBetaOrder(resp, reqBody, menu)
		}
	}
	return resp
}

func getAlphaReqBody(body *schema.OrderRequest) io.ReadCloser {
	var alphaReqProducts []schema.AlphaReqProduct
	for _, item := range body.Items {
		if item.Size == "" {
			item.Size = "8001"
		}
		alphaReqProducts = append(alphaReqProducts,
			schema.NewAlphaReqProduct(item.ID, item.Size, item.Ingredients, item.Quantity))
	}

	alphaBody, _ := service.EncodeReqRespBody(&schema.AlphaReqBody{OrderId: body.ID, Products: alphaReqProducts})

	return alphaBody
}

func getBetaReqBody(body *schema.OrderRequest) (io.ReadCloser, *schema.BetaMenu) {
	// fetch beta menu to find corresponding category IDs for item IDs
	resp := service.RequestPOSClient(http.MethodGet,
		fmt.Sprintf("http://%s"+"%s", internal.BetaClientBaseURL, internal.Menu), nil)

	menu := schema.BetaMenu{}
	if err := service.DecodeReqRespBody(resp.Body, &menu); err != nil {
		panic(err.Error())
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
	return betaBody, &menu
}
