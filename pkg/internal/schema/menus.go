package schema

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

// AlphaIngredientsMenuItem represents ingredient item structure in ingredients menu of alpha pos
type AlphaIngredientsMenuItem struct {
	IngredientID     string `json:"ingredientId"`
	Name             string `json:"name"`
	GroupDescription string `json:"groupDescription"`
}
