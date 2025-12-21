package weather

import (
	"fmt"
	"log"
)

type WeatherAPI struct{}

type RawWeather struct {
	TempF     float64
	WindMPH   float64
	Condition string
}

func (w *WeatherAPI) GetWeather(city string) (RawWeather, error) {
	return RawWeather{
		TempF:     77.0,
		WindMPH:   10.0,
		Condition: "Sunny",
	}, nil
}

type Weather struct {
	TempC     float64
	WindMPS   float64
	Condition string
}

type WeatherProvider interface {
	Get(city string) (Weather, error)
}

// Адаптер

type WeatherAdapter struct {
	api *WeatherAPI
}

func NewWeatherAdapter(api *WeatherAPI) *WeatherAdapter {
	return &WeatherAdapter{api: api}
}

func (a *WeatherAdapter) Get(city string) (Weather, error) {
	raw, err := a.api.GetWeather(city)
	if err != nil {
		return Weather{}, err
	}

	return Weather{
		TempC:     fahrenheitToCelsius(raw.TempF),
		WindMPS:   mphToMps(raw.WindMPH),
		Condition: raw.Condition,
	}, nil
}

func fahrenheitToCelsius(f float64) float64 {
	return (f - 32) * 5.0 / 9.0
}

func mphToMps(mph float64) float64 {
	return mph * 0.44704
}

func Example() {
	api := &WeatherAPI{}
	provider := NewWeatherAdapter(api)

	w, err := provider.Get("Frankfurt")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Temp: %.1f°C, wind: %.1f m/s, condition: %s\n",
		w.TempC, w.WindMPS, w.Condition)
}
