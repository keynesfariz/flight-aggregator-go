package providers

import (
	"bookcabin-app-go/src/constants"
	"bookcabin-app-go/src/models"
	"bookcabin-app-go/src/utils"
	"encoding/json"
	"time"

	"github.com/gin-gonic/gin"
)

type BatikAirRawFlight struct {
	FlightNumber      string `json:"flightNumber"`
	AirlineName       string `json:"airlineName"`
	AirlineIATA       string `json:"airlineIATA"`
	Origin            string `json:"origin"`
	Destination       string `json:"destination"`
	DepartureDateTime string `json:"departureDateTime"`
	ArrivalDateTime   string `json:"arrivalDateTime"`
	TravelTime        string `json:"travelTime"`
	NumberOfStops     int    `json:"numberOfStops"`
	Fare              struct {
		BasePrice    int    `json:"basePrice"`
		Taxes        int    `json:"taxes"`
		TotalPrice   int    `json:"totalPrice"`
		CurrencyCode string `json:"currencyCode"`
		Class        string `json:"class"`
	} `json:"fare"`
	SeatsAvailable  int      `json:"seatsAvailable"`
	AircraftModel   string   `json:"aircraftModel"`
	BaggageInfo     string   `json:"baggageInfo"`
	OnboardServices []string `json:"onboardServices"`
	Connections     []struct {
		StopAirport  string `json:"stopAirport"`
		StopDuration string `json:"stopDuration"`
	} `json:"connections,omitempty"`
}

type BatikAirRawFlightResponse struct {
	Code    int                 `json:"code"`
	Message string              `json:"message"`
	Results []BatikAirRawFlight `json:"results"`
}

type BatikAirProvider struct {
	props SearchProviderProperty
}

func NewBatikAirProvider() *BatikAirProvider {
	return &BatikAirProvider{props: SearchProviderProperty{
		Name:         "BatikAir",
		SuccessRate:  100,
		ResponseTime: [2]int{200, 400},
		MockFile:     "batik_air_search_response.json",
	}}
}

func (pvd *BatikAirProvider) Fetch(
	ctx *gin.Context,
	req models.SearchRequest,
) ([]models.Flight, error) {
	data, err := SimulateFetchWithWait(ctx.Request.Context(), pvd.props)
	if err != nil {
		return nil, err
	}

	var rawFlightResponse BatikAirRawFlightResponse
	if err := json.Unmarshal(data, &rawFlightResponse); err != nil {
		return nil, err
	}

	results := make([]models.Flight, 0)

	for _, flight := range rawFlightResponse.Results {

		flightDepartureDate, errDept := time.Parse(constants.BA_DateTimeLayout, flight.DepartureDateTime)
		flightArrivalDate, errArrv := time.Parse(constants.BA_DateTimeLayout, flight.ArrivalDateTime)

		if errDept != nil ||
			errArrv != nil ||
			flight.Origin != req.Origin ||
			flight.Destination != req.Destination ||
			flightDepartureDate.Format(time.DateOnly) != req.DepartureDate ||
			flight.SeatsAvailable < req.Passengers {
			continue
		}

		results = append(results, models.Flight{
			ID:       utils.GetFlightId(flight.AirlineName, flight.FlightNumber),
			Provider: pvd.props.Name,
			Airline: models.Airline{
				Name: flight.AirlineName,
				Code: flight.AirlineIATA,
			},
			FlightNumber: flight.FlightNumber,
			Departure: models.EventPoint{
				Airport:   flight.Origin,
				City:      getCity(flight.Origin),
				DateTime:  utils.FormatDateTime(flightDepartureDate, flight.Origin),
				Timestamp: flightDepartureDate.Unix(),
			},
			Arrival: models.EventPoint{
				Airport:   flight.Destination,
				City:      getCity(flight.Destination),
				DateTime:  utils.FormatDateTime(flightArrivalDate, flight.Destination),
				Timestamp: flightArrivalDate.Unix(),
			},
			Duration: models.Duration{
				TotalMinutes: utils.FormatDurationToMinutes(flight.TravelTime),
				Formatted:    flight.TravelTime,
			},
			Stops: flight.NumberOfStops,
			Price: models.Price{
				Amount:    flight.Fare.TotalPrice,
				Currency:  flight.Fare.CurrencyCode,
				Formatted: utils.FormatPrice(flight.Fare.TotalPrice, flight.Fare.CurrencyCode),
			},
			AvailableSeats: flight.SeatsAvailable,
			CabinClass:     "economy",
			Aircraft:       &flight.AircraftModel,
			Amenities:      &flight.OnboardServices,
			Baggage:        parseBaggageInfo(flight.BaggageInfo),
		})
	}

	return results, err
}
