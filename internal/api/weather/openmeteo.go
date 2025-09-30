package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Rochester, NY coordinates
const (
	RochesterLat = 43.1566
	RochesterLon = -77.6088
)

const openMeteoBaseURL = "https://api.open-meteo.com/v1/forecast"

// Client is a minimal HTTP client for Open-Meteo.
// It requests current and daily weather for a fixed location (Rochester, NY).
type Client struct {
	HTTP *http.Client
}

func NewClient() *Client {
	return &Client{HTTP: &http.Client{Timeout: 7 * time.Second}}
}

// OpenMeteoResponse models the subset of fields we need from Open-Meteo.
// See https://open-meteo.com/en/docs for full schema.
type OpenMeteoResponse struct {
	Latitude        float64 `json:"latitude"`
	Longitude       float64 `json:"longitude"`
	Timezone        string  `json:"timezone"`
	TimezoneAbbrev  string  `json:"timezone_abbreviation"`
	Elevation       float64 `json:"elevation"`
	CurrentUnits    struct {
		Time                string `json:"time"`
		Interval            string `json:"interval"`
		Temperature2m       string `json:"temperature_2m"`
		RelativeHumidity2m  string `json:"relative_humidity_2m"`
		ApparentTemperature string `json:"apparent_temperature"`
		IsDay               string `json:"is_day"`
		Precipitation       string `json:"precipitation"`
		WeatherCode         string `json:"weather_code"`
		CloudCover          string `json:"cloud_cover"`
		PressureMSL         string `json:"pressure_msl"`
		WindSpeed10m        string `json:"wind_speed_10m"`
		WindDirection10m    string `json:"wind_direction_10m"`
		WindGusts10m        string `json:"wind_gusts_10m"`
	} `json:"current_units"`
	Current struct {
		Time                string  `json:"time"`
		Interval            int     `json:"interval"`
		Temperature2m       float64 `json:"temperature_2m"`
		RelativeHumidity2m  float64 `json:"relative_humidity_2m"`
		ApparentTemperature float64 `json:"apparent_temperature"`
		IsDay               int     `json:"is_day"`
		Precipitation       float64 `json:"precipitation"`
		WeatherCode         int     `json:"weather_code"`
		CloudCover          float64 `json:"cloud_cover"`
		PressureMSL         float64 `json:"pressure_msl"`
		WindSpeed10m        float64 `json:"wind_speed_10m"`
		WindDirection10m    float64 `json:"wind_direction_10m"`
		WindGusts10m        float64 `json:"wind_gusts_10m"`
	} `json:"current"`
	DailyUnits struct {
		Time              string `json:"time"`
		TempMax           string `json:"temperature_2m_max"`
		TempMin           string `json:"temperature_2m_min"`
		Sunrise           string `json:"sunrise"`
		Sunset            string `json:"sunset"`
		UVIndexMax        string `json:"uv_index_max"`
	} `json:"daily_units"`
	Daily struct {
		Time       []string  `json:"time"`
		TempMax    []float64 `json:"temperature_2m_max"`
		TempMin    []float64 `json:"temperature_2m_min"`
		Sunrise    []string  `json:"sunrise"`
		Sunset     []string  `json:"sunset"`
		UVIndexMax []float64 `json:"uv_index_max"`
	} `json:"daily"`
}

// APIResponse is the simplified payload our app returns to clients.
type APIResponse struct {
	Location struct {
		Name      string  `json:"name"`
		Latitude  float64 `json:"latitude"`
		Longitude float64 `json:"longitude"`
	} `json:"location"`
	Units struct {
		Temperature      string `json:"temperature"`
		WindSpeed        string `json:"wind_speed"`
		Precipitation    string `json:"precipitation"`
		Pressure         string `json:"pressure"`
		Humidity         string `json:"humidity"`
		CloudCover       string `json:"cloud_cover"`
		WindDirection    string `json:"wind_direction"`
		WindGusts        string `json:"wind_gusts"`
	} `json:"units"`
	Current struct {
		Time             string  `json:"time"`
		Temperature      float64 `json:"temperature"`
		ApparentTemp     float64 `json:"apparent_temperature"`
		IsDay            bool    `json:"is_day"`
		Humidity         float64 `json:"humidity"`
		Precipitation    float64 `json:"precipitation"`
		WeatherCode      int     `json:"weather_code"`
		CloudCover       float64 `json:"cloud_cover"`
		Pressure         float64 `json:"pressure"`
		WindSpeed        float64 `json:"wind_speed"`
		WindDirection    float64 `json:"wind_direction"`
		WindGusts        float64 `json:"wind_gusts"`
	} `json:"current"`
	Daily []struct {
		Date        string  `json:"date"`
		TempMax     float64 `json:"temperature_max"`
		TempMin     float64 `json:"temperature_min"`
		Sunrise     string  `json:"sunrise"`
		Sunset      string  `json:"sunset"`
		UVIndexMax  float64 `json:"uv_index_max"`
	} `json:"daily"`
}

// GetRochesterWeather fetches current and daily weather from Open-Meteo for Rochester, NY.
func (c *Client) GetRochesterWeather(ctx context.Context) (*APIResponse, error) {
	params := url.Values{}
	params.Set("latitude", fmt.Sprintf("%f", RochesterLat))
	params.Set("longitude", fmt.Sprintf("%f", RochesterLon))
	params.Set("current", "temperature_2m,relative_humidity_2m,apparent_temperature,is_day,precipitation,weather_code,cloud_cover,pressure_msl,wind_speed_10m,wind_direction_10m,wind_gusts_10m")
	params.Set("daily", "temperature_2m_max,temperature_2m_min,sunrise,sunset,uv_index_max")
	params.Set("temperature_unit", "fahrenheit")
	params.Set("windspeed_unit", "mph")
	params.Set("precipitation_unit", "inch")
	params.Set("timezone", "auto")

	u := openMeteoBaseURL + "?" + params.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("open-meteo: unexpected status %d", resp.StatusCode)
	}

	var om OpenMeteoResponse
	if err := json.NewDecoder(resp.Body).Decode(&om); err != nil {
		return nil, err
	}

	out := &APIResponse{}
	out.Location.Name = "Rochester, NY"
	out.Location.Latitude = om.Latitude
	out.Location.Longitude = om.Longitude

	// Units copied straight from Open-Meteo current_units and daily_units
	out.Units.Temperature = om.CurrentUnits.Temperature2m
	out.Units.WindSpeed = om.CurrentUnits.WindSpeed10m
	out.Units.Precipitation = om.CurrentUnits.Precipitation
	out.Units.Pressure = om.CurrentUnits.PressureMSL
	out.Units.Humidity = om.CurrentUnits.RelativeHumidity2m
	out.Units.CloudCover = om.CurrentUnits.CloudCover
	out.Units.WindDirection = om.CurrentUnits.WindDirection10m
	out.Units.WindGusts = om.CurrentUnits.WindGusts10m

	out.Current.Time = om.Current.Time
	out.Current.Temperature = om.Current.Temperature2m
	out.Current.ApparentTemp = om.Current.ApparentTemperature
	out.Current.IsDay = om.Current.IsDay == 1
	out.Current.Humidity = om.Current.RelativeHumidity2m
	out.Current.Precipitation = om.Current.Precipitation
	out.Current.WeatherCode = om.Current.WeatherCode
	out.Current.CloudCover = om.Current.CloudCover
	out.Current.Pressure = om.Current.PressureMSL
	out.Current.WindSpeed = om.Current.WindSpeed10m
	out.Current.WindDirection = om.Current.WindDirection10m
	out.Current.WindGusts = om.Current.WindGusts10m

	// Daily aggregation
	for i := range om.Daily.Time {
		d := struct {
			Date       string  `json:"date"`
			TempMax    float64 `json:"temperature_max"`
			TempMin    float64 `json:"temperature_min"`
			Sunrise    string  `json:"sunrise"`
			Sunset     string  `json:"sunset"`
			UVIndexMax float64 `json:"uv_index_max"`
		}{
			Date:       om.Daily.Time[i],
			TempMax:    om.Daily.TempMax[i],
			TempMin:    om.Daily.TempMin[i],
			Sunrise:    om.Daily.Sunrise[i],
			Sunset:     om.Daily.Sunset[i],
			UVIndexMax: om.Daily.UVIndexMax[i],
		}
		out.Daily = append(out.Daily, d)
	}

	return out, nil
}

