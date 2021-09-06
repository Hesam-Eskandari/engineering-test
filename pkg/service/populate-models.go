package service

import (
	"math"

	"github.com/flypay/engineering-test/pkg/internal"
	"github.com/flypay/engineering-test/pkg/internal/schema"
)

func PopulateUnifiedOrderRespBody(reqBody *schema.OrderRequest, menu *schema.Menu, unifiedBody *schema.OrderResponse) {
	unifiedBody.ID = reqBody.ID
	unifiedBody.POS = reqBody.POS

	itemMap := make(map[string]*schema.MenuItem)
	sizeMap := make(map[string]map[string]*schema.MenuItemSize)
	ingredientMap := make(map[string][]*schema.MenuItemIngredient)
	extraMap := make(map[string]map[string]*schema.MenuItemExtra)
	for _, menuCat := range menu.Categories {
		for _, subCat := range menuCat.Subcategories {
			for _, item := range subCat.Items {
				itemMap[item.ID] = item
				sizeMap[item.ID] = make(map[string]*schema.MenuItemSize)
				for _, size := range item.Sizes {
					sizeMap[item.ID][size.ID] = size
				}
				ingredientMap[item.ID] = item.Ingredients
				extraMap[item.ID] = make(map[string]*schema.MenuItemExtra)
				for _, extra := range item.Extras {
					extraMap[item.ID][extra.ID] = extra
				}
			}
		}
	}
	items := make([]*schema.OrderResponseItem, 0, len(reqBody.Items))
	var totalPrice float32
	for _, reqItem := range reqBody.Items {
		var price float32
		item := new(schema.OrderResponseItem)
		item.Quantity = reqItem.Quantity
		item.Name = itemMap[reqItem.ID].Name
		price += sizeMap[reqItem.ID][reqItem.Size].Price * float32(item.Quantity)
		item.Size = sizeMap[reqItem.ID][reqItem.Size].Name

		ingredients := make([]string, 0, len(ingredientMap[reqItem.ID]))
	menuLevel:
		for _, menuIngredient := range itemMap[reqItem.ID].Ingredients {
			for _, reqIngredientID := range reqItem.Ingredients {
				if reqIngredientID == menuIngredient.ID {
					continue menuLevel
				}
			}
			ingredients = append(ingredients, menuIngredient.Name)
		}
		item.Ingredients = ingredients

		extras := make([]string, 0, len(reqItem.Extras))
		for _, extraID := range reqItem.Extras {
			if _, ok := extraMap[reqItem.ID][extraID]; ok {
				extras = append(extras, extraMap[reqItem.ID][extraID].Name)
				price += extraMap[reqItem.ID][extraID].Price
			}
		}
		item.Extras = extras
		item.Price = float32(math.Round(float64(price*100)) / 100)
		items = append(items, item)
		totalPrice += price
	}
	unifiedBody.TotalPrice = float32(math.Round(float64(totalPrice*100)) / 100)
	unifiedBody.Items = items
}

func PopulateUnifiedMenuFromAlphaMenu(alphaMenu *schema.AlphaMenu, unifiedMenu *schema.Menu) {
	unifiedMenu.POS = internal.POSAlpha

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
		sizes := make([]*schema.MenuItemSize, 0)
		for _, productSize := range product.Sizes {
			size := &schema.MenuItemSize{
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
		extras := make([]*schema.MenuItemExtra, 0)
		for _, productExtra := range product.Extras {
			extra := &schema.MenuItemExtra{
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

func PopulateUnifiedMenuFromBetaMenu(betaMenu *schema.BetaMenu, unifiedMenu *schema.Menu) {
	unifiedMenu.POS = internal.POSBeta
	categories := make([]*schema.MenuCategory, 0, len(betaMenu.Categories))
	for catID, betaCategory := range betaMenu.Categories {
		category := new(schema.MenuCategory)
		category.ID = catID
		category.Name = betaCategory.Name

		subcategory := new(schema.MenuSubcategory)
		subcategory.ID = catID
		subcategory.Name = betaCategory.Name

		items := make([]*schema.MenuItem, 0, len(betaCategory.Items))
		for itemID, betaItem := range betaMenu.Categories[catID].Items {
			item := new(schema.MenuItem)
			item.ID = itemID
			item.Name = betaItem.Name
			item.Description = betaItem.Description
			sizes := make([]*schema.MenuItemSize, 0, 1)
			item.Sizes = append(sizes, &schema.MenuItemSize{Name: "Regular", ID: itemID, Price: betaItem.Price})

			extras := make([]*schema.MenuItemExtra, 0, len(betaItem.AddOns))
			for _, addOn := range betaItem.AddOns {
				extra := new(schema.MenuItemExtra)
				extra.ID = addOn.ID
				extra.Name = addOn.Name
				extra.Price = addOn.Price
				extras = append(extras, extra)
			}
			item.Extras = extras
			item.Ingredients = make([]*schema.MenuItemIngredient, 0)
			items = append(items, item)
		}
		subcategory.Items = items
		category.Subcategories = append(category.Subcategories, subcategory)
		categories = append(categories, category)
	}
	unifiedMenu.Categories = categories
}
