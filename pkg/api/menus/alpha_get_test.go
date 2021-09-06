package menus

import (
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/flypay/engineering-test/pkg/internal"
	"github.com/flypay/engineering-test/pkg/internal/schema"
	"github.com/flypay/engineering-test/pkg/service"
)

// TestAlphaMenu calls unified service url to get alpha menu
func TestAlphaMenu(t *testing.T) {
	alphaMenuExpected := setUpAlphaMenu()
	t.Run("get alpha menu", func(t *testing.T) {
		alphaMenu := new(schema.Menu)
		alphaMenuResp := service.RequestPOSClient(http.MethodGet,
			fmt.Sprintf("http://:8086"+internal.BasePathAlpha+internal.Menu), nil)
		if err := service.DecodeReqRespBody(alphaMenuResp.Body, alphaMenu); err != nil {
			t.Fatalf("failed decoding body resp. err: %s", err.Error())
		}
		assertEqualMenus(t, alphaMenuExpected, alphaMenu)
	})
}

// assertEqualMenus runs deep equal on two menus and compares them
func assertEqualMenus(t *testing.T, expectedModel, retrievedModel *schema.Menu) {
	service.SortUnifiedMenu(expectedModel)
	service.SortUnifiedMenu(retrievedModel)
	if !reflect.DeepEqual(expectedModel, retrievedModel) {
		t.Fatal("The two struct are not equal")
	}
}

/*
// assertEqualMenusAlternative is  an Alternative function to check equality of two structs
func assertEqualMenusAlternative(t *testing.T, expectedMenu, retrievedMenu *schema.Menu) {
	service.SortUnifiedMenu(expectedMenu)
	service.SortUnifiedMenu(retrievedMenu)
	jsonExpectedMenu, _ := json.Marshal(expectedMenu)
	jsonRetrievedMenu, _ := json.Marshal(retrievedMenu)
	if string(jsonExpectedMenu) != string(jsonRetrievedMenu) {
		t.Fatalf("the two menus are not equal")
	}
}
*/

// setUpAlphaMenu returns expected alpha menu that should be returned
func setUpAlphaMenu() (alpha *schema.Menu) {
	alpha = new(schema.Menu)
	alpha.POS = internal.POSAlpha
	alpha.Categories = []*schema.MenuCategory{
		{
			Name: "Burgers",
			ID:   "1001",
			Subcategories: []*schema.MenuSubcategory{
				{
					Name: "Veggie Burgers",
					ID:   "2001",
					Items: []*schema.MenuItem{
						{
							Name:        "Mushroom Burger",
							Description: "Whole fried mushroom with our unique spice blend",
							ID:          "6001",
							Image:       "https://images.unsplash.com/photo-1516774266634-15661f692c19?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=2532&q=80",
							Sizes: []*schema.MenuItemSize{
								{
									Name:  "Individual",
									ID:    "8001",
									Price: 14.95,
								},
								{
									Name:  "Double Up",
									ID:    "8002",
									Price: 21.95,
								},
							},
							Ingredients: []*schema.MenuItemIngredient{
								{
									Name: "Seeded Bun",
									ID:   "9001",
								},
								{
									Name: "Mushroom Patty",
									ID:   "9002",
								},
								{
									Name: "Lettuce",
									ID:   "9003",
								},
								{
									Name: "Tomato",
									ID:   "9004",
								},
							},
							Extras: []*schema.MenuItemExtra{
								{
									Name:  "Extra Mayo",
									ID:    "9010",
									Price: 0.99,
								},
								{
									Name:  "Pickles",
									ID:    "9007",
									Price: 0.99,
								},
							},
						},
					},
				},
				{
					Name: "Chicken Burgers",
					ID:   "2002",
					Items: []*schema.MenuItem{
						{
							Name:        "Buttermilk Chicken Burger",
							Description: "Juicy tender chicken thigh coated in crispy buttermilk batter.",
							ID:          "6002",
							Image:       "https://images.unsplash.com/photo-1597900121060-cf21f1cfa5e6?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=1234&q=80",
							Sizes: []*schema.MenuItemSize{
								{
									Name:  "Individual",
									ID:    "8001",
									Price: 15.95,
								},
							},
							Ingredients: []*schema.MenuItemIngredient{
								{
									Name: "Brioche Bun",
									ID:   "9005",
								},
								{
									Name: "Chicken Patty",
									ID:   "9006",
								},
								{
									Name: "Pickles",
									ID:   "9007",
								},
								{
									Name: "Slaw",
									ID:   "9008",
								},
							},
							Extras: []*schema.MenuItemExtra{},
						},
					},
				},
			},
		},
		{
			Name: "Sides",
			ID:   "1002",
			Subcategories: []*schema.MenuSubcategory{
				{
					Name: "Potato",
					ID:   "2003",
					Items: []*schema.MenuItem{
						{
							Name:        "Fries",
							Description: "Skinny cut locally sourced potato fries.",
							ID:          "6003",
							Image:       "https://images.unsplash.com/photo-1541592106381-b31e9677c0e5?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=2550&q=80",
							Sizes: []*schema.MenuItemSize{
								{
									Name:  "Regular",
									ID:    "8005",
									Price: 3.95,
								},
								{
									Name:  "Large",
									ID:    "8006",
									Price: 5.95,
								},
							},
							Ingredients: []*schema.MenuItemIngredient{},
							Extras: []*schema.MenuItemExtra{
								{
									Name:  "Paprika Salt",
									ID:    "9009",
									Price: 0.99,
								},
							},
						},
					},
				},
			},
		},
		{
			Name: "Drinks",
			ID:   "1003",
			Subcategories: []*schema.MenuSubcategory{
				{
					Name: "Soft",
					ID:   "2004",
					Items: []*schema.MenuItem{
						{
							Name:        "Fizzy Cola",
							Description: "",
							ID:          "6004",
							Image:       "https://images.unsplash.com/photo-1592153995863-9fb8fe173740?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=2734&q=80",
							Sizes: []*schema.MenuItemSize{
								{
									Name:  "Medium",
									ID:    "8007",
									Price: 2.95,
								},
								{
									Name:  "Large",
									ID:    "8006",
									Price: 3.95,
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
