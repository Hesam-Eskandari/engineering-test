package service

import (
	"fmt"
	"log"
	"math"
	"net/http"

	"github.com/flypay/engineering-test/pkg/internal"
	"github.com/flypay/engineering-test/pkg/internal/schema"
)

func populateUnifiedRespBodyForBeta(reqBody *schema.OrderRequest, unifiedBody *schema.OrderResponse,
	menu *schema.BetaMenu) {
	unifiedBody.ID = reqBody.ID
	unifiedBody.POS = reqBody.POS

	type itemAddons struct {
		Name  string  `json:"name"`
		Price float32 `json:"price"`
	}

	newItemAddons := func(name string, price float32) *itemAddons {
		return &itemAddons{
			Name:  name,
			Price: price,
		}
	}

	type menuItem struct {
		Name  string  `json:"name"`
		Price float32 `json:"price"`
		//Quantity 	int 				`json:"quantity"`
		AllowedAddOns map[string]bool `json:"allowedAddOns"`
	}

	newMenuItem := func(name string, price float32, allowedAddOns map[string]bool) *menuItem {
		return &menuItem{
			Name:          name,
			Price:         price,
			AllowedAddOns: allowedAddOns,
		}
	}

	menuItemMap := make(map[string]*menuItem)
	addOnsMap := make(map[string]*itemAddons)
	for _, cat := range menu.Categories {
		for itemID, item := range cat.Items {
			allowedAddOns := make(map[string]bool)
			for _, addOn := range item.AddOns {
				allowedAddOns[addOn.ID] = true
				addOnsMap[addOn.ID] = newItemAddons(addOn.Name, addOn.Price)
			}
			menuItemMap[itemID] = newMenuItem(item.Name, item.Price, allowedAddOns)
		}
	}
	var orderResponseItems []schema.OrderResponseItem
	var totalPrice float32
	for _, item := range reqBody.Items {
		price := menuItemMap[item.ID].Price * float32(item.Quantity)
		extras := make([]string, 0)
		for _, extra := range append(item.Extras, item.Ingredients...) {
			if _, ok := menuItemMap[item.ID].AllowedAddOns[extra]; ok {
				price += addOnsMap[extra].Price
				extras = append(extras, addOnsMap[extra].Name)
			}
		}

		orderResponseItems = append(orderResponseItems, schema.OrderResponseItem{
			Name:     menuItemMap[item.ID].Name,
			Quantity: item.Quantity,
			Size:     "regular",
			Extras:   extras,
			Price:    float32(math.Round(float64(price*1000)) / 1000),
		})
		totalPrice += price
	}
	unifiedBody.Items = orderResponseItems
	unifiedBody.TotalPrice = float32(math.Round(float64(totalPrice*1000)) / 1000)
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
		for _, extraID := range append(item.Ingredients, item.Extras...) {
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

func PopulateUnifiedMenuFromAlphaMenu(alphaMenu *schema.AlphaMenu, unifiedMenu *schema.Menu) {
	ingredientMap := make(map[string]*schema.MenuItemIngredient)
	for _, ingredient := range alphaMenu.AlphaIngredientsMenu.Ingredients {
		ingredientMap[ingredient.IngredientID] = &schema.MenuItemIngredient{
			Name: ingredient.Name,
			ID:   ingredient.IngredientID,
		}
	}
	extraMap := make(map[string]float32)
	for _, product := range alphaMenu.AlphaProductsMenu.Products {
		for _, extra := range product.Extras {
			extraMap[extra.IngredientID] = extra.Price
		}
	}
	itemMap := make(map[string]*schema.MenuItem)
	for _, product := range alphaMenu.AlphaProductsMenu.Products {
		sizes := make([]*schema.MenuItemDetail, 0)
		for _, productSize := range product.Sizes {
			size := &schema.MenuItemDetail{
				Name:  productSize.Name,
				ID:    productSize.SizeID,
				Price: productSize.Price,
			}
			sizes = append(sizes, size)
		}
		ingredients := make([]*schema.MenuItemIngredient, 0)
		for _, productIngredientID := range product.DefaultIngredients {
			ingredient := &schema.MenuItemIngredient{
				Name: ingredientMap[productIngredientID].Name,
				ID:   productIngredientID,
			}
			ingredients = append(ingredients, ingredient)
		}
		extras := make([]*schema.MenuItemDetail, 0)
		for _, productExtra := range product.Extras {
			extra := &schema.MenuItemDetail{
				Name:  ingredientMap[productExtra.IngredientID].Name,
				ID:    productExtra.IngredientID,
				Price: productExtra.Price,
			}
			extras = append(extras, extra)
		}
		itemMap[product.ProductID] = &schema.MenuItem{
			Name:        product.Name,
			Description: product.Description,
			ID:          product.ProductID,
			Image:       product.Image,
			Sizes:       sizes,
			Ingredients: ingredients,
			Extras:      extras,
		}
	}

	unifiedMenu.POS = internal.POSAlpha
	for _, cat := range alphaMenu.Categories {
		subcategories := make([]*schema.MenuSubcategory, 0)
		for _, subCat := range cat.Subcategories {
			items := make([]*schema.MenuItem, 0)
			for _, productID := range subCat.Products {
				items = append(items, itemMap[productID])
			}
			menuSubcategory := &schema.MenuSubcategory{
				ID:    subCat.SubcategoryId,
				Name:  subCat.Name,
				Items: items,
			}
			subcategories = append(subcategories, menuSubcategory)
		}
		unifiedMenu.Categories = append(unifiedMenu.Categories, &schema.MenuCategory{
			ID:            cat.CategoryID,
			Name:          cat.Name,
			Subcategories: subcategories,
		})
	}
}
