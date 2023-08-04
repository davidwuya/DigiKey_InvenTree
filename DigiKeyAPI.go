package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"os/exec"
)

type DigikeyAPI struct {
	VercelURL   string
	DKAuthorize string
	APIKey      string
	ClientID    string
	OAuthState  string
}

func NewDigikeyAPI(APIKey, ClientID, OAuthState string) *DigikeyAPI {
	VercelURL := VercelURL
	DKAuthorize := DKAuthorize
	if APIKey == "" || ClientID == "" || OAuthState == "" {
		log.Fatal("APIKey, ClientID, and OAuthState are required")
	}
	return &DigikeyAPI{
		VercelURL:   VercelURL,
		DKAuthorize: DKAuthorize,
		APIKey:      APIKey,
		ClientID:    ClientID,
		OAuthState:  OAuthState,
	}
}

func OAuthAuthorize(d *DigikeyAPI) bool {
	RedirectURI := d.VercelURL + "callback"
	params := url.Values{
		"response_type": {"code"},
		"client_id":     {d.ClientID},
		"redirect_uri":  {RedirectURI},
		"state":         {d.OAuthState},
	}
	// open the browser to the URL with os/exec
	cmd := exec.Command("chrome", d.DKAuthorize+"?"+params.Encode())
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	return true

}

func VerifyToken(d *DigikeyAPI) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", d.VercelURL+"verify", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("x-api-key", d.APIKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	log.Debug(resp.Status)
	return resp.Status
}

func GetToken(d *DigikeyAPI) string {
	client := &http.Client{}
	req, err := http.NewRequest("GET", d.VercelURL+"token", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("x-api-key", d.APIKey)
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	// unmarshall the JSON response and get the token in "access_token"
	var tokenResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		log.Fatal(err)
	}
	log.Trace("Auth Server response status code " , resp.StatusCode)
	token := tokenResponse["access_token"].(string)
	return token
}

func GetProductDetails(d *DigikeyAPI, PartNumber string) map[string]interface{} {
	client := &http.Client{}
	// escape any problematic characters in the PartNumber
	PartNumber = url.QueryEscape(PartNumber)
	log.Trace("Escaped Part Number: " + PartNumber)
	req, err := http.NewRequest("GET", DKProductDetailEndpoint+PartNumber+"?includes=DigiKeyPartNumber,Manufacturer,ManufacturerPartNumber,ProductDescription,LimitedTaxonomy,PrimaryPhoto,ProductUrl,DetailedDescription", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Add("Authorization", "Bearer "+GetToken(d))
	req.Header.Add("X-DIGIKEY-Client-Id", d.ClientID)
	req.Header.Add("X-DIGIKEY-Locale-Site", "US")
	req.Header.Add("X-DIGIKEY-Locale-Language", "en")
	req.Header.Add("X-DIGIKEY-Locale-Currency", "USD")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	// log the JSON response to the console
	var productResponse map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&productResponse); err != nil {
		log.Fatal(err)
	}
	return productResponse
}
