package storages

type ExchangeRates struct {
	Rates map[string]float32
}

type ExchangeRate struct {
	FromCurrency string
	ToCurrency   string
	Rate         float32
}
