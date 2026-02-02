package providers

import (
	"bookcabin-app-go/src/constants"
	"bookcabin-app-go/src/models"
	"bookcabin-app-go/src/utils"
	"encoding/json"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type GarudaIndonesiaRawFlight struct {
	FlightID    string `json:"flight_id"`
	Airline     string `json:"airline"`
	AirlineCode string `json:"airline_code"`
	Departure   struct {
		Airport  string `json:"airport"`
		City     string `json:"city"`
		Time     string `json:"time"`
		Terminal string `json:"terminal"`
	} `json:"departure"`
	Arrival struct {
		Airport  string `json:"airport"`
		City     string `json:"city"`
		Time     string `json:"time"`
		Terminal string `json:"terminal"`
	} `json:"arrival"`
	DurationMinutes int    `json:"duration_minutes"`
	Stops           int    `json:"stops"`
	Aircraft        string `json:"aircraft"`
	Price           struct {
		Amount   int    `json:"amount"`
		Currency string `json:"currency"`
	} `json:"price"`
	AvailableSeats int    `json:"available_seats"`
	FareClass      string `json:"fare_class"`
	Baggage        struct {
		CarryOn int `json:"carry_on"`
		Checked int `json:"checked"`
	} `json:"baggage"`
	Amenities []string `json:"amenities,omitempty"`
	Segments  []struct {
		FlightNumber string `json:"flight_number"`
		Departure    struct {
			Airport string    `json:"airport"`
			Time    time.Time `json:"time"`
		} `json:"departure"`
		Arrival struct {
			Airport string    `json:"airport"`
			Time    time.Time `json:"time"`
		} `json:"arrival"`
		DurationMinutes int `json:"duration_minutes"`
		LayoverMinutes  int `json:"layover_minutes,omitempty"`
	} `json:"segments,omitempty"`
}

type GarudaIndonesiaRawFlightResponse struct {
	Status  string                     `json:"status"`
	Flights []GarudaIndonesiaRawFlight `json:"flights"`
}

type GarudaIndonesiaProvider struct {
	props SearchProviderProperty
}

func NewGarudaIndonesiaProvider() *GarudaIndonesiaProvider {
	return &GarudaIndonesiaProvider{props: SearchProviderProperty{
		Name:         "GarudaIndonesia",
		SuccessRate:  100,
		ResponseTime: [2]int{50, 100},
		MockFile:     "garuda_indonesia_search_response.json",
	}}
}

func (pvd *GarudaIndonesiaProvider) Fetch(
	ctx *gin.Context,
	req models.SearchRequest,
) ([]models.Flight, error) {
	data, err := SimulateFetchWithWait(ctx.Request.Context(), pvd.props)
	if err != nil {
		return nil, err
	}

	var rawFlightResponse GarudaIndonesiaRawFlightResponse
	if err := json.Unmarshal(data, &rawFlightResponse); err != nil {
		return nil, err
	}

	results := make([]models.Flight, 0)

	for _, flight := range rawFlightResponse.Flights {

		flightDepartureDate, errDept := time.Parse(constants.GA_DateTimeLayout, flight.Departure.Time)
		flightArrivalDate, errArrv := time.Parse(constants.GA_DateTimeLayout, flight.Arrival.Time)

		if errDept != nil ||
			errArrv != nil ||
			flight.Departure.Airport != req.Origin ||
			flight.Arrival.Airport != req.Destination ||
			flightDepartureDate.Format(time.DateOnly) != req.DepartureDate ||
			flight.AvailableSeats < req.Passengers {
			continue
		}

		carryOn := strconv.Itoa(flight.Baggage.CarryOn)
		checked := strconv.Itoa(flight.Baggage.Checked)

		results = append(results, models.Flight{
			ID:       utils.GetFlightId(flight.Airline, flight.FlightID),
			Provider: pvd.props.Name,
			Airline: models.Airline{
				Name: flight.Airline,
				Code: flight.AirlineCode,
			},
			FlightNumber: flight.FlightID,
			Departure: models.EventPoint{
				Airport:   flight.Departure.Airport,
				City:      flight.Departure.City,
				DateTime:  utils.FormatDateTime(flightDepartureDate, flight.Departure.Airport),
				Timestamp: flightDepartureDate.Unix(),
			},
			Arrival: models.EventPoint{
				Airport:   flight.Arrival.Airport,
				City:      flight.Arrival.City,
				DateTime:  utils.FormatDateTime(flightArrivalDate, flight.Arrival.Airport),
				Timestamp: flightArrivalDate.Unix(),
			},
			Duration: models.Duration{
				TotalMinutes: flight.DurationMinutes,
				Formatted:    utils.FormatDurationToHumans(flight.DurationMinutes),
			},
			Stops: flight.Stops,
			Price: models.Price{
				Amount:    flight.Price.Amount,
				Currency:  flight.Price.Currency,
				Formatted: utils.FormatPrice(flight.Price.Amount, flight.Price.Currency),
			},
			AvailableSeats: flight.AvailableSeats,
			CabinClass:     strings.ToLower(flight.FareClass),
			Aircraft:       &flight.Aircraft,
			Amenities:      &flight.Amenities,
			Baggage: models.Baggage{
				CarryOn: &carryOn,
				Checked: &checked,
			},
		})
	}

	return results, err
}
