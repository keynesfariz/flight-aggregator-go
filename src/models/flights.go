package models

type SearchRequest struct {
	Origin        string  `json:"origin" binding:"required"`
	Destination   string  `json:"destination" binding:"required"`
	DepartureDate string  `json:"departureDate" binding:"required"`
	ReturnDate    *string `json:"returnDate"`
	Passengers    int     `json:"passengers"`
	CabinClass    string  `json:"cabinClass"`
	Filters       Filters `json:"filters"`
	SortBy        string  `json:"sortBy"`
	SortOrder     string  `json:"sortOrder"`
}

type Filters struct {
	PriceMin           int      `json:"priceMin"`
	PriceMax           int      `json:"priceMax"`
	MaxStops           int      `json:"maxStops"`
	Airlines           []string `json:"airlines"`
	DepartureTimeRange string   `json:"departureTimeRange"`
	ArrivalTimeRange   string   `json:"arrivalTimeRange"`
	MaxDurationMinutes int      `json:"maxDurationMinutes"`
}

type Baggage struct {
	CarryOn *string `json:"carry_on"`
	Checked *string `json:"checked"`
}

type Flight struct {
	ID             string     `json:"id"`
	Provider       string     `json:"provider"`
	Airline        Airline    `json:"airline"`
	FlightNumber   string     `json:"flight_number"`
	Departure      EventPoint `json:"departure"`
	Arrival        EventPoint `json:"arrival"`
	Duration       Duration   `json:"duration"`
	Stops          int        `json:"stops"`
	Price          Price      `json:"price"`
	AvailableSeats int        `json:"available_seats"`
	CabinClass     string     `json:"cabin_class"`
	Aircraft       *string    `json:"aircraft"`
	Amenities      *[]string  `json:"amenities"`
	Baggage        Baggage    `json:"baggage"`
}

type Airline struct {
	Name string `json:"name"`
	Code string `json:"code"`
}

type EventPoint struct {
	Airport   string `json:"airport"`
	City      string `json:"city"`
	DateTime  string `json:"datetime"`
	Timestamp int64  `json:"timestamp"`
}

type Duration struct {
	TotalMinutes int    `json:"total_minutes"`
	Formatted    string `json:"formatted"`
}

type Price struct {
	Amount    int    `json:"amount"`
	Currency  string `json:"currency"`
	Formatted string `json:"formatted"`
}

type SearchResponse struct {
	Criteria SearchRequest `json:"search_criteria"`
	Metadata Metadata      `json:"metadata"`
	Flights  []Flight      `json:"flights"`
}

type Metadata struct {
	TotalResults     int  `json:"total_results"`
	ProvidersQueried int  `json:"providers_queried"`
	ProvidersSuccess int  `json:"providers_succeeded"`
	ProvidersFailed  int  `json:"providers_failed"`
	SearchTimeMs     int  `json:"search_time_ms"`
	CacheHit         bool `json:"cache_hit"`
}
