package utils

import (
	"bookcabin-app-go/src/models"
	"slices"
	"strings"
)

/* Flight Best Value algorithm using simple Weighted Scoring Model */
const (
	price_duration_stop     = .8
	amenitis_checkedbaggage = .2
)

type SearchResultNormalizer struct {
	MinPrice    int
	MaxPrice    int
	MinDuration int
	MaxDuration int
	MinStop     int
	MaxStop     int
	MinAmenity  int
	MaxAmenity  int
}

// Normalize value into float64 ranged from 0.00 to 1.00
func norm(value int, min int, max int) float64 {
	if min == max {
		return 0
	}
	return float64(value-min) / float64(max-min)
}

func GetBestValueScores(flights []models.Flight) map[string]float64 {
	var prices, durations, stops, amenities []int
	results := make(map[string]float64)

	if len(flights) == 0 {
		return results
	}

	for _, flight := range flights {
		prices = append(prices, flight.Price.Amount)
		durations = append(durations, flight.Duration.TotalMinutes)
		stops = append(stops, flight.Stops)

		if flight.Amenities == nil {
			amenities = append(amenities, 0)
		}

		if flight.Amenities != nil {
			amenity := len(*flight.Amenities)
			amenities = append(amenities, amenity)
		}
	}

	normalizer := SearchResultNormalizer{
		MinPrice:    slices.Min(prices),
		MaxPrice:    slices.Max(prices),
		MinDuration: slices.Min(durations),
		MaxDuration: slices.Max(durations),
		MinStop:     slices.Min(stops),
		MaxStop:     slices.Max(stops),
		MinAmenity:  slices.Min(amenities),
		MaxAmenity:  slices.Max(amenities),
	}

	for _, flight := range flights {
		results[flight.ID] = CalculateFlightScore(flight, normalizer)
	}

	return results
}

func CalculateFlightScore(flight models.Flight,
	normalizer SearchResultNormalizer,
) float64 {
	// Positive values: amenities, free checked baggage
	// Negative values: price, duration, stop
	var positive, negative float64

	// add amenities point
	if flight.Amenities != nil {
		positive += norm(len(*flight.Amenities), normalizer.MinAmenity, normalizer.MaxAmenity)
	}

	freeCheckedBaggagePoint := 0
	if flight.Baggage.Checked != nil {
		checkedBaggageInfo := flight.Baggage.Checked
		if !strings.Contains(strings.ToLower(*checkedBaggageInfo), "additional fee") {
			freeCheckedBaggagePoint = 1
		}
	}
	// add free checked baggage point
	positive += norm(freeCheckedBaggagePoint, 0, 1)

	// calc. price
	negative += norm(flight.Price.Amount, normalizer.MinPrice, normalizer.MaxPrice)

	// calc. duration
	negative += norm(flight.Duration.TotalMinutes, normalizer.MinDuration, normalizer.MaxDuration)

	// calc. # of stops
	negative += norm(flight.Stops, normalizer.MinStop, normalizer.MaxStop)

	// apply weight
	return positive*amenitis_checkedbaggage + (1-negative)*price_duration_stop
}
