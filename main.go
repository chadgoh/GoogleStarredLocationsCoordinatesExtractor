package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

type MapData struct {
	Type     string `json:"type"`
	Features []struct {
		Geometry struct {
			Coordinates []float64 `json:"coordinates"`
			Type        string    `json:"type"`
		} `json:"geometry"`
		Properties struct {
			Date          time.Time `json:"date"`
			GoogleMapsUrl string    `json:"google_maps_url"`
			Location      struct {
				Address     string `json:"address"`
				CountryCode string `json:"country_code,omitempty"`
				Name        string `json:"name"`
			} `json:"location,omitempty"`
			Comment string `json:"Comment,omitempty"`
		} `json:"properties"`
		Type string `json:"type"`
	} `json:"features"`
}

// simple program to extract longitude and latitude data from starred locations on google maps
// starred locations are extracted from Google Take
func main() {
	jsonFile, err := os.Open("mapdata.json")
	if err != nil {
		fmt.Println("Error opening file: ", err)
	}
	fmt.Println("Successfully Opened mapdata.json")
	defer jsonFile.Close()

	// parse file into mapdata struct
	var mapData MapData
	byteValue, _ := ioutil.ReadAll(jsonFile)
	json.Unmarshal(byteValue, &mapData)

	csvFile, err := os.Create("mapdata.csv")

	if err != nil {
		fmt.Println("Error creating csv file: ", err)
	}

	defer csvFile.Close()

	writer := csv.NewWriter(csvFile)

	defer writer.Flush()

	headers := []string{"Name", "Longitude", "Latitude"}
	err = writer.Write(headers)
	if err != nil {
		fmt.Println("Error writing headers: ", err)
	}
	// extract name and coordinates
	for _, feature := range mapData.Features {
		if feature.Properties.Location.CountryCode == "NZ" {
			// write to csv
			var record []string

			record = append(record, feature.Properties.Location.Name)
			record = append(record, fmt.Sprintf("%f", feature.Geometry.Coordinates[0]))
			record = append(record, fmt.Sprintf("%f", feature.Geometry.Coordinates[1]))
			writer.Write(record)

		} else if feature.Properties.Comment == "No location information is available for this saved place" {
			fmt.Println("No location data available for: ", feature.Properties.GoogleMapsUrl)
		}

	}

}
