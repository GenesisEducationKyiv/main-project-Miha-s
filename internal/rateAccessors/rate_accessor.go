package rateAccessors

type RateAccessor interface {
	GetCurrentRate(currencyFrom string, currencyTo string) (float64, error)
}
