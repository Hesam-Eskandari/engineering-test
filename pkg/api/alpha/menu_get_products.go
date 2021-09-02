package alpha

import (
	"context"
	"fmt"
	"net/http"

	"github.com/flypay/engineering-test/pkg/api/apiHandler"
	"github.com/flypay/engineering-test/pkg/internal"
	"github.com/flypay/engineering-test/pkg/service"
)

type getMenuProductsHandler struct {
	clientBaseURL string
}

func NewGetMenuProductsHandler() apiHandler.Handler {
	return &getMenuProductsHandler{
		clientBaseURL: internal.AlphaClientBaseURL,
	}
}

func (h *getMenuProductsHandler) URL() string {
	return internal.BasePathAlpha + internal.MenuProductsAlpha
}

func (h *getMenuProductsHandler) Methods() []string {
	return []string{http.MethodGet}
}

func (h *getMenuProductsHandler) ParseArgs(r *http.Request) (*http.Request, error) {
	ctx := context.WithValue(r.Context(), internal.MethodContextKey, http.MethodGet)
	ctx = context.WithValue(ctx, internal.POSAddressContextKey,
		fmt.Sprintf("http://%s"+"%s", h.clientBaseURL, internal.MenuProductsAlpha))
	r = r.WithContext(ctx)
	return r, nil
}

func (h *getMenuProductsHandler) Process(r *http.Request) *http.Response {
	ctx := r.Context()
	method := ctx.Value(internal.MethodContextKey).(string)
	destination := ctx.Value(internal.POSAddressContextKey).(string)
	return service.RequestClient(method, destination, nil)
}
