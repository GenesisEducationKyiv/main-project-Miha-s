package templatesTest

import (
	"btc-test-task/internal/helpers/models"
)

type TemplatesImplStub struct {
}

func (template *TemplatesImplStub) CurrencyRate(rate models.Rate) string {
	return ""
}
