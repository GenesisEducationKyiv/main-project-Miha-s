package templates

type Templates interface {
	CurrencyRate(float64) string
}
