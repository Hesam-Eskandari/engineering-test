package internal

const (
	AlphaClientBaseURL   = ":8081"
	BasePathAlpha        = "/alpha"
	Menu                 = "/menu"
	MenuCategoriesAlpha  = Menu + "/categories"
	MenuProductsAlpha    = Menu + "/products"
	MenuIngredientsAlpha = Menu + "/ingredients"
	Orders               = "/orders"

	BetaClientBaseURL = ":8082"
	BasePathBeta      = "/beta"
	MenuBeta          = BasePathBeta + "/menu"
	OrdersBeta        = BasePathBeta + "/orders"
	OrdersCreateBeta  = Orders + "/create"
)
