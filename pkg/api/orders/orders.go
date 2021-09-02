package orders

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
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
	fmt.Println("Parsing Args")
	requestBody := &schema.OrderRequest{}
	if err := DecodeReqRespBody(r.Body, requestBody); err != nil {
		return r, fmt.Errorf("[orders] error decoding request body, error: %s", err.Error())
	}
	ctx := context.WithValue(r.Context(), internal.MethodContextKey, http.MethodPost)
	switch requestBody.POS {
	case "alpha":
		ctx = context.WithValue(ctx, internal.POSAddressContextKey,
			fmt.Sprintf("http://%s"+"%s", h.alphaClientBaseURL, internal.Orders))
		ctx = context.WithValue(ctx, internal.BodyContextKey, getAlphaBody(requestBody))
	case "beta":
		ctx = context.WithValue(ctx, internal.POSAddressContextKey,
			fmt.Sprintf("http://%s"+"%s", h.betaClientBaseURL, internal.OrdersCreateBeta))
		ctx = context.WithValue(ctx, internal.BodyContextKey, getBetaBody(requestBody))
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
	return service.RequestClient(method, destination, body)
}

// DecodeReqRespBody decodes request or response arguments from request/ response body
func DecodeReqRespBody(body io.Reader, v interface{}) error {
	if body == nil {
		return errors.New("error decoding: body is nil")
	}
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("error decoding:%s", err.Error())
	}
	return nil
}

func EncodeReqBody(body interface{}) (io.ReadCloser, error) {
	if body == nil {
		return nil, errors.New("error decoding: request body is nil")
	}
	b, _ := json.Marshal(body)
	return ioutil.NopCloser(bytes.NewReader(b)), nil

}

func getAlphaBody(body *schema.OrderRequest) io.ReadCloser {
	var alphaReqProducts []schema.AlphaReqProduct
	fmt.Println("reached 1")
	for _, item := range body.Items {
		if item.Size == "" {
			item.Size = "8001"
		}
		alphaReqProducts = append(alphaReqProducts,
			schema.NewAlphaReqProduct(item.ID, item.Size, item.Ingredients, item.Quantity))
	}

	alphaBody, _ := EncodeReqBody(&schema.AlphaReqBody{OrderId: body.ID, Products: alphaReqProducts})

	return alphaBody
}

func getBetaBody(body *schema.OrderRequest) io.ReadCloser {
	// fetch beta menu to find corresponding category IDs for item IDs
	resp := service.RequestClient(http.MethodGet,
		fmt.Sprintf("http://%s"+"%s", internal.BetaClientBaseURL, internal.Menu), nil)

	menu := &schema.BetaMenu{}
	if err := DecodeReqRespBody(resp.Body, menu); err != nil {
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
		betaReqItems = append(betaReqItems, schema.NewBetaReqItem(categoryId, item.ID, item.Quantity, item.Extras))
	}
	betaBody, _ := EncodeReqBody(&schema.BetaReqBody{Id: body.ID, Items: betaReqItems})
	return betaBody
}
