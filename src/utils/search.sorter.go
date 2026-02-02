package utils

import (
	"bookcabin-app-go/src/constants"
	"bookcabin-app-go/src/models"

	"sort"
	"time"
)

func ApplySearchSorter(flights []models.Flight, sortBy string, sortOrder string) {
	if len(flights) == 0 {
		return
	}

	if sortBy == "best_value" {
		scoreResults := GetBestValueScores(flights)

		sort.Slice(flights, func(i, j int) bool {
			scoreI := scoreResults[flights[i].ID]
			scoreJ := scoreResults[flights[j].ID]

			// force descending sort bigger score = better value
			return scoreJ > scoreI
		})
		return
	}

	sort.Slice(flights, func(i, j int) bool {
		var compareVal int

		switch sortBy {
		case "price":
			if flights[i].Price.Amount < flights[j].Price.Amount {
				compareVal = -1
			} else if flights[i].Price.Amount > flights[j].Price.Amount {
				compareVal = 1
			}
		case "duration":
			if flights[i].Duration.TotalMinutes < flights[j].Duration.TotalMinutes {
				compareVal = -1
			} else if flights[i].Duration.TotalMinutes > flights[j].Duration.TotalMinutes {
				compareVal = 1
			}
		case "departure":
			departTimeI, _ := time.Parse(constants.GA_DateTimeLayout, flights[i].Departure.DateTime)
			departTimeJ, _ := time.Parse(constants.GA_DateTimeLayout, flights[j].Departure.DateTime)
			if departTimeI.Before(departTimeJ) {
				compareVal = -1
			} else if departTimeI.After(departTimeJ) {
				compareVal = 1
			}
		case "arrival":
			arriveTimeI, _ := time.Parse(constants.GA_DateTimeLayout, flights[i].Arrival.DateTime)
			arriveTimeJ, _ := time.Parse(constants.GA_DateTimeLayout, flights[j].Arrival.DateTime)
			if arriveTimeI.Before(arriveTimeJ) {
				compareVal = -1
			} else if arriveTimeI.After(arriveTimeJ) {
				compareVal = 1
			}
		}

		if sortOrder == "desc" {
			compareVal = compareVal * -1
		}

		return compareVal < 0
	})
}
