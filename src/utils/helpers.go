package utils

import (
	"bookcabin-app-go/src/constants"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func GetFlightId(airline string, flightNum string) string {
	re := regexp.MustCompile(`\W+`)
	return flightNum + "_" + re.ReplaceAllString(airline, "")
}

func CapitalizeFirst(s string) string {
	return cases.Title(language.English).String(s)
}

func FormatDurationToHumans(durationInMin int) string {
	d := time.Duration(durationInMin) * time.Minute
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	return fmt.Sprintf("%02dh %02dm", h, m)
}

func FormatDurationToMinutes(duration string) int {
	parts := strings.Fields(duration)
	if len(parts) < 2 {
		return 0
	}

	hours, _ := strconv.Atoi(strings.TrimSuffix(parts[0], "h"))
	minutes, _ := strconv.Atoi(strings.TrimSuffix(parts[1], "m"))

	return hours*60 + minutes
}

func FormatPrice(price int, currency string) string {
	p := message.NewPrinter(language.Indonesian)
	formattedPrice := p.Sprintf("%d", price)
	return currency + " " + formattedPrice
}

func FormatDateTime(dateTime time.Time, airportCode string) string {
	tz, _ := time.LoadLocation(constants.Timezones[airportCode])
	formattedDateTime := dateTime.In(tz).Format(constants.GA_DateTimeLayout)
	return formattedDateTime
}
