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
type Route struct{
	DistanceKm  float64
    DurationMin int
}
func NewGeoapifyService(apiKey string) *GeoapifyService {
	return &GeoapifyService{
		apiKey: apiKey,
		client: &http.Client{},
	}
}

type ReverseGeocodeResponse struct {
	Results []struct {
		Formated string `json:"formated"`
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

func (s *GeoapifyService) GetRoute(ctx context.Context,
	pickupLat ,pickupLon float64,
	 destinationLat, destinationLon float64)(Route,error){
	    url := fmt.Sprintf(
        "https://api.geoapify.com/v1/routing?waypoints=%f,%f|%f,%f&mode=drive&apiKey=%s",
        pickupLat,
        pickupLon,
        destinationLat,
        destinationLon,
        s.apiKey,

    )
req ,err := http.NewRequestWithContext(ctx,http.MethodGet,url,nil)
if err!= nil {
	return Route{},err
}
resp,err := s.client.Do(req)
if err != nil {
	return Route{},err
}
defer req.Body.Close()

if resp.StatusCode != http.StatusOK{
	    return Route{}, fmt.Errorf(
            "geoapify routing returned status %d",
            resp.StatusCode,
        )
}
    var result struct {
        Features []struct {
            Properties struct {
                Distance float64 `json:"distance"`
                Time     float64 `json:"time"`
            } `json:"properties"`
        } `json:"features"`
    }
	   if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return Route{}, err
    }

    if len(result.Features) == 0 {
        return Route{}, fmt.Errorf("no route found")
    }

    distanceMeters := result.Features[0].Properties.Distance
    durationSeconds := result.Features[0].Properties.Time

    return Route{
        DistanceKm:  distanceMeters / 1000,
        DurationMin: int(durationSeconds / 60),
    }, nil
	 }
