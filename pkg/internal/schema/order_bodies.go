package schema

type OrderItem struct {
	ID          string   `json:"id"`
	Quantity    int      `json:"quantity"`
	Size        string   `json:"size_id"`
	Ingredients []string `json:"ingredient_ids"`
	Extras      []string `json:"extra_ids"`
}

type OrderRequest struct {
	ID    string      `json:"id"`
	POS   string      `json:"pos"`
	Items []OrderItem `json:"items"`
}

type OrderResponseItem struct {
	Name     string   `json:"name"`
	Quantity int      `json:"quantity"`
	Size     string   `json:"size"`
	Extras   []string `json:"extras"`
	Price    float32  `json:"price"`
}

type OrderResponse struct {
	ID         string              `json:"id"`         // order id
	POS        string              `json:"pos"`        // pos
	Items      []OrderResponseItem `json:"items"`      // order items
	TotalPrice float32             `json:"totalPrice"` // total price
}

type AlphaReqBody struct {
	OrderId  string            `json:"orderId"`
	Products []AlphaReqProduct `json:"products"`
}

type AlphaReqProduct struct {
	ProductId     string   `json:"productId"`
	SizeId        string   `json:"sizeId"`
	IngredientIds []string `json:"ingredientIds"`
	Quantity      int      `json:"quantity"`
}

func NewAlphaReqProduct(productId, sizeId string, ingredientIds []string, quantity int) AlphaReqProduct {
	return AlphaReqProduct{
		ProductId:     productId,
		SizeId:        sizeId,
		IngredientIds: ingredientIds,
		Quantity:      quantity,
	}
}

type BetaReqBody struct {
	Id    string        `json:"Id"`
	Items []BetaReqItem `json:"Items"`
}

type BetaReqItem struct {
	CategoryId string   `json:"CategoryId"`
	ItemId     string   `json:"ItemId"`
	Quantity   int      `json:"Quantity"`
	AddOns     []string `json:"AddOns"`
}

func NewBetaReqItem(categoryId, itemId string, quantity int, addOns []string) BetaReqItem {
	return BetaReqItem{
		CategoryId: categoryId,
		ItemId:     itemId,
		Quantity:   quantity,
		AddOns:     addOns,
	}
}

type AlphaRespBody struct {
	OrderId  string             `json:"orderId"`
	Products []AlphaRespProduct `json:"products"`
}

type AlphaRespProduct struct {
	Name        string   `json:"name"`
	Size        string   `json:"size"`
	Ingredients []string `json:"ingredients"`
	Quantity    int      `json:"quantity"`
}

type BetaRespBody struct {
	OrderId string
	Items   []BetaRespItems
}

type BetaRespItems struct {
	Name     string
	AddOns   []string
	Quantity int
}
