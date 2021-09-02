package beta

import (
	"context"
	"fmt"
	"net/http"

	"github.com/flypay/engineering-test/pkg/api/apiHandler"
	"github.com/flypay/engineering-test/pkg/internal"
	"github.com/flypay/engineering-test/pkg/service"
)

type getMenuHandler struct {
	clientBaseURL string
}

func NewGetMenu() apiHandler.Handler {
	return &getMenuHandler{
		clientBaseURL: internal.BetaClientBaseURL,
	}
}

func (h *getMenuHandler) URL() string {
	return internal.MenuBeta
}

func (h *getMenuHandler) Methods() []string {
	return []string{http.MethodGet}
}

func (h *getMenuHandler) ParseArgs(r *http.Request) (*http.Request, error) {
	ctx := context.WithValue(r.Context(), internal.MethodContextKey, http.MethodGet)
	ctx = context.WithValue(ctx, internal.POSAddressContextKey,
		fmt.Sprintf("http://%s"+"%s", h.clientBaseURL, internal.Menu))
	r = r.WithContext(ctx)
	return r, nil
}

func (h *getMenuHandler) Process(r *http.Request) *http.Response {
	ctx := r.Context()
	method := ctx.Value(internal.MethodContextKey).(string)
	destination := ctx.Value(internal.POSAddressContextKey).(string)

	//if resp.StatusCode == http.StatusOK {
	//	return resp
	//}else {
	//	return &http.Response{
	//		StatusCode: resp.StatusCode,
	//	}
	//}
	return service.RequestClient(method, destination, nil)
}
