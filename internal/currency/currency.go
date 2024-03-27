package currency

import "github.com/google/uuid"

type Currency struct {
	CurrencyID               uuid.UUID `json:"currencyId"`
	CurrencyName             string    `json:"currencyName"`
	CurrencyPrice            float64   `json:"currencyPrice"`
	CurrencyMinPrice         float64   `json:"currencyMinPrice"`
	CurrencyMaxPrice         float64   `json:"currencyMaxPrice"`
	CurrencyPercentageChange float64   `json:"currencyPercentageChange"`
}
