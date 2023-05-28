package templates

import "btc-test-task/internal/config"

type Templates interface {
	Init(conf *config.Config) error
	CurrencyRate(float64) string
}
