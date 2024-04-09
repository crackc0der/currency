package currency

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/crackc0der/currency/config"
	"github.com/gorilla/mux"
)

type EndpointInterface interface {
	GetCurrencies(context.Context) ([]Currency, error)
	GetCurrency(context.Context, string) (*Currency, error)
	SetCurrencies(context.Context, []Currency) ([]Currency, error)
	GetChangesPerHour(context.Context, string) (float64, error)
}

type Endpoint struct {
	service *Service
	log     *slog.Logger
	config  *config.Config
}

func NewEndpoint(service *Service, log *slog.Logger, config *config.Config) *Endpoint {
	return &Endpoint{service: service, log: log, config: config}
}

func (e Endpoint) GetCurrencies(writer http.ResponseWriter, request *http.Request) {
	currencies, err := e.service.GetCurrencies(request.Context())
	if err != nil {
		e.log.Error("error in Endpoint's method GetCurrencies: " + err.Error())
	}

	if err = json.NewEncoder(writer).Encode(&currencies); err != nil {
		e.log.Error("error in Endpoint's method GetCurrencies: " + err.Error())
	}
}

func (e Endpoint) GetCurrency(writer http.ResponseWriter, request *http.Request) {
	currencyName := mux.Vars(request)["name"]

	currency, err := e.service.GetCurrency(request.Context(), currencyName)
	if err != nil {
		e.log.Error("error in Endpoint's method GetCurrency: " + err.Error())
	}

	if err = json.NewEncoder(writer).Encode(&currency); err != nil {
		e.log.Error("error in Endpoint's method GetCurrency: " + err.Error())
	}
}

func (e Endpoint) GetChangesPerHour(writer http.ResponseWriter, request *http.Request) {
	currencyName := mux.Vars(request)["name"]

	currencyChange, err := e.service.GetChangesPerHour(request.Context(), currencyName)
	if err != nil {
		e.log.Error("error in Endpoint's method GetChangesPerHour: " + err.Error())
	}

	if err = json.NewEncoder(writer).Encode(&currencyChange); err != nil {
		e.log.Error("error in Endpoint's method GetChangesPerHour: " + err.Error())
	}
}
