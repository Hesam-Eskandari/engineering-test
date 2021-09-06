package service

import (
	"github.com/flypay/engineering-test/pkg/internal/schema"
	"log"
)

func BuildRespFromAlphaOrder(reqBody *schema.OrderRequest, unifiedBody *schema.OrderResponse, alphaMenu *schema.AlphaMenu) error {
	unifiedMenu := new(schema.Menu)

	PopulateUnifiedMenuFromAlphaMenu(alphaMenu, unifiedMenu)
	PopulateUnifiedOrderRespBody(reqBody, unifiedMenu, unifiedBody)
	return nil
}

func BuildRespFromBetaOrder(reqBody *schema.OrderRequest, betaMenu *schema.BetaMenu, unifiedBody *schema.OrderResponse) error {
	unifiedMenu := new(schema.Menu)
	PopulateUnifiedMenuFromBetaMenu(betaMenu, unifiedMenu)
	PopulateUnifiedOrderRespBody(reqBody, unifiedMenu, unifiedBody)
	return nil
}

func GetAlphaMenu(method string, destination *schema.AlphaMenuAddress, alphaMenu *schema.AlphaMenu) error {

	categoriesResp := RequestPOSClient(method, destination.AlphaCategoriesAddress, nil)
	if err := DecodeReqRespBody(categoriesResp.Body, &alphaMenu.AlphaCategoriesMenu); err != nil {
		log.Printf("[categoriesResp] failed decoding:, incompatible response body: %v", err)
		return err
	}
	productsResp := RequestPOSClient(method, destination.AlphaProductsAddress, nil)
	if err := DecodeReqRespBody(productsResp.Body, &alphaMenu.AlphaProductsMenu); err != nil {
		log.Printf("[ingredientsResp] failed decoding:, incompatible response body: %v", err)
		return err
	}
	ingredientsResp := RequestPOSClient(method, destination.AlphaIngredientsAddress, nil)
	if err := DecodeReqRespBody(ingredientsResp.Body, &alphaMenu.AlphaIngredientsMenu); err != nil {
		log.Printf("[ingredientsResp] failed decoding:, incompatible response body: %v", err)
		return err
	}
	return nil
}

func GetBetaMenu(method, address string, betaMenu *schema.BetaMenu) error {
	betaMenuResp := RequestPOSClient(method, address, nil)
	if err := DecodeReqRespBody(betaMenuResp.Body, betaMenu); err != nil {
		log.Printf("[categoriesResp] failed decoding:, incompatible response body: %v", err)
		return err
	}
	return nil
}
