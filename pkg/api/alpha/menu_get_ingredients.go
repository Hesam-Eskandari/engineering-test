package alpha

import (
	"context"
	"fmt"
	"net/http"

	"github.com/flypay/engineering-test/pkg/api/apiHandler"
	"github.com/flypay/engineering-test/pkg/internal"
	"github.com/flypay/engineering-test/pkg/service"
)

type getMenuIngredientsHandler struct {
	clientBaseURL string
}

func NewGetMenuIngredients() apiHandler.Handler {
	return &getMenuIngredientsHandler{
		clientBaseURL: internal.AlphaClientBaseURL,
	}
}

func (h *getMenuIngredientsHandler) URL() string {
	return internal.BasePathAlpha + internal.MenuIngredientsAlpha
}

func (h *getMenuIngredientsHandler) Methods() []string {
	return []string{http.MethodGet}
}

func (h *getMenuIngredientsHandler) ParseArgs(r *http.Request) (*http.Request, error) {
	ctx := context.WithValue(r.Context(), internal.MethodContextKey, http.MethodGet)
	ctx = context.WithValue(ctx, internal.POSAddressContextKey,
		fmt.Sprintf("http://%s"+"%s", h.clientBaseURL, internal.MenuIngredientsAlpha))
	r = r.WithContext(ctx)
	return r, nil
}

func (h *getMenuIngredientsHandler) Process(r *http.Request) *http.Response {
	ctx := r.Context()
	method := ctx.Value(internal.MethodContextKey).(string)
	destination := ctx.Value(internal.POSAddressContextKey).(string)
	return service.RequestClient(method, destination, nil)
}
