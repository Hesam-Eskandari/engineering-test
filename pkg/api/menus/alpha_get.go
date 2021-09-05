package menus

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/flypay/engineering-test/pkg/api/apiHandler"
	"github.com/flypay/engineering-test/pkg/internal"
	"github.com/flypay/engineering-test/pkg/internal/schema"
	"github.com/flypay/engineering-test/pkg/service"
)

// getAlphaMenuHandler handles the request of getting alpha's menu
type getAlphaMenuHandler struct {
	clientBaseURL string
}

// NewGetAlphaMenu returns an instance of getAlphaMenuHandler
func NewGetAlphaMenu() apiHandler.Handler {
	return &getAlphaMenuHandler{
		clientBaseURL: internal.AlphaClientBaseURL,
	}
}

// URL returns the request URL to this handler
func (h *getAlphaMenuHandler) URL() string {
	return internal.BasePathAlpha + internal.Menu
}

// Methods returns the bounding HTTP methods
func (h *getAlphaMenuHandler) Methods() []string {
	return []string{http.MethodGet}
}

// ParseArgs parses request arguments in and generates context keys
func (h *getAlphaMenuHandler) ParseArgs(r *http.Request) (*http.Request, error) {
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

// Process processes the request
func (h *getAlphaMenuHandler) Process(r *http.Request) *http.Response {
	resp := new(http.Response)
	ctx := r.Context()
	method := ctx.Value(internal.MethodContextKey).(string)
	alphaAddress := ctx.Value(internal.POSAddressContextKey).(*schema.AlphaMenuAddress)

	alphaMenu := new(schema.AlphaMenu)
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
