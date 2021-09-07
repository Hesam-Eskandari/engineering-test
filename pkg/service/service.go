package service

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/flypay/engineering-test/pkg/internal"
	"github.com/flypay/engineering-test/pkg/repository"
	"github.com/flypay/engineering-test/pkg/schema"
)

type serviceImpl struct {
	repo repository.Repository
}

func NewServiceImpl() Service {
	return &serviceImpl{
		repo: repository.NewRepositoryImpl(),
	}
}

// GetAlphaMenu calls alpha and receives its menus
func (s *serviceImpl) GetAlphaMenu(method string, destination *schema.AlphaMenuAddress, alphaMenu *schema.AlphaMenu) error {

	categoriesResp := s.RequestPOSClient(method, destination.AlphaCategoriesAddress, nil)
	if err := s.repo.DecodeReqRespBody(categoriesResp.Body, &alphaMenu.AlphaCategoriesMenu); err != nil {
		log.Printf("[categoriesResp] failed decoding:, incompatible response body: %v", err)
		return err
	}
	productsResp := s.RequestPOSClient(method, destination.AlphaProductsAddress, nil)
	if err := s.repo.DecodeReqRespBody(productsResp.Body, &alphaMenu.AlphaProductsMenu); err != nil {
		log.Printf("[ingredientsResp] failed decoding:, incompatible response body: %v", err)
		return err
	}
	ingredientsResp := s.RequestPOSClient(method, destination.AlphaIngredientsAddress, nil)
	if err := s.repo.DecodeReqRespBody(ingredientsResp.Body, &alphaMenu.AlphaIngredientsMenu); err != nil {
		log.Printf("[ingredientsResp] failed decoding:, incompatible response body: %v", err)
		return err
	}
	return nil
}

// GetBetaMenu calls beta and receives its menu
func (s *serviceImpl) GetBetaMenu(method, address string, betaMenu *schema.BetaMenu) error {
	betaMenuResp := s.RequestPOSClient(method, address, nil)
	if err := s.repo.DecodeReqRespBody(betaMenuResp.Body, betaMenu); err != nil {
		log.Printf("[categoriesResp] failed decoding:, incompatible response body: %v", err)
		return err
	}
	return nil
}

// GetAlphaReqBody creates compatible request to get menu from alpha
func (s *serviceImpl) GetAlphaReqBody(body *schema.OrderRequest) (io.ReadCloser, *schema.AlphaMenu, error) {
	// fetch beta menu to find corresponding category IDs for item IDs
	menu := new(schema.AlphaMenu)
	err := s.GetAlphaMenu(http.MethodGet, schema.NewAlphaMenuAddress(), menu)
	if err != nil {
		log.Printf("failed getting alpha menu, err: %s", err.Error())
		return nil, menu, err
	}

	defaultIngredientIDsMap := make(map[string][]string)
	allowedExtraIDsMap := make(map[string][]string)
	for _, product := range menu.Products {
		defaultIngredientIDsMap[product.ProductID] = product.DefaultIngredients
		extras := make([]string, 0, len(product.Extras))
		for _, extra := range product.Extras {
			extras = append(extras, extra.IngredientID)
		}
		allowedExtraIDsMap[product.ProductID] = extras
	}

	var alphaReqProducts []schema.AlphaReqProduct
	for _, item := range body.Items {
		ingredientIDs := make([]string, 0, len(item.Ingredients)+len(defaultIngredientIDsMap))
	menuLevel:
		for _, defaultIngredientID := range defaultIngredientIDsMap[item.ID] {
			for _, ingredientID := range item.Ingredients {
				if defaultIngredientID == ingredientID {
					continue menuLevel
				}
			}
			ingredientIDs = append(ingredientIDs, defaultIngredientID)
		}
		for _, allowedExtraID := range allowedExtraIDsMap[item.ID] {
			for _, extraID := range item.Extras {
				if allowedExtraID == extraID {
					ingredientIDs = append(ingredientIDs, extraID)
				}
			}
		}
		alphaReqProducts = append(alphaReqProducts,
			schema.NewAlphaReqProduct(item.ID, item.Size, ingredientIDs, item.Quantity))
	}

	alphaBody, err := s.repo.EncodeReqRespBody(&schema.AlphaReqBody{OrderId: body.ID, Products: alphaReqProducts})
	if err != nil {
		log.Printf("failed encoding alpha body, err: %s", err.Error())
		return nil, menu, err
	}
	return alphaBody, menu, nil
}

// GetBetaReqBody creates compatible request to get menu from beta
func (s *serviceImpl) GetBetaReqBody(body *schema.OrderRequest) (io.ReadCloser, *schema.BetaMenu, error) {
	// fetch beta menu
	menu := new(schema.BetaMenu)
	err := s.GetBetaMenu(http.MethodGet,
		fmt.Sprintf("http://%s"+"%s", internal.BetaClientBaseURL, internal.Menu), menu)
	if err != nil {
		log.Printf("failed getting beta menu, err: %s", err.Error())
		return nil, menu, err
	}
	var betaReqItems []schema.BetaReqItem
	for _, item := range body.Items {
		var allowedAddOns []string
		var categoryId string
		for catId, cat := range menu.Categories {
			for itemId, menuItem := range cat.Items {
				if itemId == item.ID {
					categoryId = catId
				}
				for _, addOn := range menuItem.AddOns {
					for _, reqAddOn := range item.Extras {
						if addOn.ID == reqAddOn {
							allowedAddOns = append(allowedAddOns, addOn.ID)
						}
					}

				}

			}
		}

		betaReqItems = append(betaReqItems, schema.NewBetaReqItem(categoryId, item.ID,
			item.Quantity, allowedAddOns))
	}
	betaBody, err := s.repo.EncodeReqRespBody(&schema.BetaReqBody{Id: body.ID, Items: betaReqItems})
	if err != nil {
		log.Printf("failed encoding beta body, err: %s", err.Error())
		return nil, menu, err
	}
	return betaBody, menu, nil
}

// RequestPOSClient calls a given POS and returns its response
func (s *serviceImpl) RequestPOSClient(method, destination string, body io.ReadCloser) (resp *http.Response) {
	client := &http.Client{}
	req, err := http.NewRequest(method, destination, body)
	if err != nil {
		return &http.Response{StatusCode: http.StatusInternalServerError}
	}
	req.Header.Add("Content-Type", "application/json")
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(err)
		return &http.Response{StatusCode: http.StatusServiceUnavailable}
	}
	return
}
