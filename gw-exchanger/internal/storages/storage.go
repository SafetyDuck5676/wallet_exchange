package storages

type Storage interface {
	GetExchangeRates() (map[string]float64, error)
	GetExchangeRate(fromCurrency, toCurrency string) (float64, error)
}
