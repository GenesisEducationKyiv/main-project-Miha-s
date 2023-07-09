package currencyrate

import (
	"btc-test-task/internal/common/configuration/config"
	"btc-test-task/internal/common/models"
	"time"

	"github.com/pkg/errors"
)

type TimeProvider interface {
	Now() time.Time
}

type RateProvider interface {
	GetCurrentRate(currency *models.Currency) (models.Rate, error)
}

type RateCache struct {
	cacheDuration time.Duration
	rateProvider  RateProvider
	TimeProvider  TimeProvider

	lastRequestTime time.Time
	lastRate        models.Rate
}

func NewRateCache(config *config.Config, rateProvider RateProvider, timeProvider TimeProvider) (*RateCache, error) {
	duration, err := time.ParseDuration(config.RateCacheDuration)
	if err != nil {
		return nil, errors.Wrap(err, "NewRateCache")
	}
	lastRequestTime := time.Now().Add(-duration)

	return &RateCache{
		rateProvider:    rateProvider,
		TimeProvider:    timeProvider,
		cacheDuration:   duration,
		lastRequestTime: lastRequestTime,
	}, nil
}

func (cache *RateCache) GetCurrentRate(currency *models.Currency) (models.Rate, error) {
	diffDuration := cache.TimeProvider.Now().Sub(cache.lastRequestTime)
	if diffDuration >= cache.cacheDuration {
		return cache.rateProvider.GetCurrentRate(currency)
	}

	return cache.lastRate, nil
}
