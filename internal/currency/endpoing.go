package currency

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
)

type EndpointInterface interface {
	GetCurrencies(context.Context) ([]Currency, error)
	GetCurrency(context.Context, string) (*Currency, error)
	SetCurrencies(context.Context, []Currency) ([]Currency, error)
	GetChangesPerHour(context.Context, string) (float64, error)
}

type Endpoint struct {
	service Service
	log     *slog.Logger
}

func NewEndpoint(service Service, log *slog.Logger) *Endpoint {
	return &Endpoint{service: service, log: log}
}

func (e Endpoint) GetCurrencies(w http.ResponseWriter, r *http.Request) {
	currencies, err := e.service.GetCurrencies(r.Context())

	if err != nil {
		e.log.Error("Error getting currencies in method GetCurrencies: " + err.Error())
	}

	if err = json.NewEncoder(w).Encode(&currencies); err != nil {
		e.log.Error("Error encoding currency: " + err.Error())
	}
}

func (e Endpoint) GetCurrency(w http.ResponseWriter, r *http.Request) {
	currencyName := mux.Vars(r)["name"]

	currency, err := e.service.GetCurrency(r.Context(), currencyName)

	if err != nil {
		e.log.Error("Error getting currency: " + err.Error())
	}

	if err = json.NewEncoder(w).Encode(&currency); err != nil {
		e.log.Error("Error encoding currency: " + err.Error())
	}
}

func (e Endpoint) GetChangesPerHour(w http.ResponseWriter, r *http.Request) {
	currencyName := mux.Vars(r)["name"]

	currencyChange, err := e.service.GetChangesPerHour(r.Context(), currencyName)

	if err != nil {
		e.log.Error("Error getting currency: " + err.Error())
	}

	if err = json.NewEncoder(w).Encode(&currencyChange); err != nil {
		e.log.Error("Error encoding currencyChange: " + err.Error())
	}
}
