package rateAccessors

type RateAccessor interface {
	GetCurrentRate() (float64, error)
}
