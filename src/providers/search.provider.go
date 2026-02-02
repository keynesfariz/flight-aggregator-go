package providers

import (
	"bookcabin-app-go/src/constants"
	"bookcabin-app-go/src/libs"
	"bookcabin-app-go/src/models"
	"bookcabin-app-go/src/utils"
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	mockBasePath = "src/mocks/"
)

type SearchProviderProperty struct {
	Name         string
	SuccessRate  int    // in percentage
	ResponseTime [2]int // in miliseconds
	MockFile     string
}

type SearchProvider interface {
	Fetch(ctx *gin.Context, req models.SearchRequest) ([]models.Flight, error)
}

/* Fetch Simulation */
func SimulateFetchWithWait(ctx context.Context, pvd SearchProviderProperty) ([]byte, error) {
	// fetch from cache here...
	cache := libs.GetCacheClientInstance()
	cachedData, err := cache.Get(ctx, pvd.Name).Result()
	if err == nil {
		return []byte(cachedData), nil
	}

	// simulate fetching from provider
	maxRetry, _ := strconv.Atoi(libs.GetEnv("FLIGHT_PROVIDER_MAX_RETRY", "3"))

	// create random base source
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	// make some attempts to fetch flight data from a provider
	for attempt := 0; attempt <= maxRetry; attempt++ {
		if attempt > 0 {
			if err := backoffWait(ctx, attempt); err != nil {
				return nil, err
			}
		}

		if err := simulateServerLatency(ctx, pvd.ResponseTime, rng); err != nil {
			return nil, err
		}

		if fetchSucceeded(pvd.SuccessRate, rng) {
			return readMockFile(ctx, pvd)
		}
	}

	return nil, fmt.Errorf("failed fetching flight data from provider: %s", pvd.Name)
}

func fetchSucceeded(successRate int, rng *rand.Rand) bool {
	if successRate >= 100 {
		return true
	}
	return rng.Intn(100) < successRate
}

func randomDuration(minMs, maxMs int, rng *rand.Rand) time.Duration {
	if maxMs <= minMs {
		return time.Duration(minMs) * time.Millisecond
	}
	return time.Duration(minMs+rng.Intn(maxMs-minMs)) * time.Millisecond
}

func simulateServerLatency(
	ctx context.Context,
	responseTime [2]int,
	rng *rand.Rand,
) error {
	delay := randomDuration(responseTime[0], responseTime[1], rng)

	select {
	case <-time.After(delay):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func backoffWait(ctx context.Context, attempt int) error {
	baseBackoffMs, _ := strconv.Atoi(libs.GetEnv("FLIGHT_PROVIDER_BACKOFF_IN_MS", "8"))

	exponentDelayTime := baseBackoffMs * (1 << attempt)
	delay := time.Duration(exponentDelayTime) * time.Millisecond

	select {
	case <-time.After(delay):
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func readMockFile(ctx context.Context, pvd SearchProviderProperty) ([]byte, error) {
	if pvd.MockFile == "" {
		return nil, errors.New("mock file not defined")
	}

	data, err := os.ReadFile(mockBasePath + pvd.MockFile)

	if err == nil {
		cache := libs.GetCacheClientInstance()
		cache.Set(ctx, pvd.Name, string(data), 5*time.Minute)
	}

	return data, err
}

func getCity(airport string) string {
	if city, ok := constants.Cities[airport]; ok {
		return city
	}
	return "Unknown city"
}

func parseBaggageInfo(baggageInfo string) models.Baggage {
	baggageInfos := strings.Split(baggageInfo, ",")

	if len(baggageInfos) < 2 {
		return models.Baggage{
			CarryOn: nil,
			Checked: nil,
		}
	}

	carryOn := utils.CapitalizeFirst(baggageInfos[0])
	checked := utils.CapitalizeFirst(baggageInfos[1])

	return models.Baggage{
		CarryOn: &carryOn,
		Checked: &checked,
	}
}
