package schema

// Menu represents unified menu structure to return to clients
type Menu struct {
	Categories []*MenuCategory `json:"categories"`
	POS        string          `jsom:"pos"`
}

type MenuCategory struct {
	Name          string             `json:"name"`
	ID            string             `json:"id"`
	Subcategories []*MenuSubcategory `json:"subcategories"`
}

type MenuSubcategory struct {
	Name  string      `json:"name"`
	ID    string      `json:"id"`
	Items []*MenuItem `json:"items"`
}

type MenuItem struct {
	Name        string                `json:"name"`
	Description string                `json:"description"`
	ID          string                `json:"id"`
	Image       string                `json:"image"`
	Sizes       []*MenuItemDetail     `json:"sizes"`
	Ingredients []*MenuItemIngredient `json:"ingredients"`
	Extras      []*MenuItemDetail     `json:"extras"`
}

type MenuItemDetail struct {
	Name  string  `json:"name"`
	ID    string  `json:"id"`
	Price float32 `json:"price"`
}

type MenuItemIngredient struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

type NewBetaMenu struct {
	CategoriesOrder []string                     `json:"CategoriesOrder"`
	Categories      map[string]*BetaMenuCategory `json:"Categories"`
}

type BetaMenuCategory struct {
	Name  string                   `json:"Name"`
	Items map[string]*BetaMenuItem `json:"Items"`
}

type BetaMenuItem struct {
	Name        string          `json:"Name"`
	Description string          `json:"Description"`
	Price       float32         `json:"Price"`
	AddOns      []BetaMenuAddOn `json:"AddOns"`
}

type BetaMenuAddOn struct {
	Name  string  `json:"Name"`
	ID    string  `json:"Id"`
	Price float32 `json:"Price"`
}

// BetaMenu represents menu structure of beta pos
type BetaMenu struct {
	Categories map[string]struct {
		Name  string `json:"Name"`
		Items map[string]struct {
			Name     string  `json:"Name"`
			Price    float32 `json:"Price"`
			Quantity int     `json:"Quantity"`
			AddOns   []struct {
				ID    string  `json:"Id"`
				Name  string  `json:"Name"`
				Price float32 `json:"Price"`
			} `json:"AddOns"`
		} `json:"Items"`
	} `json:"Categories"`
}

// AlphaProductsMenu represents products menu structure of alpha pos
type AlphaProductsMenu struct {
	Products []AlphaProductMenuItem `json:"products"`
}

func NewEmptyAlphaProductsMenu() *AlphaProductsMenu {
	return &AlphaProductsMenu{
		Products: make([]AlphaProductMenuItem, 0),
	}
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

// AlphaIngredientsMenu represents ingredients menu structure of alpha pos
type AlphaIngredientsMenu struct {
	Ingredients []AlphaIngredientsMenuItem `json:"ingredients"`
}

func NewEmptyAlphaIngredientsMenu() *AlphaIngredientsMenu {
	return &AlphaIngredientsMenu{
		Ingredients: make([]AlphaIngredientsMenuItem, 0),
	}
}

// AlphaIngredientsMenuItem represents ingredient item structure in ingredients menu of alpha pos
type AlphaIngredientsMenuItem struct {
	IngredientID     string `json:"ingredientId"`
	Name             string `json:"name"`
	GroupDescription string `json:"groupDescription"`
}

func NewEmptyAlphaMenu() *AlphaMenu {
	return &AlphaMenu{
		AlphaCategoriesMenu:  NewEmptyAlphaCategoriesMenu(),
		AlphaIngredientsMenu: NewEmptyAlphaIngredientsMenu(),
		AlphaProductsMenu:    NewEmptyAlphaProductsMenu(),
	}
}

type AlphaMenu struct {
	*AlphaCategoriesMenu
	*AlphaIngredientsMenu
	*AlphaProductsMenu
}

type AlphaMenuAddress struct {
	AlphaCategoriesAddress  string
	AlphaIngredientsAddress string
	AlphaProductsAddress    string
}

func NewEmptyAlphaCategoriesMenu() *AlphaCategoriesMenu {
	return &AlphaCategoriesMenu{
		Categories: make([]AlphaCategoriesMenuItem, 0),
	}
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
