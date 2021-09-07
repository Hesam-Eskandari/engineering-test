package service

import (
	"io"
	"net/http"

	"github.com/flypay/engineering-test/pkg/schema"
)

type Service interface {
	// GetAlphaMenu calls alpha and receives its menus
	GetAlphaMenu(method string, destination *schema.AlphaMenuAddress, alphaMenu *schema.AlphaMenu) error
	// GetBetaMenu calls beta and receives its menu
	GetBetaMenu(method, address string, betaMenu *schema.BetaMenu) error
	// GetAlphaReqBody creates compatible request to get menu from alpha
	GetAlphaReqBody(body *schema.OrderRequest) (io.ReadCloser, *schema.AlphaMenu, error)
	// GetBetaReqBody creates compatible request to get menu from beta
	GetBetaReqBody(body *schema.OrderRequest) (io.ReadCloser, *schema.BetaMenu, error)
	// RequestPOSClient calls a given POS and returns its response
	RequestPOSClient(method, destination string, body io.ReadCloser) (resp *http.Response)
}
