package geo

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Geo struct {
	Api string
	ErrorLog *log.Logger
	InfoLog *log.Logger
}

type GeoLocationResponse struct {
	Status        string `json:"status"` // "successs" or "fail"
	Query         string `json:"query"` // "IP address"
	CountryCode   string `json:"CountryCode"` // ex. "US"
	CountryName   string `json:"CountryName"` // ex. "United States"
	Capital       string `json:"Capital"` // ex. "Washington D.C." 
	PhonePrefix   string `json:"PhonePrefix"` // ex. "+1"
	Currency      string `json:"Currency"` // ex. "USD"
	USDRate       string `json:"USDRate"` // ex. "1"
	EURRate       string `json:"EURRate"` // ex. "1.05"
	RegionCode    string `json:"RegionCode"` // ex. "FL"
	RegionName    string `json:"RegionName"` // ex. "Florida"
	City          string `json:"City"` // ex. "Miami"
	Postal        string `json:"Postal"` // ex. "33101"
	Latitude      string `json:"Latitude"` // ex. "25.7743"
	Longitude     string `json:"Longitude"` // ex. "-80.1937"
	TimeZone      string `json:"TimeZone"` // ex. "America/New_York"
	ContinentCode string `json:"ContinentCode"` // ex. "NA"
	ContinentName string `json:"ContinentName"` // ex. "North America"
	ASN           string `json:"asn"` // ex. "AS15169"
	Org           string `json:"org"` // ex. "Google LLC"
}

func (g *Geo) GetGeoLocation(ip string) string {
	apiUrl := fmt.Sprintf("%s%s", g.Api, ip)

	resp, err := http.Get(apiUrl)
	if err != nil {
		g.ErrorLog.Printf("Error fetching geolocation for %s: %v", ip, err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK { 
		g.InfoLog.Printf("Status code not 200 for %s: %v", ip, resp.StatusCode)
	}

	var geo_location GeoLocationResponse
	if err := json.NewDecoder(resp.Body).Decode(&geo_location); err != nil {
		g.ErrorLog.Printf("Error decoding geolocation response for %s: %v", ip, err)
		return ""
	}

	if geo_location.Status == "fail" {
		g.ErrorLog.Printf("Failed to fetch geolocation for %s: %v", ip, err)
		return ""
	}
	location := fmt.Sprintf("%s, %s", geo_location.City, geo_location.RegionName)

	return location
}