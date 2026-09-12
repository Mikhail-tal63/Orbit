package geocoding

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

type GeoapifyService struct {
	apiKey string
	client *http.Client
}

func NewGeoapifyService(apiKey string) *GeoapifyService {
	return &GeoapifyService{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

type ReverseGeocodeResponse struct {
	Results []struct {
		Formated string `json:"formatted"`
	} `json:"results"`
}

func (s *GeoapifyService) ReverseGeocode(ctx context.Context, lat, lon float64) (string, error) {

	endpoint := "https://api.geoapify.com/v1/geocode/reverse"

	params := url.Values{}
	params.Set("lat", fmt.Sprintf("%f", lat))
	params.Set("lon", fmt.Sprintf("%f", lon))
	params.Set("format", "json")
	params.Set("apiKey", s.apiKey)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint+"?"+params.Encode(), nil)
	if err != nil {
		return "", err
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf(
			"geoapify returned status %d",
			resp.StatusCode,
		)
	}

	var result ReverseGeocodeResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Results) == 0 {
		return "", fmt.Errorf("no address found")
	}

	return result.Results[0].Formated, nil
}
