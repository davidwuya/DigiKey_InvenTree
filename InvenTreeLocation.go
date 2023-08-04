package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"strings"
	"github.com/go-resty/resty/v2"
	"fmt"
)

type InvenTreeLocation struct {
	Pk         int    `json:"pk"`
	Name       string `json:"name"`
	Parent     int    `json:"parent"` // primary key of parent location
	PathString string `json:"pathstring"`
}


// GetAllInvenTreeLocations retrieves all InvenTree locations from the server and returns them as a slice of InvenTreeLocation structs.
// The PathString field of each location is modified to only include the last five characters.
func GetAllInvenTreeLocations() []InvenTreeLocation {
	client := resty.New()
	client.SetRetryCount(3)
	resp, err := client.R().
				SetAuthScheme("Token").
		SetAuthToken(INVENTREE_TOKEN).
		
		SetHeader("Accept", "application/json").
		Get(InvenTreeServerAddress + "stock/location/")
	if err != nil {
		log.Fatal(err)
	}
	var locations []InvenTreeLocation
	err = json.Unmarshal(resp.Body(), &locations)
	if err != nil {
		log.Fatal(err)
	}
	// Iterate over the locations and modify the PathString
	for i, location := range locations {
		// Check if PathString has more than five characters
		if len(location.PathString) > 5 {
			// Keep only the last five characters
			locations[i].PathString = location.PathString[len(location.PathString)-5:]
		}
	}
	return locations
}

func GetInvenTreeLocationByName(LocationName string, locations []InvenTreeLocation) InvenTreeLocation{
	for _, location := range locations {
		if strings.Replace(location.PathString, "/", "", 1) == LocationName {
			return location
		}
	}
	return InvenTreeLocation{}
}

func InvenTreeLocationDebugPrint(location InvenTreeLocation) {
	fmt.Println("Location Name: ", location.Name)
}

func GetInvenTreeLocationNameByPK(LocationPK int, locations []InvenTreeLocation) string {
	for _, location := range locations {
		if location.Pk == LocationPK {
			return location.PathString
		}
	}
	return ""
}