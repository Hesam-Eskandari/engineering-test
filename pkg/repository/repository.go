package repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"sort"

	"github.com/flypay/engineering-test/pkg/internal"
	"github.com/flypay/engineering-test/pkg/schema"
)

type repositoryImpl struct{}

func NewRepositoryImpl() Repository {
	return &repositoryImpl{}
}

// DecodeReqRespBody decodes request or response arguments from request/response body to a schema
func (r *repositoryImpl) DecodeReqRespBody(body io.Reader, v interface{}) error {
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
func (r *repositoryImpl) EncodeReqRespBody(body interface{}) (io.ReadCloser, error) {
	if body == nil {
		return nil, errors.New("error encoding: request body is nil")
	}
	b, _ := json.Marshal(body)
	return ioutil.NopCloser(bytes.NewReader(b)), nil

}

// SortUnifiedMenu sorts any array in unified menu alphabetically by name
func (r *repositoryImpl) SortUnifiedMenu(menu *schema.Menu) {
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
func (r *repositoryImpl) SortUnifiedResponse(resp *schema.OrderResponse) {
	sort.Slice(resp.Items, func(i, j int) bool { return resp.Items[i].Name < resp.Items[j].Name })
	for _, item := range resp.Items {
		sort.Slice(item.Extras, func(i, j int) bool { return item.Extras[i] < item.Extras[j] })
		sort.Slice(item.Ingredients, func(i, j int) bool { return item.Ingredients[i] < item.Ingredients[j] })
	}
}

// PopulateUnifiedOrderRespBody populates the unified order response body with given order
func (r *repositoryImpl) PopulateUnifiedOrderRespBody(reqBody *schema.OrderRequest, menu *schema.Menu, unifiedBody *schema.OrderResponse) {
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

// PopulateUnifiedMenuFromAlphaMenu populates the unified menu struct with alpha menus
func (r *repositoryImpl) PopulateUnifiedMenuFromAlphaMenu(alphaMenu *schema.AlphaMenu, unifiedMenu *schema.Menu) {
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

// PopulateUnifiedMenuFromBetaMenu populates the unified menu struct with beta menu
func (r *repositoryImpl) PopulateUnifiedMenuFromBetaMenu(betaMenu *schema.BetaMenu, unifiedMenu *schema.Menu) {
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
