package repository

import (
	"io"

	"github.com/flypay/engineering-test/pkg/schema"
)

type Repository interface {
	// DecodeReqRespBody decodes request or response arguments from request/response body to a schema
	DecodeReqRespBody(body io.Reader, v interface{}) error
	// EncodeReqRespBody encodes request or response arguments from a schema to request/response body
	EncodeReqRespBody(body interface{}) (io.ReadCloser, error)
	// SortUnifiedMenu sorts any array in unified menu alphabetically by name
	SortUnifiedMenu(menu *schema.Menu)
	// SortUnifiedResponse sorts all arrays in unified response body alphabetically by name
	SortUnifiedResponse(resp *schema.OrderResponse)
	// PopulateUnifiedOrderRespBody populates the unified order response body with given order
	PopulateUnifiedOrderRespBody(reqBody *schema.OrderRequest, menu *schema.Menu, unifiedBody *schema.OrderResponse)
	// PopulateUnifiedMenuFromAlphaMenu populates the unified menu struct with alpha menus
	PopulateUnifiedMenuFromAlphaMenu(alphaMenu *schema.AlphaMenu, unifiedMenu *schema.Menu)
	// PopulateUnifiedMenuFromBetaMenu populates the unified menu struct with beta menu
	PopulateUnifiedMenuFromBetaMenu(betaMenu *schema.BetaMenu, unifiedMenu *schema.Menu)
}
