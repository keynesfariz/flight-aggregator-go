# Flight Search & Aggregation System in Go

This system aggregates flight data from multiple mocked airline APIs, processes and filters the results, and returns optimized search responses to users. It is part of the BookCabin take-home assignment and is built using Go with the Gin framework and Redis for API caching.

## Tech Stack

- Go + Gin
- Redis

## Installation

```

### Run on local system

## Under the Hood

- Simulates multiple airline providers with each provider has its own configurable real-world conditions, example:

```

func NewGarudaIndonesiaProvider() \*GarudaIndonesiaProvider {
return &GarudaIndonesiaProvider{props: SearchProviderProperty{
Name: "GarudaIndonesia",
SuccessRate: 100, // percent
ResponseTime: [2]int{50, 100},
MockFile: "garuda_indonesia_search_response.json",
}}
}

```

- Parallel fetch execution on all airline providers using `sync.WaitGroup`
- Caches set on both:
  - Mocked Provider's result
  - User's search result, with **idempotency** implemented on search criteria
- Filterable by:
  - price range `priceMin` and `priceMax`
  - max number of stops `maxStops`
  - travel time range `departureTimeRange` and `arrivalTimeRange`
  - airlines `airlines`, and
  - travel duration `maxDurationMinutes`
- Sortable by:
  - Best value `best_value`
  - `price` in `asc` and `desc` order
  - `duration` in `asc` and `desc` order
  - `departure` in `asc` and `desc` order
  - `arrival` in `asc` and `desc` order
- Scoring implemented for best value sorting
- Timezone differe
```
