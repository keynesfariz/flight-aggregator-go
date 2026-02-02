package utils

import (
	"bookcabin-app-go/src/constants"
	"bookcabin-app-go/src/models"
	"slices"
	"time"
)

func ApplySearchFilters(flights *[]models.Flight, req models.SearchRequest) {
	var filteredFlights []models.Flight

	for _, flight := range *flights {

		if req.Filters.PriceMin > 0 && flight.Price.Amount < req.Filters.PriceMin {
			continue
		}

		if req.Filters.PriceMax > 0 && flight.Price.Amount > req.Filters.PriceMax {
			continue
		}

		if req.Filters.MaxStops > 0 && flight.Stops > req.Filters.MaxStops {
			continue
		}

		isAirlinePresent := slices.Contains(req.Filters.Airlines, flight.Airline.Name)
		if len(req.Filters.Airlines) > 0 && !isAirlinePresent {
			continue
		}

		if req.Filters.MaxDurationMinutes > 0 && flight.Duration.TotalMinutes > req.Filters.MaxDurationMinutes {
			continue
		}

		if req.Filters.DepartureTimeRange != "" {
			parsedFlightDeparture, errFlg := time.Parse(constants.GA_DateTimeLayout, req.Filters.DepartureTimeRange)
			parsedDepartureTimeRange, errFlt := time.Parse(time.RFC3339, req.Filters.DepartureTimeRange)

			if errFlg == nil && errFlt == nil {
				if parsedDepartureTimeRange.After(parsedFlightDeparture) {
					continue
				}
			}
		}

		if req.Filters.ArrivalTimeRange != "" {
			parsedFlightArrival, errFlg := time.Parse(constants.GA_DateTimeLayout, req.Filters.ArrivalTimeRange)
			parsedArrivalTimeRange, errFlt := time.Parse(time.RFC3339, req.Filters.ArrivalTimeRange)

			if errFlg == nil && errFlt == nil {
				if parsedFlightArrival.After(parsedArrivalTimeRange) {
					continue
				}
			}
		}

		filteredFlights = append(filteredFlights, flight)
	}

	*flights = filteredFlights
}
