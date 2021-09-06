package menus

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/flypay/engineering-test/pkg/internal"
	"github.com/flypay/engineering-test/pkg/internal/schema"
	"github.com/flypay/engineering-test/pkg/service"
)

// TestBetaMenu calls unified service url to get beta menu
func TestBetaMenu(t *testing.T) {
	betaMenuExpected := setUpBetaMenu()
	t.Run("get beta menu", func(t *testing.T) {
		betaMenu := new(schema.Menu)
		betaMenuResp := service.RequestPOSClient(http.MethodGet,
			fmt.Sprintf("http://:8086"+internal.BasePathBeta+internal.Menu), nil)
		if err := service.DecodeReqRespBody(betaMenuResp.Body, betaMenu); err != nil {
			t.Fatalf("failed decoding body resp. err: %s", err.Error())
		}
		assertEqualMenus(t, betaMenuExpected, betaMenu)
	})
}

// setUpBetaMenu returns expected beta menu that should be returned
func setUpBetaMenu() (beta *schema.Menu) {
	beta = new(schema.Menu)
	beta.POS = internal.POSBeta
	beta.Categories = []*schema.MenuCategory{
		{
			Name: "Lunch Menu",
			ID:   "qqdluj",
			Subcategories: []*schema.MenuSubcategory{
				{
					Name: "Lunch Menu",
					ID:   "qqdluj",
					Items: []*schema.MenuItem{
						{
							Name:        "Falafel Wrap",
							Description: "Made with real chickpeas",
							ID:          "hjrlho",
							Image:       "",
							Sizes: []*schema.MenuItemSize{
								{
									Name:  "Regular",
									ID:    "hjrlho",
									Price: 10.1,
								},
							},
							Ingredients: []*schema.MenuItemIngredient{},
							Extras: []*schema.MenuItemExtra{
								{
									Name:  "Extra Falafel Ball",
									ID:    "oecwhs",
									Price: 1,
								},
								{
									Name:  "Extra Hummus",
									ID:    "warnlj",
									Price: 0.5,
								},
								{
									Name:  "Extra Tabouli",
									ID:    "yzhllj",
									Price: 0.25,
								},
								{
									Name:  "Extra Chilli Sauce",
									ID:    "cbnufj",
									Price: 0.99,
								},
							},
						},
						{
							Name:        "Beef Shawarma",
							Description: "",
							ID:          "qmdehd",
							Image:       "",
							Sizes: []*schema.MenuItemSize{
								{
									Name:  "Regular",
									ID:    "qmdehd",
									Price: 12,
								},
							},
							Ingredients: []*schema.MenuItemIngredient{},
							Extras: []*schema.MenuItemExtra{
								{
									Name:  "Extra Rice",
									ID:    "duaweu",
									Price: 2.5,
								},
							},
						},
					},
				},
			},
		},
		{
			Name: "Dessert",
			ID:   "dmcshb",
			Subcategories: []*schema.MenuSubcategory{
				{
					Name: "Dessert",
					ID:   "dmcshb",
					Items: []*schema.MenuItem{
						{
							Name:        "Baklawa",
							Description: "",
							ID:          "ndytqi",
							Image:       "",
							Sizes: []*schema.MenuItemSize{
								{
									Name:  "Regular",
									ID:    "ndytqi",
									Price: 13,
								},
							},
							Ingredients: []*schema.MenuItemIngredient{},
							Extras:      []*schema.MenuItemExtra{},
						},
					},
				},
			},
		},
		{
			Name: "Sides",
			ID:   "gokpww",
			Subcategories: []*schema.MenuSubcategory{
				{
					Name: "Sides",
					ID:   "gokpww",
					Items: []*schema.MenuItem{
						{
							Name:        "Soup",
							Description: "",
							ID:          "dwaecr",
							Image:       "",
							Sizes: []*schema.MenuItemSize{
								{
									Name:  "Regular",
									ID:    "dwaecr",
									Price: 4.99,
								},
							},
							Ingredients: []*schema.MenuItemIngredient{},
							Extras:      []*schema.MenuItemExtra{},
						},
						{
							Name:        "Samosa",
							Description: "",
							ID:          "fllrsi",
							Image:       "",
							Sizes: []*schema.MenuItemSize{
								{
									Name:  "Regular",
									ID:    "fllrsi",
									Price: 3,
								},
							},
							Ingredients: []*schema.MenuItemIngredient{},
							Extras:      []*schema.MenuItemExtra{},
						},
						{
							Name:        "Hummus",
							Description: "",
							ID:          "gtbiop",
							Image:       "",
							Sizes: []*schema.MenuItemSize{
								{
									Name:  "Regular",
									ID:    "gtbiop",
									Price: 10.2,
								},
							},
							Ingredients: []*schema.MenuItemIngredient{},
							Extras:      []*schema.MenuItemExtra{},
						},
						{
							Name:        "Jam Donut",
							Description: "",
							ID:          "ukhqnd",
							Image:       "",
							Sizes: []*schema.MenuItemSize{
								{
									Name:  "Regular",
									ID:    "ukhqnd",
									Price: 1.35,
								},
							},
							Ingredients: []*schema.MenuItemIngredient{},
							Extras:      []*schema.MenuItemExtra{},
						},
					},
				},
			},
		},
		{
			Name: "Drinks",
			ID:   "lzoaud",
			Subcategories: []*schema.MenuSubcategory{
				{
					Name: "Drinks",
					ID:   "lzoaud",
					Items: []*schema.MenuItem{
						{
							Name:        "Bottled Water",
							Description: "",
							ID:          "ugqcbb",
							Image:       "",
							Sizes: []*schema.MenuItemSize{
								{
									Name:  "Regular",
									ID:    "ugqcbb",
									Price: 4,
								},
							},
							Ingredients: []*schema.MenuItemIngredient{},
							Extras:      []*schema.MenuItemExtra{},
						},
						{
							Name:        "Juice",
							Description: "",
							ID:          "pgzigb",
							Image:       "",
							Sizes: []*schema.MenuItemSize{
								{
									Name:  "Regular",
									ID:    "pgzigb",
									Price: 5,
								},
							},
							Ingredients: []*schema.MenuItemIngredient{},
							Extras:      []*schema.MenuItemExtra{},
						},
						{
							Name:        "Bottled Beer",
							Description: "",
							ID:          "sjalnl",
							Image:       "",
							Sizes: []*schema.MenuItemSize{
								{
									Name:  "Regular",
									ID:    "sjalnl",
									Price: 6,
								},
							},
							Ingredients: []*schema.MenuItemIngredient{},
							Extras:      []*schema.MenuItemExtra{},
						},
					},
				},
			},
		},
	}
	return
}
