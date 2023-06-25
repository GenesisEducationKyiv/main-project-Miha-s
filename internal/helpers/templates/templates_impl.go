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
	newSimpleTextTemplates := new(SimpleTextTemplates)
	err := newSimpleTextTemplates.init(conf)
	if err != nil {
		return nil, err
	}

	return newSimpleTextTemplates, nil
}

func (template *SimpleTextTemplates) init(conf *config.Config) error {
	template.CurrencyFrom = conf.CurrencyFrom
	template.CurrencyTo = conf.CurrencyTo
	return nil
}

func (template *SimpleTextTemplates) CurrencyRate(rate float64) string {
	return fmt.Sprintf("Current exchage rate from %v to %v is %.2f", template.CurrencyFrom, template.CurrencyTo, rate)
}
