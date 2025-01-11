package geonames

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rusystem/crm-api/pkg/domain"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	Username = "demodemo"
	BaseURL  = "http://api.geonames.org"
)

var (
	countryPath = "/countryInfoJSON"
	searchPath  = "/searchJSON"
	dataPath    = "/getJSON"
)

type Client struct {
	BaseURL    string
	Username   string
	HTTPClient *http.Client
}

func NewGeonamesClient(timeout time.Duration) (*Client, error) {
	if timeout == 0 {
		return nil, errors.New("timeout can`t be zero")
	}

	return &Client{
		BaseURL:  BaseURL,
		Username: Username,
		HTTPClient: &http.Client{
			Timeout: timeout,
			Transport: &loggingRoundTripper{
				logger: os.Stdout,
				next:   http.DefaultTransport,
			},
		},
	}, nil
}

func (c *Client) sendRequest(req *http.Request, data interface{}) error {
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json; charset=utf-8")

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		var errResp ErrorResponse
		if err = json.NewDecoder(resp.Body).Decode(&errResp); err == nil {
			return errors.New(errResp.Info())
		}

		return fmt.Errorf("unknown error, status code %d\n", resp.StatusCode)
	}

	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil
	}

	return nil
}

func (c *Client) FetchCountries(ctx context.Context, lang string) ([]domain.Country, error) {
	var countries []domain.Country

	params := url.Values{}
	params.Add("username", c.Username)
	params.Add("lang", lang)

	fullUrl := fmt.Sprintf("%s%s?%s", c.BaseURL, countryPath, params.Encode())

	req, err := http.NewRequest("GET", fullUrl, nil)
	if err != nil {
		return nil, err
	}

	req.WithContext(ctx)

	var cr CountriesResponse
	if err = c.sendRequest(req, &cr); err != nil {
		return nil, err
	}

	for _, country := range cr.Geonames {
		countries = append(countries, domain.Country{
			Id:   country.CountryID,
			Name: country.CountryName,
		})
	}

	sort.Slice(countries, func(i, j int) bool {
		return countries[i].Name < countries[j].Name
	})

	return countries, err
}

func (c *Client) FetchRegions(ctx context.Context, code, lang string) ([]domain.Region, error) {
	citiesMap := make(map[string]domain.Region)

	params := url.Values{}
	params.Add("username", c.Username)
	params.Add("lang", lang)
	params.Add("country", strings.ToUpper(code))
	params.Add("featureCode", "ADM1")
	params.Add("maxRows", "1000")

	startRow := 0

	for {
		params.Set("startRow", fmt.Sprintf("%d", startRow))
		fullUrl := fmt.Sprintf("%s%s?%s", c.BaseURL, searchPath, params.Encode())

		req, err := http.NewRequest("GET", fullUrl, nil)
		if err != nil {
			return nil, err
		}

		req = req.WithContext(ctx)

		// Выполняем запрос и декодируем ответ
		var cr RegionResponse
		if err = c.sendRequest(req, &cr); err != nil {
			return nil, err
		}

		// Если нет данных, выходим из цикла
		if len(cr.Geonames) == 0 {
			break
		}

		// Обрабатываем полученные города
		for _, region := range cr.Geonames {
			if !isValidGeoName(region.Name) {
				continue
			}

			// Уникальный ключ для городов
			key := fmt.Sprintf("%s_%s", region.Name, region.AdminCode1)

			// Добавляем только уникальные города
			if _, exists := citiesMap[key]; !exists {
				citiesMap[key] = domain.Region{
					Name: region.Name,
					Id:   region.AdminCode1,
				}
			}
		}

		startRow += 1000
	}

	var cities []domain.Region
	for _, city := range citiesMap {
		cities = append(cities, city)
	}

	sort.Slice(cities, func(i, j int) bool {
		return cities[i].Name < cities[j].Name
	})

	return cities, nil
}

func (c *Client) FetchCitiesByRegion(ctx context.Context, countryCode, adminCode1, lang string) ([]domain.City, error) {
	var cities []domain.City
	nameSet := make(map[string]bool)
	startRow := 0

	params := url.Values{}
	params.Add("adminCode1", adminCode1)
	params.Add("country", countryCode)
	params.Add("featureClass", "P")
	params.Add("lang", lang)
	params.Add("username", c.Username)
	params.Add("orderby", "population")
	params.Add("maxRows", "1000")
	params.Add("startRow", strconv.Itoa(startRow))

	fullURL := fmt.Sprintf("%s%s?%s", c.BaseURL, searchPath, params.Encode())

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	req.WithContext(ctx)

	for {
		var resp CityResponse
		if err = c.sendRequest(req, &resp); err != nil {
			return nil, err
		}

		for _, geo := range resp.Geonames {
			if _, exists := nameSet[geo.Name]; exists {
				continue
			}

			if !isValidGeoName(geo.Name) {
				continue
			}

			nameSet[geo.Name] = true
			cities = append(cities, domain.City{
				Name: geo.Name,
				Id:   geo.GeonameId,
			})
		}

		// Переход к следующей странице
		startRow += 1000
		if startRow >= resp.TotalResultsCount {
			break
		}
	}

	// Сортировка в алфавитном порядке
	sort.Slice(cities, func(i, j int) bool {
		return cities[i].Name < cities[j].Name
	})

	return cities, nil
}

// isValidGeoName проверяет, содержит ли название города английские символы
func isValidGeoName(name string) bool {
	re := regexp.MustCompile(`[a-zA-Z]`)
	return !re.MatchString(name)
}
