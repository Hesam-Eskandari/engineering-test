package service

import (
	"io"
	"net/http"

	"github.com/flypay/engineering-test/pkg/repository"
	"github.com/flypay/engineering-test/pkg/schema"
)

type serviceMock struct {
	repo repository.Repository
}

func NewServiceMock() Service {
	return &serviceMock{
		repo: repository.NewRepositoryImpl(),
	}
}

// GetAlphaMenu calls alpha and receives its menus
func (s *serviceMock) GetAlphaMenu(method string, destination *schema.AlphaMenuAddress, alphaMenu *schema.AlphaMenu) error {
	alphaMenu.AlphaCategoriesMenu = schema.NewAlphaMenuMock().AlphaCategoriesMenu
	alphaMenu.AlphaIngredientsMenu = schema.NewAlphaMenuMock().AlphaIngredientsMenu
	alphaMenu.AlphaProductsMenu = schema.NewAlphaMenuMock().AlphaProductsMenu
	return nil
}

// GetBetaMenu calls beta and receives its menu
func (s *serviceMock) GetBetaMenu(method, address string, betaMenu *schema.BetaMenu) error {
	betaMenu.Categories = schema.NewBetaMenuMock().Categories
	return nil
}

// GetAlphaReqBody creates compatible request to get menu from alpha
func (s *serviceMock) GetAlphaReqBody(body *schema.OrderRequest) (io.ReadCloser, *schema.AlphaMenu, error) {
	alphaMenu := schema.NewAlphaMenuMock()
	resp, err := s.repo.EncodeReqRespBody(alphaMenu)
	return resp, alphaMenu, err
}

// GetBetaReqBody creates compatible request to get menu from beta
func (s *serviceMock) GetBetaReqBody(body *schema.OrderRequest) (io.ReadCloser, *schema.BetaMenu, error) {
	betaMenu := schema.NewBetaMenuMock()
	resp, err := s.repo.EncodeReqRespBody(betaMenu)
	return resp, betaMenu, err
}

// RequestPOSClient calls a given POS and returns its response
func (s *serviceMock) RequestPOSClient(method, destination string, body io.ReadCloser) (resp *http.Response) {
	resp = new(http.Response)
	resp.StatusCode = http.StatusOK
	return
}
