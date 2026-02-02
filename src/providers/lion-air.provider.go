package providers

import (
	"bookcabin-app-go/src/constants"
	"bookcabin-app-go/src/models"
	"bookcabin-app-go/src/utils"
	"encoding/json"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type LionAirRawFlight struct {
	ID      string `json:"id"`
	Carrier struct {
		Name string `json:"name"`
		Iata string `json:"iata"`
	} `json:"carrier"`
	Route struct {
		From struct {
			Code string `json:"code"`
			Name string `json:"name"`
			City string `json:"city"`
		} `json:"from"`
		To struct {
			Code string `json:"code"`
			Name string `json:"name"`
			City string `json:"city"`
		} `json:"to"`
	} `json:"route"`
	Schedule struct {
		Departure         string `json:"departure"`
		DepartureTimezone string `json:"departure_timezone"`
		Arrival           string `json:"arrival"`
		ArrivalTimezone   string `json:"arrival_timezone"`
	} `json:"schedule"`
	FlightTime int  `json:"flight_time"`
	IsDirect   bool `json:"is_direct"`
	Pricing    struct {
		Total    int    `json:"total"`
		Currency string `json:"currency"`
		FareType string `json:"fare_type"`
	} `json:"pricing"`
	SeatsLeft int    `json:"seats_left"`
	PlaneType string `json:"plane_type"`
	Services  struct {
		WifiAvailable    bool `json:"wifi_available"`
		MealsIncluded    bool `json:"meals_included"`
		BaggageAllowance struct {
			Cabin string `json:"cabin"`
			Hold  string `json:"hold"`
		} `json:"baggage_allowance"`
	} `json:"services"`
	StopCount int `json:"stop_count,omitempty"`
	Layovers  []struct {
		Airport         string `json:"airport"`
		DurationMinutes int    `json:"duration_minutes"`
	} `json:"layovers,omitempty"`
}

type LionAirRawFlightResponse struct {
	Success bool `json:"success"`
	Data    struct {
		AvailableFlights []LionAirRawFlight `json:"available_flights"`
	} `json:"data"`
}

type LionAirProvider struct {
	props SearchProviderProperty
}

func NewLionAirProvider() *LionAirProvider {
	return &LionAirProvider{props: SearchProviderProperty{
		Name:         "LionAir",
		SuccessRate:  100,
		ResponseTime: [2]int{50, 100},
		MockFile:     "garuda_indonesia_search_response.json",
	}}
}

func (pvd *LionAirProvider) Fetch(
	ctx *gin.Context,
	req models.SearchRequest,
) ([]models.Flight, error) {
	data, err := SimulateFetchWithWait(ctx.Request.Context(), pvd.props)
	if err != nil {
		return nil, err
	}

	var rawFlightResponse LionAirRawFlightResponse
	if err := json.Unmarshal(data, &rawFlightResponse); err != nil {
		return nil, err
	}

	results := make([]models.Flight, 0)

	for _, flight := range rawFlightResponse.Data.AvailableFlights {

		departLoc, _ := time.LoadLocation(flight.Schedule.DepartureTimezone)
		arriveLoc, _ := time.LoadLocation(flight.Schedule.ArrivalTimezone)

		flightDepartureDate, errDept := time.ParseInLocation(constants.LA_DateTimeLayout, flight.Schedule.Departure, departLoc)
		flightArrivalDate, errArrv := time.ParseInLocation(constants.LA_DateTimeLayout, flight.Schedule.Arrival, arriveLoc)

		if errDept != nil ||
			errArrv != nil ||
			flight.Route.From.Code != req.Origin ||
			flight.Route.To.Code != req.Destination ||
			flightDepartureDate.Format(time.DateOnly) != req.DepartureDate ||
			flight.SeatsLeft < req.Passengers {
			continue
		}

		amenities, baggage := buildLionAirAmenities(flight)

		results = append(results, models.Flight{
			ID:       utils.GetFlightId(flight.Carrier.Name, flight.ID),
			Provider: pvd.props.Name,
			Airline: models.Airline{
				Name: flight.Carrier.Name,
				Code: flight.Carrier.Iata,
			},
			FlightNumber: flight.ID,
			Departure: models.EventPoint{
				Airport:   flight.Route.From.Code,
				City:      flight.Route.From.City,
				DateTime:  utils.FormatDateTime(flightDepartureDate, flight.Route.From.Code),
				Timestamp: flightDepartureDate.Unix(),
			},
			Arrival: models.EventPoint{
				Airport:   flight.Route.To.Code,
				City:      flight.Route.To.City,
				DateTime:  utils.FormatDateTime(flightArrivalDate, flight.Route.To.Code),
				Timestamp: flightArrivalDate.Unix(),
			},
			Duration: models.Duration{
				TotalMinutes: flight.FlightTime,
				Formatted:    utils.FormatDurationToHumans(flight.FlightTime),
			},
			Stops: flight.StopCount,
			Price: models.Price{
				Amount:    flight.Pricing.Total,
				Currency:  flight.Pricing.Currency,
				Formatted: utils.FormatPrice(flight.Pricing.Total, flight.Pricing.Currency),
			},
			AvailableSeats: flight.SeatsLeft,
			CabinClass:     strings.ToLower(flight.Pricing.FareType),
			Aircraft:       &flight.PlaneType,
			Amenities:      &amenities,
			Baggage:        baggage,
		})
	}

	return results, err
}

func buildLionAirAmenities(flight LionAirRawFlight) ([]string, models.Baggage) {
	var amenities []string
	if flight.Services.MealsIncluded {
		amenities = append(amenities, "Meals included")
	}
	if flight.Services.WifiAvailable {
		amenities = append(amenities, "Wifi available")
	}

	return amenities, models.Baggage{
		CarryOn: &flight.Services.BaggageAllowance.Cabin,
		Checked: &flight.Services.BaggageAllowance.Hold,
	}
}
