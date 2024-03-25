package currency

import "github.com/google/uuid"

type Currency struct {
	CurrencyID               uuid.UUID `json:"currency_id"`
	CurrencyName             string    `json:"currency_name"`
	CurrencyPrice            float64   `json:"currency_price"`
	CurrencyMinPrice         float64   `json:"currency_min_price"`
	CurrencyMaxPrice         float64   `json:"currency_max_price"`
	CurrencyPercentageChange float64   `json:"currency_percentage_change"`
}
