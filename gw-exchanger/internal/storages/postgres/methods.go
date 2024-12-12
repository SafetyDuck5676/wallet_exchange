package postgres

import (
	"fmt"
)

func (ps *PostgresStorage) GetExchangeRates() (map[string]float64, error) {
	rows, err := ps.db.Query("SELECT from_currency, rate FROM exchange_rates")
	if err != nil {
		return nil, fmt.Errorf("failed to query exchange rates: %w", err)
	}
	defer rows.Close()

	rates := make(map[string]float64)
	for rows.Next() {
		var currency string
		var rate float64
		if err := rows.Scan(&currency, &rate); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		rates[currency] = rate
	}

	return rates, nil
}

func (ps *PostgresStorage) GetExchangeRate(fromCurrency, toCurrency string) (float64, error) {
	var rate float64
	query := "SELECT rate FROM exchange_rates WHERE from_currency = $1 AND to_currency = $2"
	err := ps.db.QueryRow(query, fromCurrency, toCurrency).Scan(&rate)
	if err != nil {
		return 0, fmt.Errorf("failed to get exchange rate: %w", err)
	}

	return rate, nil
}
