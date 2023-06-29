package rateAccessorsTest

type RateAccessorStub struct {
	RateError error
	Rate      float64
}

func (accessor *RateAccessorStub) GetCurrentRate(currencyFrom string, currencyTo string) (float64, error) {
	return accessor.Rate, accessor.RateError
}
