package alpha

import (
	"context"
	"fmt"
	"github.com/flypay/engineering-test/pkg/internal/schema"
	"log"
	"net/http"

	"github.com/flypay/engineering-test/pkg/api/apiHandler"
	"github.com/flypay/engineering-test/pkg/internal"
	"github.com/flypay/engineering-test/pkg/service"
)

type getMenuHandler struct {
	clientBaseURL string
}

func NewGetMenuCategories() apiHandler.Handler {
	return &getMenuHandler{
		clientBaseURL: internal.AlphaClientBaseURL,
	}
}

func (h *getMenuHandler) URL() string {
	return internal.BasePathAlpha + internal.Menu
}

func (h *getMenuHandler) Methods() []string {
	return []string{http.MethodGet}
}

// ParseArgs parses arguments in the request body and generates context keys
func (h *getMenuHandler) ParseArgs(r *http.Request) (*http.Request, error) {
	// no argument in the body of the request we accept
	// lets just create context keys
	ctx := context.WithValue(r.Context(), internal.MethodContextKey, http.MethodGet)
	ctx = context.WithValue(ctx, internal.POSAddressContextKey,
		&schema.AlphaMenuAddress{
			AlphaCategoriesAddress:  fmt.Sprintf("http://%s"+"%s", h.clientBaseURL, internal.MenuCategoriesAlpha),
			AlphaIngredientsAddress: fmt.Sprintf("http://%s"+"%s", h.clientBaseURL, internal.MenuIngredientsAlpha),
			AlphaProductsAddress:    fmt.Sprintf("http://%s"+"%s", h.clientBaseURL, internal.MenuProductsAlpha),
		})
	r = r.WithContext(ctx)
	return r, nil
}

func (h *getMenuHandler) Process(r *http.Request) *http.Response {
	resp := new(http.Response)
	ctx := r.Context()
	method := ctx.Value(internal.MethodContextKey).(string)
	alphaAddress := ctx.Value(internal.POSAddressContextKey).(*schema.AlphaMenuAddress)

	alphaMenu := schema.NewEmptyAlphaMenu()
	if err := service.GetAlphaMenu(method, alphaAddress, alphaMenu); err != nil {
		log.Printf("error getting menu from alpha pos. err: %s", err.Error())
		resp.StatusCode = http.StatusInternalServerError
		resp.Status = fmt.Sprintf("failed getting menu from alpha pos")
		return resp
	}
	unifiedMenu := schema.Menu{}

	service.PopulateUnifiedMenuFromAlphaMenu(alphaMenu, &unifiedMenu)
	body, err := service.EncodeReqRespBody(unifiedMenu)
	if err != nil {
		fmt.Printf("err reached, err: %v", err.Error())
	}
	resp.Body = body
	resp.StatusCode = http.StatusOK
	return resp
}
