
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
**Response**
```
{
  "search_criteria": {
    "origin": "CGK",
    "destination": "DPS",
    "departureDate": "2025-12-15",
    "returnDate": null,
    "passengers": 1,
    "cabinClass": "economy",
    "filters": {
      "priceMin": 0,
      "priceMax": 0,
      "maxStops": 0,
      "airlines": null,
      "departureTimeRange": "",
      "arrivalTimeRange": "",
      "maxDurationMinutes": 0
    },
    "sortBy": "best_value",
    "sortOrder": "asc"
  },
  "metadata": {
    "total_results": 9,
    "providers_queried": 4,
    "providers_succeeded": 4,
    "providers_failed": 0,
    "search_time_ms": 6711,
    "cache_hit": false
  },
  "flights": [
	// ... Example of one result
    {
      "id": "QZ7250_AirAsia",
      "provider": "AirAsia",
      "airline": {
        "name": "AirAsia",
        "code": "QZ"
      },
      "flight_number": "QZ7250",
      "departure": {
        "airport": "CGK",
        "city": "Jakarta",
        "datetime": "2025-12-15T15:15:00+07:00",
        "timestamp": 1765786500
      },
      "arrival": {
        "airport": "DPS",
        "city": "Denpasar",
        "datetime": "2025-12-15T20:35:00+08:00",
        "timestamp": 1765802100
      },
      "duration": {
        "total_minutes": 259,
        "formatted": "04h 19m"
      },
      "stops": 1,
      "price": {
        "amount": 485000,
        "currency": "IDR",
        "formatted": "IDR 485.000"
      },
      "available_seats": 88,
      "cabin_class": "economy",
      "aircraft": null,
      "amenities": null,
      "baggage": {
        "carry_on": "Cabin Baggage Only",
        "checked": " Checked Bags Additional Fee"
      }
    }
  ]
}
```

## Under the Hood

- Simulates multiple airline providers with each provider has its own **configurable real-world conditions** and **retry logic** with exponential backoff set to 8 ms.
- Parallel fetch execution on all airline providers using `sync.WaitGroup`
- Caches set on both:
	- Mocked Provider's result
	- User's search result, with **idempotency** implemented on search criteria, marked with `"cache_hit": true` on response' metadata
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
