package templates

import (
	"btc-test-task/internal/helpers/config"
	"fmt"
)

type SimpleTextTemplates struct {
	CurrencyFrom string
	CurrencyTo   string
}

func NewSimpleTextTemplates(conf *config.Config) (*SimpleTextTemplates, error) {
	return &SimpleTextTemplates{
		CurrencyFrom: conf.CurrencyFrom,
		CurrencyTo:   conf.CurrencyTo,
	}, nil
}

func (template *SimpleTextTemplates) CurrencyRate(rate float64) string {
	return fmt.Sprintf("Current exchage rate from %v to %v is %.2f", template.CurrencyFrom, template.CurrencyTo, rate)
}
