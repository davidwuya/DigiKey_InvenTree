package main

import (
	"encoding/json"
	"github.com/go-resty/resty/v2"
	log "github.com/sirupsen/logrus"
	"strings"
)

type InvenTreeCompany struct {
	Pk              int    `json:"pk"`
	Name            string `json:"name"`
	Is_supplier     bool   `json:"is_supplier"`
	Is_manufacturer bool   `json:"is_manufacturer"`
	Currency        string `json:"currency"`
}

func GetAllInvenTreeCompanies() []InvenTreeCompany {
	client := resty.New()
	client.SetRetryCount(3)
	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetAuthScheme("Token").
		SetAuthToken(INVENTREE_TOKEN).
		Get(InvenTreeServerAddress + "company/")
	if err != nil {
		log.Fatal(err)
	}
	var companies []InvenTreeCompany
	err = json.Unmarshal(resp.Body(), &companies)
	if err != nil {
		log.Fatal(err)
	}
	return companies
}

func GetInvenTreeCompanyByName(name string, companies []InvenTreeCompany) (InvenTreeCompany, []InvenTreeCompany) {
	for _, c := range companies {
		if strings.EqualFold(c.Name, name) {
			return c, companies
		}
	}
	log.Error("Company not found: ", name)
	c:= CreateInvenTreeCompany(&InvenTreeCompany{
		Name:            name,
		Is_supplier:     false,
		Is_manufacturer: true,
		Currency:        "USD",
	})
	companyCache := append(companies, c)
	return c, companyCache
}

func CreateInvenTreeCompany(c *InvenTreeCompany) InvenTreeCompany {
	client := resty.New()
	client.SetRetryCount(3)
	resp, err := client.R().
		SetAuthScheme("Token").
		SetAuthToken(INVENTREE_TOKEN).
		// POST to company endpoint, do not send pk
		SetBody(InvenTreeCompany{
			Name:            c.Name,
			Is_supplier:     c.Is_supplier,
			Is_manufacturer: c.Is_manufacturer,
			Currency:        c.Currency,
		}).
		Post(InvenTreeServerAddress + "company/")
	if err != nil {
		log.Fatal(err)
	}

	var company InvenTreeCompany
	err = json.Unmarshal(resp.Body(), &company)
	if err != nil {
		log.Fatal(err)
	}
	log.Info("Created company: ", company.Name)
	return company

}
