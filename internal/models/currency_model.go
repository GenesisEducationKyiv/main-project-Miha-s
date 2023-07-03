package models

import "fmt"

type Currency struct {
	CurrencyFrom string
	CurrencyTo   string
}

func (curr *Currency) String() string {
	return fmt.Sprintf("%s->%s", curr.CurrencyFrom, curr.CurrencyTo)
}
