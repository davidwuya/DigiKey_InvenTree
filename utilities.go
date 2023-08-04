package main

import (
	"bufio"
	"fmt"
	log "github.com/sirupsen/logrus"
	"os"
	"regexp"
	"strconv"
)

func decodeBarcode(barcode string) string {
	dkPartNumberRegex := regexp.MustCompile(`\$P(.*?)\$1P`)
	mouserPartNumberRegex := regexp.MustCompile(`\$1P(.*?)\$Q`)

	dkPartNumberMatch := dkPartNumberRegex.FindStringSubmatch(barcode)
	mouserPartNumberMatch := mouserPartNumberRegex.FindStringSubmatch(barcode)

	var partNumber string
	if dkPartNumberMatch != nil {
		log.Debug("Digi-Key Barcode Found")
		partNumber = dkPartNumberMatch[1]
	} else if mouserPartNumberMatch != nil {
		log.Debug("Mouser Barcode Found")
		partNumber = mouserPartNumberMatch[1]
	} else {
		log.Debug("No Barcode Found, assuming part number")
		partNumber = barcode
	}

	return partNumber
}

func CacheAll() ([]InvenTreeLocation, []InvenTreePartCategory, []InvenTreeCompany, []InvenTreeStockItem, []InvenTreePart) {
	// Cache all the things
	AllInvenTreeLocations := GetAllInvenTreeLocations()
	AllInvenTreeCategories := GetAllInvenTreeCategories()
	AllInvenTreeCompanies := GetAllInvenTreeCompanies()
	AllInvenTreeStockItems := GetAllInvenTreeStockItems()
	AllInvenTreeParts := GetAllInvenTreeParts()
	return AllInvenTreeLocations, AllInvenTreeCategories, AllInvenTreeCompanies, AllInvenTreeStockItems, AllInvenTreeParts
}

func ProcessBarcode() string {
	fmt.Print("Enter barcode: ")
	barcode := ProcessInput()
	if barcode != "" {
		return decodeBarcode(barcode)
	} else {
		log.Warn("No barcode entered.")
	}
	return ""
}

func ProcessInput() string{
	var input string
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		input = scanner.Text()
	}
	if len(input) > 0 {
		return input
	} else {
		log.Error("No input entered.")
	}
	return ""
}

func ProcessQuantity() float32{
	fmt.Print("Enter quantity: ")
	quantity := ProcessInput()
	if quantity != "" {
		q, err := strconv.ParseFloat(quantity, 32)
		if err != nil {
			log.Fatal(err)
		}
		return float32(q)
	} else {
		log.Error("No quantity entered.")
	}
	return 0
}

func ParseLocation() InvenTreeLocation{
	fmt.Print("Enter location: ")
	location := ProcessInput()
	if location != "" {
		return InvenTreeLocation{
			Name: location,
		}
	} else {
		log.Error("No location entered.")
	}
	// return empty location if error
	return InvenTreeLocation{}
}