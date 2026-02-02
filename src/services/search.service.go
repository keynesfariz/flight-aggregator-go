package services

import (
	"bookcabin-app-go/src/libs"
	"bookcabin-app-go/src/models"
	"bookcabin-app-go/src/providers"
	"bookcabin-app-go/src/utils"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type SearchService struct {
	providers []providers.SearchProvider
}

func NewSearchService() *SearchService {
	return &SearchService{[]providers.SearchProvider{
		providers.NewAirAsiaProvider(),
		providers.NewBatikAirProvider(),
		providers.NewGarudaIndonesiaProvider(),
		providers.NewLionAirProvider(),
	}}
}

func (s *SearchService) Search(ctx *gin.Context, req models.SearchRequest) (models.SearchResponse, error) {
	start := time.Now()

	cache := libs.GetCacheClientInstance()
	cacheKey := getCacheKeyFromSearchRequest(req)
	result, err := cache.Get(ctx, cacheKey).Result()

	if err == nil {
		var cachedResult models.SearchResponse
		err := json.Unmarshal([]byte(result), &cachedResult)

		if err == nil {
			cachedResult.Metadata.CacheHit = true
			return cachedResult, err
		}
	}

	var wg sync.WaitGroup
	resultsCh := make(chan []models.Flight, len(s.providers))
	errorsCh := make(chan error, len(s.providers))

	for _, p := range s.providers {
		wg.Add(1)
		go func(p providers.SearchProvider) {
			defer wg.Done()
			flights, err := p.Fetch(ctx, req)
			if err != nil {
				errorsCh <- err
				return
			}
			resultsCh <- flights
		}(p)
	}

	wg.Wait()
	close(resultsCh)
	close(errorsCh)

	// merge & normalize results
	flights := []models.Flight{}
	for f := range resultsCh {
		flights = append(flights, f...)
	}

	// filter
	utils.ApplySearchFilters(&flights, req)

	// sorting, scoring
	utils.ApplySearchSorter(flights, req.SortBy, req.SortOrder)

	results := models.SearchResponse{
		Criteria: req,
		Metadata: models.Metadata{
			TotalResults:     len(flights),
			ProvidersQueried: len(s.providers),
			ProvidersSuccess: len(s.providers) - len(errorsCh),
			ProvidersFailed:  len(errorsCh),
			SearchTimeMs:     int(time.Since(start).Milliseconds()),
		},
		Flights: flights,
	}

	resultToCache, err := json.Marshal(results)
	if err == nil {
		cache.Set(ctx, cacheKey, resultToCache, 5*time.Minute)
	}

	return results, nil
}

func getCacheKeyFromSearchRequest(req models.SearchRequest) string {
	reqJson, _ := json.Marshal(req)
	hash := sha256.Sum256(reqJson)
	return fmt.Sprintf("Q:%x", hash)
}
