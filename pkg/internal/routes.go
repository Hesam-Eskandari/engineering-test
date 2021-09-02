package internal

const (
	Menu = "/menu"

	AlphaClientBaseURL   = ":8081"
	BasePathAlpha        = "/alpha"
	MenuCategoriesAlpha  = Menu + "/categories"
	MenuProductsAlpha    = Menu + "/products"
	MenuIngredientsAlpha = Menu + "/ingredients"
	Orders               = "/orders"

	BetaClientBaseURL = ":8082"
	BasePathBeta      = "/beta"
	MenuBeta          = BasePathBeta + Menu
	OrdersCreateBeta  = Orders + "/create"
)
