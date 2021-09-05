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

// getBetaMenuHandler handles the request of getting alpha's menu
type getBetaMenuHandler struct {
	clientBaseURL string
}

// NewGetBetaMenu returns an instance of getBetaMenuHandler
func NewGetBetaMenu() apiHandler.Handler {
	return &getBetaMenuHandler{
		clientBaseURL: internal.BetaClientBaseURL,
	}
}

// URL returns the request URL to this handler
func (h *getBetaMenuHandler) URL() string {
	return internal.MenuBeta
}

// Methods returns the bounding HTTP methods
func (h *getBetaMenuHandler) Methods() []string {
	return []string{http.MethodGet}
}

// ParseArgs parses request arguments in and generates context keys
func (h *getBetaMenuHandler) ParseArgs(r *http.Request) (*http.Request, error) {
	// no argument in the body of the request we accept
	// lets just create context keys
	ctx := context.WithValue(r.Context(), internal.MethodContextKey, http.MethodGet)
	ctx = context.WithValue(ctx, internal.POSAddressContextKey,
		fmt.Sprintf("http://%s"+"%s", h.clientBaseURL, internal.Menu))
	r = r.WithContext(ctx)
	return r, nil
}

// Process processes the request
func (h *getBetaMenuHandler) Process(r *http.Request) *http.Response {
	resp := new(http.Response)
	ctx := r.Context()
	method := ctx.Value(internal.MethodContextKey).(string)
	betaAddress := ctx.Value(internal.POSAddressContextKey).(string)

	betaMenu := new(schema.BetaMenu)
	if err := service.GetBetaMenu(method, betaAddress, betaMenu); err != nil {
		log.Printf("error getting menu from beta pos. err: %s", err.Error())
		resp.StatusCode = http.StatusInternalServerError
		resp.Status = fmt.Sprintf("failed getting menu from beta pos")
		return resp
	}
	unifiedMenu := new(schema.Menu)
	service.PopulateUnifiedMenuFromBetaMenu(betaMenu, unifiedMenu)
	body, err := service.EncodeReqRespBody(unifiedMenu)
	if err != nil {
		fmt.Printf("err reached, err: %v", err.Error())
	}
	resp.Body = body
	resp.StatusCode = http.StatusOK
	return resp
}
