package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"os/exec"
	"strings"
)

type DigiKeyPart struct {
	DigiKeyPartNumber      string   `json:"DigiKeyPartNumber"`
	Manufacturer           string   `json:"Manufacturer"`
	ManufacturerPartNumber string   `json:"ManufacturerPartNumber"`
	ProductDescription     string   `json:"ProductDescription"`
	LimitedTaxonomy        []string `json:"LimitedTaxonomy"`
	PrimaryPhoto           string   `json:"PrimaryPhoto"`
	ProductUrl             string   `json:"ProductUrl"`
	DetailedDescription    string   `json:"DetailedDescription"`
}

func (d *DigiKeyPart) String() string {
	return fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s,%s", d.DigiKeyPartNumber, d.Manufacturer, d.ManufacturerPartNumber, d.ProductDescription, strings.Join(d.LimitedTaxonomy, ","), d.PrimaryPhoto, d.ProductUrl, d.DetailedDescription)
}

func (d *DigiKeyPart) splitTaxonomy() {
	splitTaxonomy := strings.Split(d.LimitedTaxonomy[0], " - ")
	splitTaxonomy = append(splitTaxonomy, d.LimitedTaxonomy[1])
	d.LimitedTaxonomy = append([]string{splitTaxonomy[len(splitTaxonomy)-1]}, splitTaxonomy[:len(splitTaxonomy)-1]...)
}

func (d *DigiKeyPart) extractValues(response map[string]interface{}) {
	for key, value := range response {
		switch v := value.(type) {
		case map[string]interface{}:
			d.extractValues(v)
		case []interface{}:
			for _, item := range v {
				if m, ok := item.(map[string]interface{}); ok {
					d.extractValues(m)
				}
			}
		case string:
			switch key {
			case "Value":
				if response["Parameter"] == "Categories" {
					d.LimitedTaxonomy = append(d.LimitedTaxonomy, v)
				} else if response["Parameter"] == "Manufacturer" {
					d.Manufacturer = v
				}
			case "ProductUrl":
				d.ProductUrl = v
			case "PrimaryPhoto":
				d.PrimaryPhoto = v
			case "DetailedDescription":
				d.DetailedDescription = v
			case "ManufacturerPartNumber":
				d.ManufacturerPartNumber = v
			case "DigiKeyPartNumber":
				d.DigiKeyPartNumber = v
			case "ProductDescription":
				d.ProductDescription = v
			}
		}
	}
}

func (d *DigiKeyPart) parseResponse(response map[string]interface{}) {
	log.Debug("Parsing response from Digi-Key API")
	d.extractValues(response)
	d.splitTaxonomy()
	log.Info("Processed Part ", d.ProductDescription)
}

func (d *DigiKeyPart) PrettyPrint() {
	fmt.Println("Digi-Key Part Number:", d.DigiKeyPartNumber)
	fmt.Println("Manufacturer Part Number:", d.ManufacturerPartNumber)
	fmt.Println("Manufacturer:", d.Manufacturer)
	fmt.Println("Product Description:", d.ProductDescription)
	fmt.Println("Detailed Description:", d.DetailedDescription)
	fmt.Println("Limited Taxonomy:", d.LimitedTaxonomy)
	fmt.Println("Primary Photo:", d.PrimaryPhoto)
	fmt.Println("Product URL:", d.ProductUrl)
}

func NewDigiKeyPart(response map[string]interface{}) *DigiKeyPart {
	d := new(DigiKeyPart)
	d.parseResponse(response)
	return d
}

func WriteLabels(d *DigiKeyPart) {
	// so I gave up on implementing WriteLabel in Go and just called the python script
	// I'm sorry

	// install python dependencies from requirements.txt
	cmd_dep := exec.Command("pip", "install", "-r", "requirements.txt")
	err_dep := cmd_dep.Run()
	if err_dep != nil {
		log.Fatal(err_dep)
	}
	cmd := exec.Command("python", "labelwriter.py", "-m", d.ManufacturerPartNumber, "-l", d.LimitedTaxonomy[len(d.LimitedTaxonomy)-2], "-d", d.ProductDescription)
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}
