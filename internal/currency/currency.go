package currency

import (
	"time"
)

type Currency struct {
	CurrencyID            int64     `json:"currencyId"`
	CurrencyName          string    `json:"currencyName"`
	CurrencyPrice         float64   `json:"currencyPrice"`
	CurrencyMinPrice      float64   `json:"currencyMinPrice"`
	CurrencyMaxPrice      float64   `json:"currencyMaxPrice"`
	CurrencyChangePerHour float64   `json:"currencyChangePerHour"`
	CurrencyLastUpdate    time.Time `json:"currencyLastUpdate"`
}

type DataCurrencyMonitor struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
	Data    Data   `json:"data"`
}

type Data struct {
	BTCRUB string `json:"BTCRUB"`
	ETHRUB string `json:"ETHRUB"`
}

type LastUpdate struct {
	BTC time.Time
	ETH time.Time
}
