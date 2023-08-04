package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"sync"
)

func init() {

	log.SetOutput(os.Stdout)
	log.SetFormatter(&log.TextFormatter{
		ForceColors: true,
		FullTimestamp: false,
	})
	log.SetLevel(log.WarnLevel)
}

func ProcessParts() {
	// Initialize Digikey API and cache all inventree parts and locations
	var wg sync.WaitGroup
	wg.Add(2) // We have two concurrent tasks

	var AllInvenTreeLocations []InvenTreeLocation
	var AllInvenTreeCategories []InvenTreePartCategory
	var AllInvenTreeCompanies []InvenTreeCompany
	var AllInvenTreeStockItems []InvenTreeStockItem
	var AllInvenTreeParts []InvenTreePart
	go func() {
		defer wg.Done() // Done with caching
		AllInvenTreeLocations, AllInvenTreeCategories, AllInvenTreeCompanies, AllInvenTreeStockItems, AllInvenTreeParts = CacheAll()
	}()

	var d *DigikeyAPI
	var PN string
	go func() {
		defer wg.Done() // Done with initializing Digikey and processing input
		d = NewDigikeyAPI(APIKey, CLIENT_ID, OAuthState)
		PN = ProcessBarcode()
	}()

	wg.Wait() // Wait for both tasks to finish

	// Get product details and create DigiKeyPart object
	ThisDigiKeyPart := NewDigiKeyPart(GetProductDetails(d, PN))
	ThisDigiKeyPart.splitTaxonomy()

	var ThisPart InvenTreePart
	var ThisManufacturerPart ManufacturerPart
	var ThisStockItem InvenTreeStockItem
	var ThisSupplierPart SupplierPart

	ThisStockItem = GetInvenTreeStockItemByPartName(ThisDigiKeyPart.ProductDescription, AllInvenTreeStockItems, AllInvenTreeParts)
	ThisPart = GetInvenTreePartByName(ThisDigiKeyPart.ProductDescription, AllInvenTreeParts)
	ThisSupplierPart = GetSupplierPartByName(ThisDigiKeyPart)
	ThisManufacturerPart = GetManufacturerPartByName(ThisDigiKeyPart)

	if ThisPart.Pk == 0 {
		log.Info("Part does not exist.")
		ThisPart, AllInvenTreeCompanies = CreateInvenTreePart(ThisDigiKeyPart, AllInvenTreeCategories, AllInvenTreeCompanies)
		UploadThumbnailToInvenTreePart(&ThisPart, ThisDigiKeyPart)
		log.Info("Part created.")
	}

	if ThisManufacturerPart.Pk == 0 {
		log.Trace("Manufacturer part does not exist.")
		ThisManufacturerPart = CreateManufacturerPart(ThisPart, ThisDigiKeyPart, AllInvenTreeCompanies)
		log.Trace("Manufacturer part created.")
	}

	if ThisSupplierPart.Part == 0 {
		log.Trace("Supplier part does not exist.")
		ThisSupplierPart = CreateSupplierPart(&ThisPart, &ThisManufacturerPart, ThisDigiKeyPart)
		log.Trace("Supplier part created.")
	}

	if ThisStockItem.Pk == 0 {
		log.Trace("Stock item does not exist.")
		fmt.Print("Enter the location: ")
		LocString := ProcessInput()
		location := GetInvenTreeLocationByName(LocString, AllInvenTreeLocations)
		CreateStockItem(ThisSupplierPart, location, int(ProcessQuantity()))
		WriteLabels(ThisDigiKeyPart)
	}

	if ThisStockItem.Pk != 0 {
		fmt.Println("Stock item exists.")
		PrettyPrintStockItem(ThisStockItem, AllInvenTreeLocations)
		fmt.Println("Enter quantity to add to or subtract from stock. Enter 0 to rewrite labels.")
		StockUpdateQty := int(ProcessQuantity())
		if StockUpdateQty != 0 {
			UpdateStock(ThisStockItem, StockUpdateQty)
			fmt.Println("Stock updated.")
		} else {
			WriteLabels(ThisDigiKeyPart)
		}
	}

}

func main() {
	for {
		ProcessParts()
	}
}
