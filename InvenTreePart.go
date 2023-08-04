package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"io"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

type InvenTreePart struct {
	Pk                int    `json:"pk"` // primary key of part
	Name              string `json:"name"`
	Description       string `json:"description"`
	Category          int    `json:"category"`
	IPN               string `json:"IPN"` // Manufacturer Part Number
	Active            bool   `json:"active"`
	Component         bool   `json:"component"`
	Virtual           bool   `json:"virtual"`
	Purchasable       bool   `json:"purchasable"`
	Assembly          bool   `json:"assembly"`
	DigiKeyPartNumber string
	PrimaryUrl        string
	PrimaryPhoto      string
	Manufacturer      int `json:"manufacturer"` // primary key of associated manufacturer
	SKU               string
}

type SupplierPart struct {
	Part             int    `json:"part"`
	Supplier         int    `json:"supplier"`
	SKU              string `json:"SKU"`
	Manufacturer     string `json:"manufacturer"`
	Description      string `json:"description"`
	Link             string `json:"link"`
	ManufacturerPart int    `json:"manufacturer_part"`
}

type ManufacturerPart struct {
	Pk           int    `json:"pk"`           // primary key of manufacturer part
	Part         int    `json:"part"`         // primary key of associated part
	Manufacturer int    `json:"manufacturer"` // primary key of associated manufacturer
	MPN          string `json:"MPN"`          // Manufacturer Part Number
	Description  string `json:"description"`
	Link         string `json:"link"`
}

func GetAllInvenTreeParts() []InvenTreePart {
	client := resty.New()
	client.SetRetryCount(3)
	resp, err := client.R().
		SetAuthScheme("Token").
		SetAuthToken(INVENTREE_TOKEN).
		SetHeader("Accept", "application/json").
		Get(InvenTreeServerAddress + "part/")
	if err != nil {
		log.Fatal(err)
	}
	var parts []InvenTreePart
	err = json.Unmarshal(resp.Body(), &parts)
	if err != nil {
		log.Fatal(err)
	}
	return parts
}

func GetInvenTreePartByName(name string, parts []InvenTreePart) InvenTreePart {
	for _, part := range parts {
		if part.Name == name {
			return part
		}
	}
	return InvenTreePart{}
}

func GetInvenTreePartByPK(pk int, parts []InvenTreePart) InvenTreePart {
	for _, part := range parts {
		if part.Pk == pk {
			return part
		}
	}
	return InvenTreePart{}
}
func GetInvenTreePartByMPN(mpn string, parts []InvenTreePart) InvenTreePart {
	for _, part := range parts {
		if part.IPN == mpn {
			return part
		}
	}
	return InvenTreePart{}
}



func UploadThumbnailToInvenTreePart(part *InvenTreePart, d *DigiKeyPart) bool {
	ImageURL := d.PrimaryPhoto
	TempFile := "temp.jpg"
	log.Debug("Downloading image: ", ImageURL)
	// Download the image
	ImageClient := &http.Client{}
	req, err := http.NewRequest("GET", ImageURL, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.36")
	resp, err := ImageClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Create the temporary file
	log.Trace("Creating temporary file: ", TempFile)
	out, err := os.Create(TempFile)
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	// Write the downloaded image to the file
	log.Trace("Writing image to temporary file: ", TempFile)
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	// Open the temporary image file
	log.Trace("Reading temporary file: ", TempFile)
	file, err := os.Open(TempFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	EndPoint := fmt.Sprintf("part/thumbs/%d/", part.Pk)

	// Upload image to InvenTree
	log.Trace("Uploading thumbnail to InvenTree: ", EndPoint)
	client := resty.New()
	resp2, err := client.R().
		SetFileReader("image", "image.jpg", file).
		SetAuthScheme("Token").
		SetAuthToken(INVENTREE_TOKEN).
		Put(InvenTreeServerAddress + EndPoint)

	if err != nil {
		log.Fatal(err)
		

	}
	if resp2.StatusCode() == 200 {
		log.Info("Uploaded thumbnail to InvenTree")
		return true
	} else {
		log.Error("Failed to upload thumbnail to InvenTree")
		return false
	}
}

func CreateInvenTreePart(d *DigiKeyPart, categories []InvenTreePartCategory, companies []InvenTreeCompany) (InvenTreePart, []InvenTreeCompany) {
	category := FindInvenTreeCategoryFromDigiKeyPart(d, categories)
	log.Trace("Category: ", category)
	manufacturer, companies := GetInvenTreeCompanyByName(d.Manufacturer, companies)
	log.Trace("Manufacturer: ", manufacturer)

	client := resty.New()
	resp, err := client.R().
		SetAuthScheme("Token").
		SetAuthToken(INVENTREE_TOKEN).
		SetBody(InvenTreePart{
			Name:         d.ProductDescription,
			Description:  d.DetailedDescription,
			IPN:          d.ManufacturerPartNumber,
			Active:       true,
			Component:    true,
			Assembly:     false,
			Virtual:      false,
			Category:     category.Pk,
			Manufacturer: manufacturer.Pk,
		}).
		Post(InvenTreeServerAddress + "part/")

	if err != nil {
		log.Fatal(err)
	}
	var part InvenTreePart
	log.Trace("Response: ", resp.String())
	err = json.Unmarshal(resp.Body(), &part)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Created InvenTree part: ", part.Name)
	return part, companies
}

func InvenTreePartDebugPrint(part InvenTreePart) {
	fmt.Println("InvenTreePart:")
	fmt.Println("PK: ", part.Pk)
	fmt.Println("Name: ", part.Name)
	fmt.Println("Manufacturer: ", part.Manufacturer)
	fmt.Println("Description: ", part.Description)
	fmt.Println("IPN: ", part.IPN)
	fmt.Println("Category: ", part.Category)
	fmt.Println("Active: ", part.Active)
	fmt.Println("Component: ", part.Component)
	fmt.Println("Assembly: ", part.Assembly)
	fmt.Println("Virtual: ", part.Virtual)
}

func CreateManufacturerPart(part InvenTreePart, d *DigiKeyPart, companies []InvenTreeCompany) ManufacturerPart {
	client := resty.New()
	manufaturer, _ := GetInvenTreeCompanyByName(d.Manufacturer, companies)
	log.Trace("CreateManufacturerPart Manufacturer: ", manufaturer)
	log.Trace("CreateManufacturerPart Part: ", part)
	resp, err := client.R().
		SetAuthScheme("Token").
		SetAuthToken(INVENTREE_TOKEN).
		SetBody(ManufacturerPart{
			Part:         part.Pk,
			Manufacturer: manufaturer.Pk,
			MPN:          part.IPN,
			Description:  part.Description,
			Link:         d.ProductUrl,
		}).
		Post(InvenTreeServerAddress + "company/part/manufacturer/")
	if err != nil {
		log.Fatal(err)
	}
	log.Trace("CreateManufacturerPart Response: ", resp.String())
	var manufacturerPart ManufacturerPart
	err = json.Unmarshal(resp.Body(), &manufacturerPart)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Created InvenTree manufacturer part: ", manufacturerPart.MPN)
	return manufacturerPart
}

func CreateSupplierPart(part *InvenTreePart, mfr_part *ManufacturerPart, d *DigiKeyPart) SupplierPart {
	client := resty.New()
	resp, err := client.R().
		SetAuthScheme("Token").
		SetAuthToken(INVENTREE_TOKEN).
		SetBody(SupplierPart{
			Part:             part.Pk,
			Supplier:         1, // Digi-Key
			SKU:              d.DigiKeyPartNumber,
			Description:      part.Description,
			Link:             d.ProductUrl,
			ManufacturerPart: mfr_part.Pk,
		}).
		Post(InvenTreeServerAddress + "company/part/")
	if err != nil {
		log.Fatal(err)
	}
	var supplierPart SupplierPart
	err = json.Unmarshal(resp.Body(), &supplierPart)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Created InvenTree supplier part: ", supplierPart.Description)
	return supplierPart
}

func GetSupplierPartByName(d *DigiKeyPart) SupplierPart {
	client := resty.New()
	resp, err := client.R().
		SetAuthScheme("Token").
		SetAuthToken(INVENTREE_TOKEN).
		Get(InvenTreeServerAddress + "company/part/?MPN=" + d.ManufacturerPartNumber)
	if err != nil {
		log.Fatal(err)
	}
	var supplierParts []SupplierPart
	err = json.Unmarshal(resp.Body(), &supplierParts)
	if err != nil {
		log.Fatal(err)
	}
	if len(supplierParts) == 0 {
		return SupplierPart{}
	}
	return supplierParts[0]
}

func SupplierPartDebugPrint(s SupplierPart) {
	fmt.Println("SupplierPart:")
	fmt.Println("Part: ", s.Part)
	fmt.Println("Supplier: ", s.Supplier)
	fmt.Println("SKU: ", s.SKU)
	fmt.Println("Description: ", s.Description)
	fmt.Println("Link: ", s.Link)
	fmt.Println("ManufacturerPart: ", s.ManufacturerPart)
}

func GetManufacturerPartByName(d *DigiKeyPart) ManufacturerPart{
	client := resty.New()
	resp, err := client.R().
		SetAuthScheme("Token").
		SetAuthToken(INVENTREE_TOKEN).
		Get(InvenTreeServerAddress + "company/part/manufacturer/?MPN=" + d.ManufacturerPartNumber)
	if err != nil {
		log.Fatal(err)
	}
	var manufacturerParts []ManufacturerPart
	err = json.Unmarshal(resp.Body(), &manufacturerParts)
	if err != nil {
		log.Fatal(err)
	}
	if len(manufacturerParts) == 0 {
		return ManufacturerPart{}
	}
	return manufacturerParts[0]
}

func ManufacturerPartDebugPrint(m ManufacturerPart) {
	fmt.Println("ManufacturerPart:")
	fmt.Println("PK", m.Pk)
	fmt.Println("Part: ", m.Part)
	fmt.Println("Manufacturer: ", m.Manufacturer)
	fmt.Println("MPN: ", m.MPN)
	fmt.Println("Description: ", m.Description)
	fmt.Println("Link: ", m.Link)
}