package menus

import (
	"context"
	"fmt"
	"net/http"

	"github.com/flypay/engineering-test/pkg/api/apiHandler"
	"github.com/flypay/engineering-test/pkg/internal"
	"github.com/flypay/engineering-test/pkg/service"
)

type getBetaMenuHandler struct {
	clientBaseURL string
}

func NewGetBetaMenu() apiHandler.Handler {
	return &getBetaMenuHandler{
		clientBaseURL: internal.BetaClientBaseURL,
	}
}

func (h *getBetaMenuHandler) URL() string {
	return internal.MenuBeta
}

func (h *getBetaMenuHandler) Methods() []string {
	return []string{http.MethodGet}
}

func (h *getBetaMenuHandler) ParseArgs(r *http.Request) (*http.Request, error) {
	ctx := context.WithValue(r.Context(), internal.MethodContextKey, http.MethodGet)
	ctx = context.WithValue(ctx, internal.POSAddressContextKey,
		fmt.Sprintf("http://%s"+"%s", h.clientBaseURL, internal.Menu))
	r = r.WithContext(ctx)
	return r, nil
}

func (h *getBetaMenuHandler) Process(r *http.Request) *http.Response {
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
	return service.RequestPOSClient(method, destination, nil)
}
