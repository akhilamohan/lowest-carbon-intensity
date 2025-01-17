package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"
)

type CarbonIntensityForcast struct {
	From      string `json:"from"`
	To        string `json:"to"`
	Intensity struct {
		Forecast int `json:"forecast"`
	} `json:"intensity"`
}

type CarbonIntensitySlot struct {
	ValidFrom time.Time `json:"valid_from"`
	ValidTo   time.Time `json:"valid_to"`
	Carbon    struct {
		Intensity int `json:"intensity"`
	} `json:"carbon"`
}

// ForecastResponse represents the response structure from the forecast API
type ForecastResponse struct {
	Data []CarbonIntensityForcast `json:"data"`
}

// FetchForecast retrieves carbon intensity forecast for the next 24 hours
func FetchForecast() ([]CarbonIntensityForcast, error) {
	currentTime := time.Now()
	timeString := currentTime.Format("2006-01-02T15:04Z")

	url := fmt.Sprintf("https://api.carbonintensity.org.uk/intensity/%s/fw24h", timeString)
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	var forecast ForecastResponse
	if err := json.NewDecoder(resp.Body).Decode(&forecast); err != nil {
		log.Println("decode error", err)
		return nil, err
	}

	return forecast.Data, nil
}

// Finds the time slots with the lowest carbon intensity
func findLowestIntensitySlots(carbonIntensityList []CarbonIntensityForcast, duration int,
	contiguous bool) []CarbonIntensitySlot {

	sort.Slice(carbonIntensityList, func(i, j int) bool {
		return carbonIntensityList[i].Intensity.Forecast < carbonIntensityList[j].Intensity.Forecast
	})

	var results []CarbonIntensitySlot

	if contiguous {
		for _, ci := range carbonIntensityList {
			toTime, err := time.Parse("2006-01-02T15:04Z", ci.To)
			if err != nil {
				log.Println(err)
				continue
			}
			fromTime, err := time.Parse("2006-01-02T15:04Z", ci.From)
			if err != nil {
				log.Println(err)
				continue
			}
			timeDiff := int(toTime.Sub(fromTime).Minutes())
			if timeDiff >= duration {
				var slot CarbonIntensitySlot
				slot.ValidFrom = fromTime
				slot.ValidTo = fromTime.Add(time.Duration(duration) * time.Minute)
				slot.Carbon.Intensity = ci.Intensity.Forecast
				results = append(results, slot)
				break
			}
		}
	} else {
		curr_duration := 0
		for _, ci := range carbonIntensityList {
			toTime, err := time.Parse("2006-01-02T15:04Z", ci.To)
			if err != nil {
				log.Println(err)
				continue
			}
			fromTime, err := time.Parse("2006-01-02T15:04Z", ci.From)
			if err != nil {
				log.Println(err)
				continue
			}
			timeDiff := int(toTime.Sub(fromTime).Minutes())
			if curr_duration < duration {
				var slot CarbonIntensitySlot
				slot.ValidFrom = fromTime
				slot.Carbon.Intensity = ci.Intensity.Forecast

				tmp_duration := curr_duration + timeDiff
				if tmp_duration > duration {
					slot.ValidTo = fromTime.Add(time.Duration(duration-curr_duration) * time.Minute)
					curr_duration = duration
				} else {
					slot.ValidTo = toTime
					curr_duration = tmp_duration
				}
				results = append(results, slot)
			}
			if curr_duration == duration {
				break
			}
		}
	}
	log.Println("results   ", results)
	return results
}

// Handles the API endpoint
func handleSlots(w http.ResponseWriter, r *http.Request) {
	duration := 30
	contiguous := false

	if durationParam := r.URL.Query().Get("duration"); durationParam != "" {
		if d, err := strconv.Atoi(durationParam); err == nil {
			duration = d
		}
	}

	if contiguousParam := r.URL.Query().Get("contiguous"); contiguousParam != "" {
		if c, err := strconv.ParseBool(contiguousParam); err == nil {
			contiguous = c
		}
	}

	// Fetch the forecast data
	forecast, err := FetchForecast()
	if err != nil {
		http.Error(w, "Failed to fetch forecast data", http.StatusInternalServerError)
		return
	}

	// Calculate the lowest intensity slots
	slots := findLowestIntensitySlots(forecast, duration, contiguous)

	// Prepare the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(slots)
}

func main() {
	http.HandleFunc("/slots", handleSlots)
	port := "3000"
	log.Printf("Starting server on port %s...\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Println("Failed to start server:", err)
		os.Exit(1)
	}
}
