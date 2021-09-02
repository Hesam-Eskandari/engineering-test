package alpha

import (
	"context"
	"fmt"
	"github.com/flypay/engineering-test/pkg/api/apiHandler"
	"github.com/flypay/engineering-test/pkg/internal"
	"github.com/flypay/engineering-test/pkg/service"
	"net/http"
	"path"
)

type getMenuCategoriesHandler struct {
	clientBaseURL string
}

func NewGetMenuCategories() apiHandler.Handler {
	return &getMenuCategoriesHandler{
		clientBaseURL: internal.AlphaClientBaseURL,
	}
}

func (h *getMenuCategoriesHandler) URL() string {
	return internal.BasePathAlpha + internal.MenuCategoriesAlpha
}

func (h *getMenuCategoriesHandler) Methods() []string {
	return []string{http.MethodGet}
}

func (h *getMenuCategoriesHandler) ParseArgs(r *http.Request) (*http.Request, error) {
	fmt.Println("Parsing Args")
	ctx := context.WithValue(r.Context(), internal.MethodContextKey, http.MethodGet)
	ctx = context.WithValue(ctx, internal.POSAddressContextKey,
		fmt.Sprintf("http://%s"+"%s", h.clientBaseURL, internal.MenuCategoriesAlpha))
	fmt.Println("path:", path.Base(r.URL.Path))
	//switch path.Base(r.URL.Path) {
	//case internal.MenuCategoriesAlpha:

	//case internal.MenuProductsAlpha:
	//	ctx = context.WithValue(ctx, internal.POSAddressContextKey,
	//		fmt.Sprintf("http://%s" + "%s", h.clientBaseURL, internal.MenuProductsAlpha))
	//case internal.MenuIngredientsAlpha:
	//	ctx = context.WithValue(ctx, internal.POSAddressContextKey,
	//		fmt.Sprintf("http://%s" + "%s", h.clientBaseURL, internal.MenuIngredientsAlpha))
	//default:
	//	return r, fmt.Errorf("")
	//}
	r = r.WithContext(ctx)
	return r, nil
}

func (h *getMenuCategoriesHandler) Process(r *http.Request) *http.Response {
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
