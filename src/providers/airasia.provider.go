package providers

import (
	"bookcabin-app-go/src/models"
	"bookcabin-app-go/src/utils"
	"encoding/json"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samber/lo"
)

type AirAsiaRawFlight struct {
	FlightCode    string    `json:"flight_code"`
	Airline       string    `json:"airline"`
	FromAirport   string    `json:"from_airport"`
	ToAirport     string    `json:"to_airport"`
	DepartTime    time.Time `json:"depart_time"`
	ArriveTime    time.Time `json:"arrive_time"`
	DurationHours float64   `json:"duration_hours"`
	DirectFlight  bool      `json:"direct_flight"`
	PriceIdr      int       `json:"price_idr"`
	Seats         int       `json:"seats"`
	CabinClass    string    `json:"cabin_class"`
	BaggageNote   string    `json:"baggage_note"`
	Stops         []struct {
		Airport         string `json:"airport"`
		WaitTimeMinutes int    `json:"wait_time_minutes"`
	} `json:"stops,omitempty"`
}

type AirAsiaRawFlightResponse struct {
	Status  string             `json:"status"`
	Flights []AirAsiaRawFlight `json:"flights"`
}

type AirAsiaProvider struct {
	props SearchProviderProperty
}

func NewAirAsiaProvider() *AirAsiaProvider {
	return &AirAsiaProvider{props: SearchProviderProperty{
		Name:         "AirAsia",
		SuccessRate:  90,
		ResponseTime: [2]int{50, 150},
		MockFile:     "airasia_search_response.json",
	}}
}

func (pvd *AirAsiaProvider) Fetch(
	ctx *gin.Context,
	req models.SearchRequest,
) ([]models.Flight, error) {
	data, err := SimulateFetchWithWait(ctx.Request.Context(), pvd.props)
	if err != nil {
		return nil, err
	}

	var rawFlightResponse AirAsiaRawFlightResponse
	if err := json.Unmarshal(data, &rawFlightResponse); err != nil {
		return nil, err
	}

	results := make([]models.Flight, 0)

	for _, flight := range rawFlightResponse.Flights {

		if flight.FromAirport != req.Origin ||
			flight.ToAirport != req.Destination ||
			flight.DepartTime.Format(time.DateOnly) != req.DepartureDate ||
			flight.Seats < req.Passengers {
			continue
		}

		durationInt := int(flight.DurationHours * 60)

		results = append(results, models.Flight{
			ID:       utils.GetFlightId(flight.Airline, flight.FlightCode),
			Provider: pvd.props.Name,
			Airline: models.Airline{
				Name: flight.Airline,
				Code: getAirAsiaAirlineCode(flight.FlightCode),
			},
			FlightNumber: flight.FlightCode,
			Departure: models.EventPoint{
				Airport:   flight.FromAirport,
				City:      getCity(flight.FromAirport),
				DateTime:  utils.FormatDateTime(flight.DepartTime, flight.FromAirport),
				Timestamp: flight.DepartTime.Unix(),
			},
			Arrival: models.EventPoint{
				Airport:   flight.ToAirport,
				City:      getCity(flight.ToAirport),
				DateTime:  utils.FormatDateTime(flight.ArriveTime, flight.ToAirport),
				Timestamp: flight.ArriveTime.Unix(),
			},
			Duration: models.Duration{
				TotalMinutes: durationInt,
				Formatted:    utils.FormatDurationToHumans(durationInt),
			},
			Stops: len(flight.Stops),
			Price: models.Price{
				Amount:    flight.PriceIdr,
				Currency:  "IDR",
				Formatted: utils.FormatPrice(flight.PriceIdr, "IDR"),
			},
			AvailableSeats: flight.Seats,
			CabinClass:     strings.ToLower(flight.CabinClass),
			Aircraft:       nil,
			Amenities:      nil,
			Baggage:        parseBaggageInfo(flight.BaggageNote),
		})
	}

	return results, err
}

func getAirAsiaAirlineCode(flightCode string) string {
	return lo.Substring(flightCode, 0, 2)
}
