package service

import (
	"log"
	"net/http"

	"github.com/flypay/engineering-test/pkg/internal/schema"
)

func BuildRespFromAlphaOrder(resp *http.Response, reqBody *schema.OrderRequest) *http.Response {

	// response body from alpha pos
	var unifiedBody schema.OrderResponse // desired response body to return
	var err error

	populateUnifiedRespBodyForAlpha(reqBody, &unifiedBody)

	if resp.Body, err = EncodeReqRespBody(unifiedBody); err != nil {
		log.Fatalf("failed encoding: incompatible schema: %v", unifiedBody)
	}
	return resp
}

func BuildRespFromBetaOrder(resp *http.Response, reqBody *schema.OrderRequest,
	alphaMenu *schema.BetaMenu) *http.Response {

	body := schema.BetaRespBody{}

	var unifiedBody schema.OrderResponse // desired response body to return
	var err error
	if err = DecodeReqRespBody(resp.Body, &body); err != nil {
		log.Fatalf("failed decoding:, incompatible response body: %v", resp.Body)
	}

	populateUnifiedRespBodyForBeta(reqBody, &unifiedBody, alphaMenu)
	if resp.Body, err = EncodeReqRespBody(unifiedBody); err != nil {
		log.Fatalf("failed encoding: incompatible schema: %v", body)
	}
	return resp
}

func GetAlphaMenu(method string, destination *schema.AlphaMenuAddress, alphaMenu *schema.AlphaMenu) error {

	categoriesResp := RequestPOSClient(method, destination.AlphaCategoriesAddress, nil)
	if err := DecodeReqRespBody(categoriesResp.Body, alphaMenu.AlphaCategoriesMenu); err != nil {
		log.Printf("[categoriesResp] failed decoding:, incompatible response body: %v", err)
		return err
	}
	productsResp := RequestPOSClient(method, destination.AlphaProductsAddress, nil)
	if err := DecodeReqRespBody(productsResp.Body, alphaMenu.AlphaProductsMenu); err != nil {
		log.Printf("[ingredientsResp] failed decoding:, incompatible response body: %v", err)
		return err
	}
	ingredientsResp := RequestPOSClient(method, destination.AlphaIngredientsAddress, nil)
	if err := DecodeReqRespBody(ingredientsResp.Body, alphaMenu.AlphaIngredientsMenu); err != nil {
		log.Printf("[ingredientsResp] failed decoding:, incompatible response body: %v", err)
		return err
	}
	return nil
}
