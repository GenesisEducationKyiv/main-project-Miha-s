package currencyrate

import (
	"btc-test-task/internal/common/configuration/config"
	"btc-test-task/internal/common/configuration/logger"
	"btc-test-task/internal/common/models"
	"time"
)

type CacheProvider interface {
	Get(k string) (interface{}, bool)
	Add(k string, x interface{}, d time.Duration) error
}

type RateProvider interface {
	GetCurrentRate(currency *models.Currency) (models.Rate, error)
}

type RateCache struct {
	rateProvider RateProvider
	cache        CacheProvider
	expiration   time.Duration
}

func NewRateCache(config *config.Config, rateProvider RateProvider, cache CacheProvider) (*RateCache, error) {
	return &RateCache{
		rateProvider: rateProvider,
		cache:        cache,
		expiration:   config.RateCacheDuration,
	}, nil
}

func (c *RateCache) GetCurrentRate(currency *models.Currency) (models.Rate, error) {
	cachedRate, exists := c.cache.Get(c.key(currency))
	if exists {
		return cachedRate.(models.Rate), nil
	}
	rate, err := c.rateProvider.GetCurrentRate(currency)
	if err != nil {
		return models.Rate{}, err
	}

	err = c.cache.Add(c.key(currency), rate, c.expiration)
	if err != nil {
		logger.Log.Errorf("Failed to save cache %v", err)
	}

	return rate, nil
}

func (c *RateCache) key(currency *models.Currency) string {
	return currency.From + currency.To
}
