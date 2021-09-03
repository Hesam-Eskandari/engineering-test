package service

import (
	"fmt"
	"log"
	"math"
	"net/http"

	"github.com/flypay/engineering-test/pkg/internal"
	"github.com/flypay/engineering-test/pkg/internal/schema"
)

func BuildRespFromAlpha(resp *http.Response, reqBody *schema.OrderRequest) *http.Response {

	var body schema.AlphaRespBody        // response body from alpha pos
	var unifiedBody schema.OrderResponse // desired response body to return
	var err error

	populateUnifiedRespBodyForAlpha(reqBody, &unifiedBody)

	if resp.Body, err = EncodeReqRespBody(unifiedBody); err != nil {
		log.Fatalf("failed encoding: incompatible schema: %v", body)
	}
	return resp
}

func BuildRespFromBeta(resp *http.Response, reqBody *schema.OrderRequest) *http.Response {
	var body schema.BetaRespBody
	var err error
	if err = DecodeReqRespBody(resp.Body, &body); err != nil {
		log.Fatalf("failed decoding:, incompatible response body: %v", resp.Body)
	}
	if resp.Body, err = EncodeReqRespBody(body); err != nil {
		log.Fatalf("failed encoding: incompatible schema: %v", body)
	}
	return resp
}

func populateUnifiedRespBodyForAlpha(reqBody *schema.OrderRequest, unifiedBody *schema.OrderResponse) {
	var err error

	// get products menu from alpha pos
	productsResp := RequestPOSClient(http.MethodGet,
		fmt.Sprintf("http://%s"+"%s", internal.AlphaClientBaseURL, internal.MenuProductsAlpha), nil)
	var products schema.AlphaProductsMenu
	if err = DecodeReqRespBody(productsResp.Body, &products); err != nil {
		log.Fatalf("failed decoding:, incompatible response body: %v", err)
	}
	// get ingredients menu from alpha pos
	ingredientsResp := RequestPOSClient(http.MethodGet,
		fmt.Sprintf("http://%s"+"%s", internal.AlphaClientBaseURL, internal.MenuIngredientsAlpha), nil)
	var ingredients schema.AlphaIngredientsMenu
	if err = DecodeReqRespBody(ingredientsResp.Body, &ingredients); err != nil {
		log.Fatalf("failed decoding:, incompatible response body: %v", ingredientsResp.Body)
	}
	// ingredientItem represents a structure that includes ingredient name and price
	type ingredientItem struct {
		name  string
		price float32
	}

	// newIngredientItem is a function literal to construct a new ingredientItem instance
	newIngredientItem := func(name string, price float32) *ingredientItem {
		return &ingredientItem{name, price}
	}

	// productSizeItem represents a structure that includes product size name and product price
	type productSizeItem struct {
		name  string
		price float32
	}

	// productItem represents a structure that includes product name, sizes, and an array of allowed extras
	type productItem struct {
		name          string
		sizes         map[string]*productSizeItem
		allowedExtras map[string]bool
	}

	// newProductItem is a function literal to construct a new productItem instance
	newProductItem := func(name string, sizesMap map[string]*productSizeItem, allowedExtras map[string]bool) *productItem {
		return &productItem{
			name:          name,
			sizes:         sizesMap,
			allowedExtras: allowedExtras,
		}
	}

	// ingredientsMap maps ingredient ids to ingredientItem (includes ingredient name and price)
	ingredientsMap := make(map[string]*ingredientItem)

	// productsMap maps product ids to productItem (includes product name, sizes, and an array of allowed extras)
	productsMap := make(map[string]*productItem)
	for _, ingredient := range ingredients.Ingredients {
		ingredientsMap[ingredient.IngredientID] = newIngredientItem(ingredient.Name, 0)
	}
	for _, product := range products.Products {
		allowedExtras := make(map[string]bool)
		for _, extra := range product.Extras {
			allowedExtras[extra.IngredientID] = true
			if _, ok := ingredientsMap[extra.IngredientID]; ok {
				ingredientsMap[extra.IngredientID].price = extra.Price
			} else {
				log.Printf("unknown ingredient id: %v, pos: %v", extra.IngredientID, reqBody.POS)
			}
		}
		sizesMap := make(map[string]*productSizeItem)
		for _, size := range product.Sizes {
			sizesMap[size.SizeID] = &productSizeItem{size.Name, size.Price}
		}
		productsMap[product.ProductID] = newProductItem(product.Name, sizesMap, allowedExtras)
	}

	// populate unifiedBody
	var totalPrice float32
	unifiedBody.ID = reqBody.ID
	unifiedBody.POS = reqBody.POS
	for _, item := range reqBody.Items {
		var price float32
		extras := make([]string, 0, len(item.Extras))
		for _, extraID := range item.Ingredients {
			fmt.Println("extra", extraID)
			if _, ok := productsMap[item.ID].allowedExtras[extraID]; ok {
				extras = append(extras, ingredientsMap[extraID].name)
				price += ingredientsMap[extraID].price
			}
		}
		price += float32(item.Quantity) * productsMap[item.ID].sizes[item.Size].price
		unifiedBody.Items = append(unifiedBody.Items, schema.OrderResponseItem{
			Name:     productsMap[item.ID].name,
			Quantity: item.Quantity,
			Size:     productsMap[item.ID].sizes[item.Size].name,
			Extras:   extras,
			Price:    float32(math.Round(float64(price*1000)) / 1000),
		})
		totalPrice += price
	}
	unifiedBody.TotalPrice = float32(math.Round(float64(totalPrice*1000)) / 1000)
	return
}
