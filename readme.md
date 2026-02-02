
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
    },
    {
      "id": "ID7042_BatikAir",
      "provider": "BatikAir",
      "airline": {
        "name": "Batik Air",
        "code": "ID"
      },
      "flight_number": "ID7042",
      "departure": {
        "airport": "CGK",
        "city": "Jakarta",
        "datetime": "2025-12-15T18:45:00+07:00",
        "timestamp": 1765799100
      },
      "arrival": {
        "airport": "DPS",
        "city": "Denpasar",
        "datetime": "2025-12-15T23:50:00+08:00",
        "timestamp": 1765813800
      },
      "duration": {
        "total_minutes": 185,
        "formatted": "3h 5m"
      },
      "stops": 1,
      "price": {
        "amount": 950000,
        "currency": "IDR",
        "formatted": "IDR 950.000"
      },
      "available_seats": 41,
      "cabin_class": "economy",
      "aircraft": "Airbus A320",
      "amenities": [
        "Snack"
      ],
      "baggage": {
        "carry_on": "7Kg Cabin",
        "checked": " 20Kg Checked"
      }
    },
    {
      "id": "GA410_GarudaIndonesia",
      "provider": "GarudaIndonesia",
      "airline": {
        "name": "Garuda Indonesia",
        "code": "GA"
      },
      "flight_number": "GA410",
      "departure": {
        "airport": "CGK",
        "city": "Jakarta",
        "datetime": "2025-12-15T09:30:00+07:00",
        "timestamp": 1765765800
      },
      "arrival": {
        "airport": "DPS",
        "city": "Denpasar",
        "datetime": "2025-12-15T12:25:00+08:00",
        "timestamp": 1765772700
      },
      "duration": {
        "total_minutes": 115,
        "formatted": "01h 55m"
      },
      "stops": 0,
      "price": {
        "amount": 1450000,
        "currency": "IDR",
        "formatted": "IDR 1.450.000"
      },
      "available_seats": 15,
      "cabin_class": "economy",
      "aircraft": "Airbus A330-300",
      "amenities": [
        "wifi",
        "power_outlet",
        "meal",
        "entertainment"
      ],
      "baggage": {
        "carry_on": "1",
        "checked": "2"
      }
    },
    {
      "id": "GA400_GarudaIndonesia",
      "provider": "GarudaIndonesia",
      "airline": {
        "name": "Garuda Indonesia",
        "code": "GA"
      },
      "flight_number": "GA400",
      "departure": {
        "airport": "CGK",
        "city": "Jakarta",
        "datetime": "2025-12-15T06:00:00+07:00",
        "timestamp": 1765753200
      },
      "arrival": {
        "airport": "DPS",
        "city": "Denpasar",
        "datetime": "2025-12-15T08:50:00+08:00",
        "timestamp": 1765759800
      },
      "duration": {
        "total_minutes": 110,
        "formatted": "01h 50m"
      },
      "stops": 0,
      "price": {
        "amount": 1250000,
        "currency": "IDR",
        "formatted": "IDR 1.250.000"
      },
      "available_seats": 28,
      "cabin_class": "economy",
      "aircraft": "Boeing 737-800",
      "amenities": [
        "wifi",
        "meal",
        "entertainment"
      ],
      "baggage": {
        "carry_on": "1",
        "checked": "2"
      }
    },
    {
      "id": "ID6520_BatikAir",
      "provider": "BatikAir",
      "airline": {
        "name": "Batik Air",
        "code": "ID"
      },
      "flight_number": "ID6520",
      "departure": {
        "airport": "CGK",
        "city": "Jakarta",
        "datetime": "2025-12-15T13:30:00+07:00",
        "timestamp": 1765780200
      },
      "arrival": {
        "airport": "DPS",
        "city": "Denpasar",
        "datetime": "2025-12-15T16:20:00+08:00",
        "timestamp": 1765786800
      },
      "duration": {
        "total_minutes": 110,
        "formatted": "1h 50m"
      },
      "stops": 0,
      "price": {
        "amount": 1180000,
        "currency": "IDR",
        "formatted": "IDR 1.180.000"
      },
      "available_seats": 18,
      "cabin_class": "economy",
      "aircraft": "Boeing 737-800",
      "amenities": [
        "Meal",
        "Beverage",
        "Entertainment"
      ],
      "baggage": {
        "carry_on": "7Kg Cabin",
        "checked": " 20Kg Checked"
      }
    },
    {
      "id": "ID6514_BatikAir",
      "provider": "BatikAir",
      "airline": {
        "name": "Batik Air",
        "code": "ID"
      },
      "flight_number": "ID6514",
      "departure": {
        "airport": "CGK",
        "city": "Jakarta",
        "datetime": "2025-12-15T07:15:00+07:00",
        "timestamp": 1765757700
      },
      "arrival": {
        "airport": "DPS",
        "city": "Denpasar",
        "datetime": "2025-12-15T10:00:00+08:00",
        "timestamp": 1765764000
      },
      "duration": {
        "total_minutes": 105,
        "formatted": "1h 45m"
      },
      "stops": 0,
      "price": {
        "amount": 1100000,
        "currency": "IDR",
        "formatted": "IDR 1.100.000"
      },
      "available_seats": 32,
      "cabin_class": "economy",
      "aircraft": "Airbus A320",
      "amenities": [
        "Snack",
        "Beverage"
      ],
      "baggage": {
        "carry_on": "7Kg Cabin",
        "checked": " 20Kg Checked"
      }
    },
    {
      "id": "QZ524_AirAsia",
      "provider": "AirAsia",
      "airline": {
        "name": "AirAsia",
        "code": "QZ"
      },
      "flight_number": "QZ524",
      "departure": {
        "airport": "CGK",
        "city": "Jakarta",
        "datetime": "2025-12-15T10:00:00+07:00",
        "timestamp": 1765767600
      },
      "arrival": {
        "airport": "DPS",
        "city": "Denpasar",
        "datetime": "2025-12-15T12:45:00+08:00",
        "timestamp": 1765773900
      },
      "duration": {
        "total_minutes": 105,
        "formatted": "01h 45m"
      },
      "stops": 0,
      "price": {
        "amount": 720000,
        "currency": "IDR",
        "formatted": "IDR 720.000"
      },
      "available_seats": 54,
      "cabin_class": "economy",
      "aircraft": null,
      "amenities": null,
      "baggage": {
        "carry_on": "Cabin Baggage Only",
        "checked": " Checked Bags Additional Fee"
      }
    },
    {
      "id": "QZ520_AirAsia",
      "provider": "AirAsia",
      "airline": {
        "name": "AirAsia",
        "code": "QZ"
      },
      "flight_number": "QZ520",
      "departure": {
        "airport": "CGK",
        "city": "Jakarta",
        "datetime": "2025-12-15T04:45:00+07:00",
        "timestamp": 1765748700
      },
      "arrival": {
        "airport": "DPS",
        "city": "Denpasar",
        "datetime": "2025-12-15T07:25:00+08:00",
        "timestamp": 1765754700
      },
      "duration": {
        "total_minutes": 100,
        "formatted": "01h 40m"
      },
      "stops": 0,
      "price": {
        "amount": 650000,
        "currency": "IDR",
        "formatted": "IDR 650.000"
      },
      "available_seats": 67,
      "cabin_class": "economy",
      "aircraft": null,
      "amenities": null,
      "baggage": {
        "carry_on": "Cabin Baggage Only",
        "checked": " Checked Bags Additional Fee"
      }
    },
    {
      "id": "QZ532_AirAsia",
      "provider": "AirAsia",
      "airline": {
        "name": "AirAsia",
        "code": "QZ"
      },
      "flight_number": "QZ532",
      "departure": {
        "airport": "CGK",
        "city": "Jakarta",
        "datetime": "2025-12-15T19:30:00+07:00",
        "timestamp": 1765801800
      },
      "arrival": {
        "airport": "DPS",
        "city": "Denpasar",
        "datetime": "2025-12-15T22:10:00+08:00",
        "timestamp": 1765807800
      },
      "duration": {
        "total_minutes": 100,
        "formatted": "01h 40m"
      },
      "stops": 0,
      "price": {
        "amount": 595000,
        "currency": "IDR",
        "formatted": "IDR 595.000"
      },
      "available_seats": 72,
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
