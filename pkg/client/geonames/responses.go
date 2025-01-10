package geonames

import "fmt"

type ErrorResponse struct {
	Message     string `json:"message"`
	Description string `json:"description"`
	Error       string `json:"error"`
}

type SuccessResponse struct {
	Href      string `json:"href"`
	Method    string `json:"method"`
	Templated bool   `json:"templated"`
}

func (e *ErrorResponse) Info() string {
	return fmt.Sprintf("message: %s, description: %s, error: %s\n", e.Message, e.Description, e.Error)
}

type Country struct {
	CountryID   string `json:"countryCode"`
	CountryName string `json:"countryName"`
}

type CountriesResponse struct {
	Geonames []Country `json:"geonames"`
}

type Region struct {
	AdminCode1 string `json:"adminCode1"`
	Name       string `json:"name"`
}

type RegionResponse struct {
	Geonames []Region `json:"geonames"`
}

type CityResponse struct {
	Geonames []struct {
		Name        string `json:"name"`
		ToponymName string `json:"toponymName"`
		GeonameId   int64  `json:"geonameId"`
	} `json:"geonames"`
	TotalResultsCount int `json:"totalResultsCount"`
}
