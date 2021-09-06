package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"sort"

	"github.com/flypay/engineering-test/pkg/internal/schema"
)

// DecodeReqRespBody decodes request or response arguments from request/response body to a schema
func DecodeReqRespBody(body io.Reader, v interface{}) error {
	if body == nil {
		return errors.New("error decoding: body is nil")
	}
	if v == nil {
		return errors.New("error decoding: interface is nil")
	}
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("error decoding:%s", err.Error())
	}
	return nil
}

// EncodeReqRespBody encodes request or response arguments from a schema to request/response body
func EncodeReqRespBody(body interface{}) (io.ReadCloser, error) {
	if body == nil {
		return nil, errors.New("error encoding: request body is nil")
	}
	b, _ := json.Marshal(body)
	return ioutil.NopCloser(bytes.NewReader(b)), nil

}

// SortUnifiedMenu sorts any array in unified menu alphabetically by name
func SortUnifiedMenu(menu *schema.Menu) {
	sort.Slice(menu.Categories, func(i, j int) bool { return menu.Categories[i].Name < menu.Categories[j].Name })
	for _, cat := range menu.Categories {
		for _, subCat := range cat.Subcategories {
			for _, item := range subCat.Items {
				sort.Slice(item.Extras, func(i, j int) bool { return item.Extras[i].Name < item.Extras[j].Name })
				sort.Slice(item.Ingredients, func(i, j int) bool { return item.Ingredients[i].Name < item.Ingredients[j].Name })
				sort.Slice(item.Sizes, func(i, j int) bool { return item.Sizes[i].Name < item.Sizes[j].Name })
			}
			sort.Slice(subCat.Items, func(i, j int) bool { return subCat.Items[i].Name < subCat.Items[j].Name })
		}
		sort.Slice(cat.Subcategories, func(i, j int) bool { return cat.Subcategories[i].Name < cat.Subcategories[j].Name })
	}
}

// SortUnifiedResponse sorts all arrays in unified response body alphabetically by name
func SortUnifiedResponse(resp *schema.OrderResponse) {
	sort.Slice(resp.Items, func(i, j int) bool { return resp.Items[i].Name < resp.Items[j].Name })
	for _, item := range resp.Items {
		sort.Slice(item.Extras, func(i, j int) bool { return item.Extras[i] < item.Extras[j] })
		sort.Slice(item.Ingredients, func(i, j int) bool { return item.Ingredients[i] < item.Ingredients[j] })
	}
}
