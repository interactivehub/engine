package currency

import "github.com/Rhymond/go-money"

var (
	HubCurrency = money.AddCurrency("HUB", "HUB", "1 $", ".", ",", 2)
)

type HubMoney struct {
	*money.Money
}

func NewHubMoney(amount float64) *HubMoney {
	return &HubMoney{
		money.NewFromFloat(amount, HubCurrency.Code),
	}
}
