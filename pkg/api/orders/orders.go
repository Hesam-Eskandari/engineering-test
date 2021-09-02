package orders

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/flypay/engineering-test/pkg/api/apiHandler"
	"github.com/flypay/engineering-test/pkg/internal"
	"github.com/flypay/engineering-test/pkg/service"
	"io"
	"io/ioutil"
	"net/http"
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
	requestBody := &internal.OrderRequest{}
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

// DecodeReqRespBody decodes request arguments from request body
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

type alphaReqBody struct {
	OrderId  string            `json:"orderId"`
	Products []alphaReqProduct `json:"products"`
}

type alphaReqProduct struct {
	ProductId     string   `json:"productId"`
	SizeId        string   `json:"sizeId"`
	IngredientIds []string `json:"ingredientIds"`
	Quantity      int      `json:"quantity"`
}

func newAlphaReqProduct(productId, sizeId string, ingredientIds []string, quantity int) alphaReqProduct {
	return alphaReqProduct{
		ProductId:     productId,
		SizeId:        sizeId,
		IngredientIds: ingredientIds,
		Quantity:      quantity,
	}
}

func getAlphaBody(body *internal.OrderRequest) io.ReadCloser {
	var alphaReqProducts []alphaReqProduct
	fmt.Println("reached 1")
	for _, item := range body.Items {
		if item.Size == "" {
			item.Size = "8001"
		}
		alphaReqProducts = append(alphaReqProducts, newAlphaReqProduct(item.ID, item.Size, item.Ingredients, item.Quantity))
	}

	alphaBody, _ := EncodeReqBody(&alphaReqBody{OrderId: body.ID, Products: alphaReqProducts})

	return alphaBody
}

type betaReqBody struct {
	Id    string        `json:"Id"`
	Items []betaReqItem `json:"Items"`
}

type betaReqItem struct {
	CategoryId string   `json:"CategoryId"`
	ItemId     string   `json:"ItemId"`
	Quantity   int      `json:"Quantity"`
	AddOns     []string `json:"AddOns"`
}

func newBetaReqItem(categoryId, itemId string, quantity int, addOns []string) betaReqItem {
	return betaReqItem{
		CategoryId: categoryId,
		ItemId:     itemId,
		Quantity:   quantity,
		AddOns:     addOns,
	}
}

type betaMenu struct {
	Categories map[string]struct {
		Name  string `json:"Name"`
		Items map[string]struct {
			Name     string  `json:"Name"`
			Price    float32 `json:"Price"`
			Quantity int     `json:"Quantity"`
			AddOns   []struct {
				ID    string  `json:"Id"`
				Name  string  `json:"Name"`
				Price float32 `json:"Price"`
			} `json:"AddOns"`
		} `json:"Items"`
	} `json:"Categories"`
}

func getBetaBody(body *internal.OrderRequest) io.ReadCloser {
	resp := service.RequestClient(http.MethodGet,
		fmt.Sprintf("http://%s"+"%s", internal.BetaClientBaseURL, internal.Menu), nil)
	fmt.Println(resp.StatusCode)
	menu := &betaMenu{}
	if err := DecodeReqRespBody(resp.Body, menu); err != nil {
		panic(err.Error())
	}
	var betaReqItems []betaReqItem
	for _, item := range body.Items {
		var categoryId string
		for catId, cat := range menu.Categories {
			for itemId := range cat.Items {
				if itemId == item.ID {
					categoryId = catId
				}
			}
		}
		betaReqItems = append(betaReqItems, newBetaReqItem(categoryId, item.ID, item.Quantity, item.Extras))
	}
	betaBody, _ := EncodeReqBody(&betaReqBody{Id: body.ID, Items: betaReqItems})
	return betaBody
}
