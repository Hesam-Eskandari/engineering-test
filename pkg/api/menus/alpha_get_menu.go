package menus

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/flypay/engineering-test/pkg/api/apiHandler"
	"github.com/flypay/engineering-test/pkg/internal"
	"github.com/flypay/engineering-test/pkg/repository"
	"github.com/flypay/engineering-test/pkg/schema"
	"github.com/flypay/engineering-test/pkg/service"
)

// getAlphaMenuHandler handles the request of getting alpha's menu
type getAlphaMenuHandler struct {
	service service.Service
	repo    repository.Repository
}

// NewGetAlphaMenu returns an instance of getAlphaMenuHandler
func NewGetAlphaMenu() apiHandler.Handler {
	return &getAlphaMenuHandler{
		service: service.NewServiceImpl(),
		repo:    repository.NewRepositoryImpl(),
	}
}

// NewGetMockAlphaMenu returns an instance of getAlphaMenuHandler
func NewGetMockAlphaMenu() apiHandler.Handler {
	return &getAlphaMenuHandler{
		service: service.NewServiceMock(),
		repo:    repository.NewRepositoryImpl(),
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
	ctx = context.WithValue(ctx, internal.POSAddressContextKey, schema.NewAlphaMenuAddress())
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
	if err := h.service.GetAlphaMenu(method, alphaAddress, alphaMenu); err != nil {
		log.Printf("error getting menu from alpha pos. err: %s", err.Error())
		resp.StatusCode = http.StatusInternalServerError
		resp.Status = fmt.Sprintf("failed getting menu from alpha pos")
		return resp
	}
	fmt.Println("alphaMenu", alphaMenu)
	unifiedMenu := new(schema.Menu)
	h.repo.PopulateUnifiedMenuFromAlphaMenu(alphaMenu, unifiedMenu)
	h.repo.SortUnifiedMenu(unifiedMenu)
	body, err := h.repo.EncodeReqRespBody(unifiedMenu)
	if err != nil {
		fmt.Printf("err reached, err: %v", err.Error())
	}
	resp.Body = body
	resp.StatusCode = http.StatusOK
	return resp
}
