package templates

import (
	"btc-test-task/internal/helpers/config"
	"btc-test-task/internal/helpers/models"
	"fmt"
)

type SimpleTextTemplates struct {
	CurrencyFrom string
	CurrencyTo   string
}

func NewSimpleTextTemplates(conf *config.Config) *SimpleTextTemplates {
	return &SimpleTextTemplates{
		CurrencyFrom: conf.CurrencyFrom,
		CurrencyTo:   conf.CurrencyTo,
	}
}

func (template *SimpleTextTemplates) CurrencyRate(rate models.Rate) string {
	return fmt.Sprintf("Current exchage rate from %v to %v is %.2f", template.CurrencyFrom, template.CurrencyTo, rate.Value)
}
