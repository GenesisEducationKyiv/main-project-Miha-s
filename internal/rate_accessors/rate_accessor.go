package rate_accessors

import (
	"btc-test-task/internal/config"
)

type RateAccessor interface {
	Init(*config.Config) error
	GetCurrentRate() (float64, error)
}
