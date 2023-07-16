package currencyrate

import (
	"btc-test-task/internal/common/configuration/config"
	"btc-test-task/internal/common/configuration/logger"
	"btc-test-task/internal/common/models"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/pkg/errors"
)

type TimeProvider interface {
	Now() time.Time
}

type RateProvider interface {
	GetCurrentRate(currency *models.Currency) (models.Rate, error)
}

type RateCache struct {
	rateProvider RateProvider
	cache        *cache.Cache
}

func NewRateCache(config *config.Config, rateProvider RateProvider, timeProvider TimeProvider) (*RateCache, error) {
	duration, err := time.ParseDuration(config.RateCacheDuration)
	if err != nil {
		return nil, errors.Wrap(err, "NewRateCache")
	}

	return &RateCache{
		rateProvider: rateProvider,
		cache:        cache.New(duration, duration),
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

	err = c.cache.Add(c.key(currency), rate, cache.DefaultExpiration)
	if err != nil {
		logger.Log.Errorf("Failed to save cache %v", err)
	}

	return rate, nil
}

func (c *RateCache) key(currency *models.Currency) string {
	return currency.From + currency.To
}
