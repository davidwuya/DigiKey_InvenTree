package main

import (
	"encoding/json"
	"log"
	"github.com/go-resty/resty/v2"
)

type InvenTreePartCategory struct {
	Pk     int    `json:"pk"`
	Name   string `json:"name"`
	Parent int    `json:"parent"`
}

func GetAllInvenTreeCategories() []InvenTreePartCategory {
	client := resty.New()
	resp, err := client.R().
		SetAuthScheme("Token").
		SetAuthToken(INVENTREE_TOKEN).
		SetHeader("Accept", "application/json").
		Get(InvenTreeServerAddress + "part/category/")
	if err != nil {
		log.Fatal(err)
	}

	var categories []InvenTreePartCategory
	err = json.Unmarshal(resp.Body(), &categories)

	if err != nil {
		log.Fatal(err)
	}

	return categories
}

func GetInvenTreeCategoryByName(name string, categories []InvenTreePartCategory) InvenTreePartCategory {
	for _, category := range categories {
		if category.Name == name {
			return category
		}
	}
	return InvenTreePartCategory{}

}

func CreateInvenTreeCategory(c *InvenTreePartCategory) InvenTreePartCategory {
	client := resty.New()
	resp, err := client.R().
		SetAuthScheme("Token").
		SetAuthToken(INVENTREE_TOKEN).
		// POST to /part/category endpoint
		SetBody(InvenTreePartCategory{
			Name:   c.Name,
			Parent: c.Parent,
		}).
		Post(InvenTreeServerAddress + "part/category/")
	if err != nil {
		log.Fatal(err)
	}

	var category InvenTreePartCategory
	err = json.Unmarshal(resp.Body(), &category)
	if err != nil {
		log.Fatal(err)
	}
	return category

}

func FindInvenTreeCategoryFromDigiKeyPart(d *DigiKeyPart, categories []InvenTreePartCategory) InvenTreePartCategory {
	parentPk := 1
	var category InvenTreePartCategory

	for _, categoryName := range d.LimitedTaxonomy {
		category = GetInvenTreeCategoryByName(categoryName, categories)
		if (category == InvenTreePartCategory{}) {
			newCategory := InvenTreePartCategory{
				Name:   categoryName,
				Parent: parentPk,
			}
			category = CreateInvenTreeCategory(&newCategory)
			categories = append(categories, category)
		}
		parentPk = category.Pk
	}

	return category
}
