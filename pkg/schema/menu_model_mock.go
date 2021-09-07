package schema

import "github.com/flypay/engineering-test/pkg/internal"

func NewAlphaMenuMock() (menu *AlphaMenu) {
	menu = new(AlphaMenu)
	menu.AlphaCategoriesMenu = AlphaCategoriesMenu{
		Categories: []AlphaCategoriesMenuItem{
			{
				CategoryID:  "1001",
				Name:        "Burgers",
				Description: "",
				Subcategories: []AlphaSubcategoriesMenuItem{
					{
						SubcategoryId: "2001",
						Name:          "Veggie Burgers",
						Products: []string{
							"6001",
						},
					},
					{
						SubcategoryId: "2002",
						Name:          "Chicken Burgers",
						Products: []string{
							"6002",
						},
					},
				},
			},
			{
				CategoryID:  "1002",
				Name:        "Sides",
				Description: "",
				Subcategories: []AlphaSubcategoriesMenuItem{
					{
						SubcategoryId: "2003",
						Name:          "Potato",
						Products: []string{
							"6003",
						},
					},
				},
			},
			{
				CategoryID:  "1003",
				Name:        "Drinks",
				Description: "",
				Subcategories: []AlphaSubcategoriesMenuItem{
					{
						SubcategoryId: "2004",
						Name:          "Soft",
						Products: []string{
							"6004",
						},
					},
				},
			},
		},
	}
	menu.AlphaIngredientsMenu = AlphaIngredientsMenu{
		Ingredients: []AlphaIngredientsMenuItem{
			{
				IngredientID:     "9001",
				Name:             "Seeded Bun",
				GroupDescription: "Bread",
			},
			{
				IngredientID:     "9002",
				Name:             "Mushroom Patty",
				GroupDescription: "Patties",
			},
			{
				IngredientID:     "9003",
				Name:             "Lettuce",
				GroupDescription: "Salad",
			},
			{
				IngredientID:     "9004",
				Name:             "Tomato",
				GroupDescription: "Salad",
			},
			{
				IngredientID:     "9005",
				Name:             "Brioche Bun",
				GroupDescription: "Bread",
			},
			{
				IngredientID:     "9006",
				Name:             "Chicken Patty",
				GroupDescription: "Patties",
			},
			{
				IngredientID:     "9007",
				Name:             "Pickles",
				GroupDescription: "Salad",
			},
			{
				IngredientID:     "9008",
				Name:             "Slaw",
				GroupDescription: "Salad",
			},
			{
				IngredientID:     "9009",
				Name:             "Paprika Salt",
				GroupDescription: "Condiments",
			},
			{
				IngredientID:     "9010",
				Name:             "Extra Mayo",
				GroupDescription: "Condiments",
			},
		},
	}
	menu.AlphaProductsMenu = AlphaProductsMenu{
		Products: []AlphaProductMenuItem{
			{
				ProductID:   "6001",
				Name:        "Mushroom Burger",
				Description: "Whole fried mushroom with our unique spice blend",
				Image:       "https://images.unsplash.com/photo-1516774266634-15661f692c19?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=2532&q=80",
				Sizes: []AlphaProductMenuItemSize{
					{
						SizeID: "8001",
						Name:   "Individual",
						Price:  14.95,
					},
					{
						SizeID: "8002",
						Name:   "Double Up",
						Price:  21.95,
					},
				},
				DefaultIngredients: []string{
					"9001", "9002", "9003", "9004",
				},
				Extras: []AlphaProductMenuItemExtra{
					{
						IngredientID: "9010",
						Price:        0.99,
					},
					{
						IngredientID: "9007",
						Price:        0.99,
					},
				},
			},
			{
				ProductID:   "6002",
				Name:        "Buttermilk Chicken Burger",
				Description: "Juicy tender chicken thigh coated in crispy buttermilk batter.",
				Image:       "https://images.unsplash.com/photo-1597900121060-cf21f1cfa5e6?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=1234&q=80",
				Sizes: []AlphaProductMenuItemSize{
					{
						SizeID: "8001",
						Name:   "Individual",
						Price:  15.95,
					},
				},
				DefaultIngredients: []string{
					"9005", "9006", "9007", "9008",
				},
				Extras: []AlphaProductMenuItemExtra{},
			},
			{
				ProductID:   "6003",
				Name:        "Fries",
				Description: "Skinny cut locally sourced potato fries.",
				Image:       "https://images.unsplash.com/photo-1541592106381-b31e9677c0e5?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=2550&q=80",
				Sizes: []AlphaProductMenuItemSize{
					{
						SizeID: "8005",
						Name:   "Regular",
						Price:  3.95,
					},
					{
						SizeID: "8006",
						Name:   "Large",
						Price:  5.95,
					},
				},
				DefaultIngredients: []string{},
				Extras: []AlphaProductMenuItemExtra{
					{
						IngredientID: "9009",
						Price:        0.99,
					},
				},
			},
			{
				ProductID:   "6004",
				Name:        "Fizzy Cola",
				Description: "",
				Image:       "https://images.unsplash.com/photo-1592153995863-9fb8fe173740?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=2734&q=80",
				Sizes: []AlphaProductMenuItemSize{
					{
						SizeID: "8007",
						Name:   "Medium",
						Price:  2.95,
					},
					{
						SizeID: "8006",
						Name:   "Large",
						Price:  3.95,
					},
				},
				DefaultIngredients: []string{},
				Extras:             []AlphaProductMenuItemExtra{},
			},
		},
	}
	return
}

func NewBetaMenuMock() (menu *BetaMenu) {
	menu = new(BetaMenu)
	categories := make(map[string]*BetaMenuCategory)
	launchMenuItems := make(map[string]*BetaMenuItem)
	launchMenuItems["hjrlho"] = &BetaMenuItem{
		Name:        "Falafel Wrap",
		Description: "Made with real chickpeas",
		Price:       10.1,
		AddOns: []BetaMenuAddOn{
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
	}
	launchMenuItems["qmdehd"] = &BetaMenuItem{
		Name:        "Beef Shawarma",
		Description: "",
		Price:       12,
		AddOns: []BetaMenuAddOn{
			{
				Name:  "Extra Rice",
				ID:    "duaweu",
				Price: 2.5,
			},
		},
	}

	dessertItems := make(map[string]*BetaMenuItem)
	dessertItems["ndytqi"] = &BetaMenuItem{
		Name:        "Baklawa",
		Description: "",
		Price:       13,
		AddOns:      []BetaMenuAddOn{},
	}

	sideItems := make(map[string]*BetaMenuItem)
	sideItems["gtbiop"] = &BetaMenuItem{
		Name:  "Hummus",
		Price: 10.20,
	}
	sideItems["dwaecr"] = &BetaMenuItem{
		Name:  "Soup",
		Price: 4.99,
	}
	sideItems["ukhqnd"] = &BetaMenuItem{
		Name:  "Jam Donut",
		Price: 1.35,
	}
	sideItems["fllrsi"] = &BetaMenuItem{
		Name:  "Samosa",
		Price: 3,
	}

	drinkItems := make(map[string]*BetaMenuItem)
	drinkItems["ugqcbb"] = &BetaMenuItem{
		Name:  "Bottled Water",
		Price: 4,
	}
	drinkItems["pgzigb"] = &BetaMenuItem{
		Name:  "Juice",
		Price: 5,
	}
	drinkItems["sjalnl"] = &BetaMenuItem{
		Name:  "Bottled Beer",
		Price: 6,
	}

	categories["qqdluj"] = &BetaMenuCategory{
		Name:  "Lunch Menu",
		Items: launchMenuItems,
	}
	categories["dmcshb"] = &BetaMenuCategory{
		Name:  "Dessert",
		Items: dessertItems,
	}
	categories["gokpww"] = &BetaMenuCategory{
		Name:  "Sides",
		Items: sideItems,
	}
	categories["lzoaud"] = &BetaMenuCategory{
		Name:  "Drinks",
		Items: drinkItems,
	}
	menu.Categories = categories
	return
}

// NewUnifiedMenuPopulatedAlphaMenuMock returns a unified menu populated by mocked alpha menu
func NewUnifiedMenuPopulatedAlphaMenuMock() (alpha *Menu) {
	alpha = new(Menu)
	alpha.POS = internal.POSAlpha
	alpha.Categories = []*MenuCategory{
		{
			Name: "Burgers",
			ID:   "1001",
			Subcategories: []*MenuSubcategory{
				{
					Name: "Veggie Burgers",
					ID:   "2001",
					Items: []*MenuItem{
						{
							Name:        "Mushroom Burger",
							Description: "Whole fried mushroom with our unique spice blend",
							ID:          "6001",
							Image:       "https://images.unsplash.com/photo-1516774266634-15661f692c19?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=2532&q=80",
							Sizes: []*MenuItemSize{
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
							Ingredients: []*MenuItemIngredient{
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
							Extras: []*MenuItemExtra{
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
					Items: []*MenuItem{
						{
							Name:        "Buttermilk Chicken Burger",
							Description: "Juicy tender chicken thigh coated in crispy buttermilk batter.",
							ID:          "6002",
							Image:       "https://images.unsplash.com/photo-1597900121060-cf21f1cfa5e6?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=1234&q=80",
							Sizes: []*MenuItemSize{
								{
									Name:  "Individual",
									ID:    "8001",
									Price: 15.95,
								},
							},
							Ingredients: []*MenuItemIngredient{
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
							Extras: []*MenuItemExtra{},
						},
					},
				},
			},
		},
		{
			Name: "Sides",
			ID:   "1002",
			Subcategories: []*MenuSubcategory{
				{
					Name: "Potato",
					ID:   "2003",
					Items: []*MenuItem{
						{
							Name:        "Fries",
							Description: "Skinny cut locally sourced potato fries.",
							ID:          "6003",
							Image:       "https://images.unsplash.com/photo-1541592106381-b31e9677c0e5?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=2550&q=80",
							Sizes: []*MenuItemSize{
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
							Ingredients: []*MenuItemIngredient{},
							Extras: []*MenuItemExtra{
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
			Subcategories: []*MenuSubcategory{
				{
					Name: "Soft",
					ID:   "2004",
					Items: []*MenuItem{
						{
							Name:        "Fizzy Cola",
							Description: "",
							ID:          "6004",
							Image:       "https://images.unsplash.com/photo-1592153995863-9fb8fe173740?ixlib=rb-1.2.1&ixid=eyJhcHBfaWQiOjEyMDd9&auto=format&fit=crop&w=2734&q=80",
							Sizes: []*MenuItemSize{
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
							Ingredients: []*MenuItemIngredient{},
							Extras:      []*MenuItemExtra{},
						},
					},
				},
			},
		},
	}
	return
}

// NewUnifiedMenuPopulatedBetaMenuMock returns a unified menu populated by mocked beta menu
func NewUnifiedMenuPopulatedBetaMenuMock() (beta *Menu) {
	beta = new(Menu)
	beta.POS = internal.POSBeta
	beta.Categories = []*MenuCategory{
		{
			Name: "Lunch Menu",
			ID:   "qqdluj",
			Subcategories: []*MenuSubcategory{
				{
					Name: "Lunch Menu",
					ID:   "qqdluj",
					Items: []*MenuItem{
						{
							Name:        "Falafel Wrap",
							Description: "Made with real chickpeas",
							ID:          "hjrlho",
							Image:       "",
							Sizes: []*MenuItemSize{
								{
									Name:  "Regular",
									ID:    "hjrlho",
									Price: 10.1,
								},
							},
							Ingredients: []*MenuItemIngredient{},
							Extras: []*MenuItemExtra{
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
							Sizes: []*MenuItemSize{
								{
									Name:  "Regular",
									ID:    "qmdehd",
									Price: 12,
								},
							},
							Ingredients: []*MenuItemIngredient{},
							Extras: []*MenuItemExtra{
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
			Subcategories: []*MenuSubcategory{
				{
					Name: "Dessert",
					ID:   "dmcshb",
					Items: []*MenuItem{
						{
							Name:        "Baklawa",
							Description: "",
							ID:          "ndytqi",
							Image:       "",
							Sizes: []*MenuItemSize{
								{
									Name:  "Regular",
									ID:    "ndytqi",
									Price: 13,
								},
							},
							Ingredients: []*MenuItemIngredient{},
							Extras:      []*MenuItemExtra{},
						},
					},
				},
			},
		},
		{
			Name: "Sides",
			ID:   "gokpww",
			Subcategories: []*MenuSubcategory{
				{
					Name: "Sides",
					ID:   "gokpww",
					Items: []*MenuItem{
						{
							Name:        "Soup",
							Description: "",
							ID:          "dwaecr",
							Image:       "",
							Sizes: []*MenuItemSize{
								{
									Name:  "Regular",
									ID:    "dwaecr",
									Price: 4.99,
								},
							},
							Ingredients: []*MenuItemIngredient{},
							Extras:      []*MenuItemExtra{},
						},
						{
							Name:        "Samosa",
							Description: "",
							ID:          "fllrsi",
							Image:       "",
							Sizes: []*MenuItemSize{
								{
									Name:  "Regular",
									ID:    "fllrsi",
									Price: 3,
								},
							},
							Ingredients: []*MenuItemIngredient{},
							Extras:      []*MenuItemExtra{},
						},
						{
							Name:        "Hummus",
							Description: "",
							ID:          "gtbiop",
							Image:       "",
							Sizes: []*MenuItemSize{
								{
									Name:  "Regular",
									ID:    "gtbiop",
									Price: 10.2,
								},
							},
							Ingredients: []*MenuItemIngredient{},
							Extras:      []*MenuItemExtra{},
						},
						{
							Name:        "Jam Donut",
							Description: "",
							ID:          "ukhqnd",
							Image:       "",
							Sizes: []*MenuItemSize{
								{
									Name:  "Regular",
									ID:    "ukhqnd",
									Price: 1.35,
								},
							},
							Ingredients: []*MenuItemIngredient{},
							Extras:      []*MenuItemExtra{},
						},
					},
				},
			},
		},
		{
			Name: "Drinks",
			ID:   "lzoaud",
			Subcategories: []*MenuSubcategory{
				{
					Name: "Drinks",
					ID:   "lzoaud",
					Items: []*MenuItem{
						{
							Name:        "Bottled Water",
							Description: "",
							ID:          "ugqcbb",
							Image:       "",
							Sizes: []*MenuItemSize{
								{
									Name:  "Regular",
									ID:    "ugqcbb",
									Price: 4,
								},
							},
							Ingredients: []*MenuItemIngredient{},
							Extras:      []*MenuItemExtra{},
						},
						{
							Name:        "Juice",
							Description: "",
							ID:          "pgzigb",
							Image:       "",
							Sizes: []*MenuItemSize{
								{
									Name:  "Regular",
									ID:    "pgzigb",
									Price: 5,
								},
							},
							Ingredients: []*MenuItemIngredient{},
							Extras:      []*MenuItemExtra{},
						},
						{
							Name:        "Bottled Beer",
							Description: "",
							ID:          "sjalnl",
							Image:       "",
							Sizes: []*MenuItemSize{
								{
									Name:  "Regular",
									ID:    "sjalnl",
									Price: 6,
								},
							},
							Ingredients: []*MenuItemIngredient{},
							Extras:      []*MenuItemExtra{},
						},
					},
				},
			},
		},
	}
	return
}
