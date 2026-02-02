
# Flight Search & Aggregation System in Go

This system aggregates flight data from multiple mocked airline APIs, processes and filters the results, and returns optimized search responses to users. It is part of the BookCabin take-home assignment and is built using Go with the Gin framework and Redis for API caching.  

## Tech Stack

- Go + Gin
- Redis
  
## Setup

**1. Install dependencies**

```
go mod tidy
```

**2. Run on local system**
```
go run main.go
```  

## Making an API Request

After running the application, make a **POST** request to  `/search/` with `application/json` on the app, usually it's served on `localhost:8080`

**Endpoint**
```
POST /search
```

**Required Request Body** 
```
{
  "origin": "CGK",
  "destination": "DPS",
  "departureDate": "2025-12-15"
}
```

**Full Request Body**
```
{
  "origin": "CGK",
  "destination": "DPS",
  "departureDate": "2025-12-15",
  "returnDate": null,
  "passengers": 1,
  "cabinClass": "economy",
  "sortBy": "best_value",
  "sortOrder": "asc",
  "filters": {
    "priceMin": 0,
    "priceMax": 0,
    "maxStops": 0,
    "airlines": [],
    "departureTimeRange": "",
    "arrivalTimeRange": "",
    "maxDurationMinutes": 0
  }
}
```

## Under the Hood

- Simulates multiple airline providers with each provider has its own **configurable real-world conditions** and **retry logic** with exponential backoff set to 8 ms.
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
- **Best Value** sorting is implemented using Weighted Scoring Model
- Timezone conversions to Go's `2006-01-02T15:04:05-07:00` format
- Support for displaying currency in IDR formatting with thousands separator
