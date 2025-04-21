package bot

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

const URL string = "https://api.openweathermap.org/data/2.5/weather?"
const GEOURL string = "http://api.openweathermap.org/geo/1.0/zip?"

type WeatherData struct {
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temp      float64 `json:"temp"`
		Humidity  int     `json:"humidity"`
		FeelsLike float64 `json:"feels_like"`
	} `json:"main"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
	Name string `json:"name"`
}

type CoordinatesData struct {
	Zip     string  `json:"zip"`
	Lat     float64 `json:"lat"`
	Lon     float64 `json:"lon"`
	Country string  `json:"country"`
	Name    string  `json:"name"`
}

func getCurrentWeather(message string) *discordgo.MessageSend {
	// Match 5-digit US ZIP code
	r, _ := regexp.Compile(`!zip\s+(.+)$`)
	matches := r.FindStringSubmatch(message)
	if len(matches) == 0 {
		return &discordgo.MessageSend{
			Content: "Sorry that ZIP code doesn't look right",
		}
	}
	zip := matches[len(matches)-1]

	// If ZIP not found, return an error
	if zip == "" {
		return &discordgo.MessageSend{
			Content: "Sorry that ZIP code doesn't look right",
		}
	}
	geoURL := fmt.Sprintf("%szip=%s,%s&appid=%s", GEOURL, zip, CountryCode, OpenWeatherToken)
	// Create new HTTP client & set timeout
	client := http.Client{Timeout: 5 * time.Second}

	// Query OpenWeather API
	response, err := client.Get(geoURL)
	if err != nil {
		return &discordgo.MessageSend{
			Content: "Sorry, there was an error trying to get the coordinates",
		}
	}
	geo, _ := io.ReadAll(response.Body)
	defer response.Body.Close()
	var geoData CoordinatesData
	json.Unmarshal(geo, &geoData)
	lat := geoData.Lat
	lon := geoData.Lon

	// Build full URL to query OpenWeather
	weatherURL := fmt.Sprintf("%slat=%f&lon=%f&appid=%s&units=metric", URL, lat, lon, OpenWeatherToken)

	// Query OpenWeather API
	response, err = client.Get(weatherURL)
	if err != nil {
		return &discordgo.MessageSend{
			Content: "Sorry, there was an error trying to get the weather",
		}
	}

	// Open HTTP response body
	body, _ := io.ReadAll(response.Body)
	defer response.Body.Close()

	// Convert JSON
	var data WeatherData
	json.Unmarshal([]byte(body), &data)

	// Pull out desired weather info & Convert to string if necessary
	city := geoData.Name
	conditions := data.Weather[0].Description
	temperature := strconv.FormatFloat(data.Main.Temp, 'f', 2, 64)
	feelsLike := strconv.FormatFloat(data.Main.FeelsLike, 'f', 2, 64)
	humidity := strconv.Itoa(data.Main.Humidity)
	wind := strconv.FormatFloat(data.Wind.Speed, 'f', 2, 64)

	// Build Discord embed response
	embed := &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{{
			Type:        discordgo.EmbedTypeRich,
			Title:       "Current Weather",
			Description: "Weather for " + city,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Conditions",
					Value:  conditions,
					Inline: true,
				},
				{
					Name:   "Temperature",
					Value:  temperature + "°C",
					Inline: true,
				},
				{
					Name:   "Feels Like",
					Value:  feelsLike + "°C",
					Inline: true,
				},
				{
					Name:   "Humidity",
					Value:  humidity + "%",
					Inline: true,
				},
				{
					Name:   "Wind",
					Value:  wind + " mph",
					Inline: true,
				},
			},
		},
		},
	}

	return embed
}

func processCountryCode(countryCode string) string {
	newCountryCode := CountryCode
	r, _ := regexp.Compile(`!country\s+(\w{2})$`)
	matches := r.FindStringSubmatch(countryCode)
	if len(matches) == 0 {
		return newCountryCode
	}
	code := matches[len(matches)-1]

	// If ZIP not found, return an error
	if code == "" {
		return newCountryCode
	}
	newCountryCode = code
	return newCountryCode
}
