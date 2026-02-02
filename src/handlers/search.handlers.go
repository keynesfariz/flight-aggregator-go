package handlers

import (
	"bookcabin-app-go/src/models"
	"bookcabin-app-go/src/services"
	"errors"
	"net/http"
	"sort"
	"time"

	"github.com/gin-gonic/gin"
)

func SearchFlights(ctx *gin.Context) {
	var req models.SearchRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := validateAndNormalizeSearchRequest(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	searchService := services.NewSearchService()
	res, err := searchService.Search(ctx, req)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func validateAndNormalizeSearchRequest(req *models.SearchRequest) error {
	if err := validateDate(req.DepartureDate); err != nil {
		return errors.New("invalid departureDate")
	}

	if req.ReturnDate != nil {
		if err := validateDate(*req.ReturnDate); err != nil {
			return errors.New("invalid returnDate")
		}
	}

	// sort Filters.Airlines for flight search caching purpose
	sort.Strings(req.Filters.Airlines)

	if req.Passengers <= 0 {
		req.Passengers = 1
	}

	if req.CabinClass == "" {
		req.CabinClass = "economy"
	}

	if req.SortBy == "" {
		req.SortBy = "best_value"
	}

	if req.SortOrder == "" {
		req.SortOrder = "asc"
	}

	return nil
}

func validateDate(date string) error {
	_, err := time.Parse(time.DateOnly, date)
	return err
}
