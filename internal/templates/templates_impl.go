package templates

import (
	"btc-test-task/internal/config"
	"fmt"
)

type TemplatesImpl struct {
	CurrencyFrom string
	CurrencyTo   string
}

func (template *TemplatesImpl) Init(conf *config.Config) error {
	template.CurrencyFrom = conf.CurrencyFrom
	template.CurrencyTo = conf.CurrencyTo
	return nil
}

func (template *TemplatesImpl) CurrencyRate(rate float64) string {
	return fmt.Sprintf("Current exchage rate from %v to %v is %.2f", template.CurrencyFrom, template.CurrencyTo, rate)
}
