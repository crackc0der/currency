package currency

import (
	"context"
	"encoding/json"
	"exchange_course/config"
	"io"
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
		e.log.Error("error in method GetCurrencies: " + err.Error())
	}

	if err = json.NewEncoder(writer).Encode(&currencies); err != nil {
		e.log.Error("error in method GetCurrencies: " + err.Error())
	}
}

func (e Endpoint) GetCurrency(writer http.ResponseWriter, request *http.Request) {
	currencyName := mux.Vars(request)["name"]

	currency, err := e.service.GetCurrency(request.Context(), currencyName)
	if err != nil {
		e.log.Error("error in method GetCurrency: " + err.Error())
	}

	if err = json.NewEncoder(writer).Encode(&currency); err != nil {
		e.log.Error("error in method GetCurrency: " + err.Error())
	}
}

func (e Endpoint) GetChangesPerHour(writer http.ResponseWriter, request *http.Request) {
	currencyName := mux.Vars(request)["name"]

	currencyChange, err := e.service.GetChangesPerHour(request.Context(), currencyName)
	if err != nil {
		e.log.Error("error in method GetChangesPerHour: " + err.Error())
	}

	if err = json.NewEncoder(writer).Encode(&currencyChange); err != nil {
		e.log.Error("error in method GetChangesPerHour: " + err.Error())
	}
}

func (e Endpoint) SetCurrencies() {
}

func (e Endpoint) CurrencyMonitor() {
	url := "https://currate.ru/api/?get=rates&pairs=BTCRUB,ETHRUB&key=" + e.config.APIKey
	var data DataCurrencyMonitor
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		e.log.Error("error in method CurrencyMonitor: " + err.Error())
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		e.log.Error("error in method CurrencyMonitor: " + err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		e.log.Error("error in method CurrencyMonitor: " + err.Error())
	}

	err = json.Unmarshal(body, &data)

	if err != nil {
		e.log.Error("error in method CurrencyMonitor: ", err)
	}

	err = e.service.SetCurrencies(context.Background(), data.Data)
	if err != nil {
		e.log.Error("error in method CurrentMonitor: %w", err)
	}
}
