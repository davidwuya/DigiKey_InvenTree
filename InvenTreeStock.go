package main

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"fmt"
)

type InvenTreeStockItem struct {
	Pk       int     `json:"pk"`       // primary key of stock item
	Part     int     `json:"part"`     // primary key of associated part
	Location int     `json:"location"` // primary key of associated location
	Quantity float32 `json:"quantity"` // quantity of stock item
	Status   int     `json:"status"`   // 10 means OK
}

func GetAllInvenTreeStockItems() []InvenTreeStockItem {
	client := resty.New()
	client.SetRetryCount(3)
	resp, err := client.R().
		SetAuthScheme("Token").
		SetAuthToken(INVENTREE_TOKEN).
		SetHeader("Accept", "application/json").
		Get(InvenTreeServerAddress + "stock/")
	if err != nil {
		log.Fatal(err)
	}
	var stockItems []InvenTreeStockItem
	err = json.Unmarshal(resp.Body(), &stockItems)
	if err != nil {
		log.Fatal(err)
	}
	return stockItems
}

func GetInvenTreeStockItemByPartName(partName string, AllInvenTreeStockItems []InvenTreeStockItem, AllInvenTreeParts []InvenTreePart) InvenTreeStockItem {
	for _, stockItem := range AllInvenTreeStockItems {
		if stockItem.Part == GetInvenTreePartByName(partName, AllInvenTreeParts).Pk {
			return stockItem
		}
	}
	return InvenTreeStockItem{}
}

func CreateStockItem(part SupplierPart, location InvenTreeLocation, quantity int) InvenTreeStockItem {
	client := resty.New()
	client.SetRetryCount(3)
	resp, err := client.R().
		SetAuthScheme("Token").
		SetAuthToken(INVENTREE_TOKEN).
		SetBody(InvenTreeStockItem{
			Part:     part.Part,
			Location: location.Pk,
			Quantity: float32(quantity),
			Status:   10,
		}).
		SetHeader("Accept", "application/json").
		Post(InvenTreeServerAddress + "stock/")
	if err != nil {
		log.Fatal(err)
	}
	var stockItem InvenTreeStockItem
	err = json.Unmarshal(resp.Body(), &stockItem)
	if err != nil {
		log.Fatal(err)
	}
	return stockItem
}

func UpdateStock(stockItem InvenTreeStockItem, quantity int) bool {
	client := resty.New()
	client.SetRetryCount(3)
	Endpoint := ""
	if quantity < 0 {
		Endpoint = "stock/remove/"
		quantity = -quantity
	} else {
		Endpoint = "stock/add/"
	}

	resp, err := client.R().
		SetAuthScheme("Token").
		SetAuthToken(INVENTREE_TOKEN).
		SetBody(map[string]interface{}{
			"items": []map[string]interface{}{
				{
					"pk":       stockItem.Pk,
					"quantity": quantity,
				},
			},
		}).
		SetHeader("Accept", "application/json").
		Post(InvenTreeServerAddress + Endpoint)
	if err != nil {
		log.Fatal(err)
	}
	if resp.IsSuccess() {
		log.Info("Stock updated successfully")
	} else {
		log.Error("Stock update failed")
	}
	return resp.IsSuccess()
}

func PrettyPrintStockItem(stockItem InvenTreeStockItem, AllInvenTreeLocations []InvenTreeLocation) {
	fmt.Println("Location: ", GetInvenTreeLocationNameByPK(stockItem.Location, AllInvenTreeLocations))
	fmt.Println("Current Stock Quantity: ", stockItem.Quantity)
}