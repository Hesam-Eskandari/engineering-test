package schema

import (
	"fmt"

	"github.com/flypay/engineering-test/pkg/internal"
)

// Menu represents unified menu structure to return to clients
type Menu struct {
	POS        string          `json:"pos"`
	Categories []*MenuCategory `json:"categories"`
}

// MenuCategory represents unified category structure for unified menu
type MenuCategory struct {
	Name          string             `json:"name"`
	ID            string             `json:"id"`
	Subcategories []*MenuSubcategory `json:"subcategories"`
}

// MenuSubcategory represents unified subcategory structure for unified menu
type MenuSubcategory struct {
	Name  string      `json:"name"`
	ID    string      `json:"id"`
	Items []*MenuItem `json:"items"`
}

// MenuItem represents unified item structure for unified menu
type MenuItem struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	ID          string                `json:"id"`
	Image       string                `json:"image"`
	Sizes       []*MenuItemSize       `json:"sizes"`
	Ingredients []*MenuItemIngredient `json:"ingredients"`
	Extras      []*MenuItemExtra      `json:"extras"`
}

// MenuItemExtra represents unified extra structure for unified menu
type MenuItemExtra struct {
	Name  string  `json:"name"`
	ID    string  `json:"id"`
	Price float32 `json:"price"`
}

// MenuItemSize represents unified size structure for unified menu
type MenuItemSize struct {
	Name  string  `json:"name"`
	ID    string  `json:"id"`
	Price float32 `json:"price"`
}

// MenuItemIngredient represents unified item ingredient structure for unified menu
type MenuItemIngredient struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

// BetaMenu represents menu structure of beta pos
type BetaMenu struct {
	Categories map[string]*BetaMenuCategory `json:"Categories"`
}

// BetaMenuCategory represents category structure for beta menu
type BetaMenuCategory struct {
	Name  string                   `json:"Name"`
	Items map[string]*BetaMenuItem `json:"Items"`
}

// BetaMenuItem represents item structure for beta menu
type BetaMenuItem struct {
	Name        string          `json:"Name"`
	Description string          `json:"Description"`
	Price       float32         `json:"Price"`
	AddOns      []BetaMenuAddOn `json:"AddOns"`
}

// BetaMenuAddOn represents addon structure for beta menu
type BetaMenuAddOn struct {
	Name  string  `json:"Name"`
	ID    string  `json:"Id"`
	Price float32 `json:"Price"`
}

// AlphaMenu represents menu structure containing all menus of alpha
type AlphaMenu struct {
	AlphaCategoriesMenu
	AlphaIngredientsMenu
	AlphaProductsMenu
}

// AlphaCategoriesMenu represents category menu structure of alpha pos
type AlphaCategoriesMenu struct {
	Categories []AlphaCategoriesMenuItem `json:"categories"`
}

// AlphaCategoriesMenuItem represents category item structure in category menu of alpha pos
type AlphaCategoriesMenuItem struct {
	CategoryID    string                       `json:"categoryId"`
	Name          string                       `json:"name"`
	Description   string                       `json:"description"`
	Subcategories []AlphaSubcategoriesMenuItem `json:"subcategories"`
}

// AlphaSubcategoriesMenuItem represents subcategory item structure un category menu of alpha pos
type AlphaSubcategoriesMenuItem struct {
	SubcategoryId string   `json:"subcategoryId"`
	Name          string   `json:"name"`
	Products      []string `json:"products"`
}

// AlphaProductsMenu represents products menu structure of alpha pos
type AlphaProductsMenu struct {
	Products []AlphaProductMenuItem `json:"products"`
}

// AlphaIngredientsMenu represents ingredients menu structure of alpha pos
type AlphaIngredientsMenu struct {
	Ingredients []AlphaIngredientsMenuItem `json:"ingredients"`
}

// AlphaIngredientsMenuItem represents ingredient item structure in ingredients menu of alpha pos
type AlphaIngredientsMenuItem struct {
	IngredientID     string `json:"ingredientId"`
	Name             string `json:"name"`
	GroupDescription string `json:"groupDescription"`
}

// AlphaProductMenuItem represents product item structure in products menu of alpha pos
type AlphaProductMenuItem struct {
	ProductID          string                      `json:"productId"`
	Name               string                      `json:"name"`
	Description        string                      `json:"description"`
	Image              string                      `json:"image"`
	Sizes              []AlphaProductMenuItemSize  `json:"sizes"`
	DefaultIngredients []string                    `json:"defaultIngredients"`
	Extras             []AlphaProductMenuItemExtra `json:"extras"`
}

// AlphaProductMenuItemSize represents product item size in products menu of alpha pos
type AlphaProductMenuItemSize struct {
	SizeID string  `json:"sizeId"`
	Name   string  `json:"name"`
	Price  float32 `json:"price"`
}

// AlphaProductMenuItemExtra represents product item extras in products menu of alpha pos
type AlphaProductMenuItemExtra struct {
	IngredientID string  `json:"ingredientId"`
	Price        float32 `json:"price"`
}

// AlphaMenuAddress represents URL structure containing all URLs of alpha to get menu
type AlphaMenuAddress struct {
	AlphaCategoriesAddress  string
	AlphaIngredientsAddress string
	AlphaProductsAddress    string
}

func NewAlphaMenuAddress() *AlphaMenuAddress {
	return &AlphaMenuAddress{
		AlphaCategoriesAddress:  fmt.Sprintf("http://%s"+"%s", internal.AlphaClientBaseURL, internal.MenuCategoriesAlpha),
		AlphaIngredientsAddress: fmt.Sprintf("http://%s"+"%s", internal.AlphaClientBaseURL, internal.MenuIngredientsAlpha),
		AlphaProductsAddress:    fmt.Sprintf("http://%s"+"%s", internal.AlphaClientBaseURL, internal.MenuProductsAlpha),
	}
}
