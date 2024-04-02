package currency

import (
	"context"
	"encoding/json"
	"exchange_course/config"
	"fmt"
	"io/ioutil"
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
	config  *config.Config
}

func NewEndpoint(service Service, log *slog.Logger, config *config.Config) *Endpoint {
	return &Endpoint{service: service, log: log, config: config}
}

func (e Endpoint) GetCurrencies(writer http.ResponseWriter, request *http.Request) {
	currencies, err := e.service.GetCurrencies(request.Context())
	if err != nil {
		e.log.Error("Error getting currencies in method GetCurrencies: " + err.Error())
	}

	if err = json.NewEncoder(writer).Encode(&currencies); err != nil {
		e.log.Error("Error encoding currency: " + err.Error())
	}
}

func (e Endpoint) GetCurrency(writer http.ResponseWriter, request *http.Request) {
	currencyName := mux.Vars(request)["name"]

	currency, err := e.service.GetCurrency(request.Context(), currencyName)
	if err != nil {
		e.log.Error("Error getting currency: " + err.Error())
	}

	if err = json.NewEncoder(writer).Encode(&currency); err != nil {
		e.log.Error("Error encoding currency: " + err.Error())
	}
}

func (e Endpoint) GetChangesPerHour(writer http.ResponseWriter, request *http.Request) {
	currencyName := mux.Vars(request)["name"]

	currencyChange, err := e.service.GetChangesPerHour(request.Context(), currencyName)
	if err != nil {
		e.log.Error("Error getting currency: " + err.Error())
	}

	if err = json.NewEncoder(writer).Encode(&currencyChange); err != nil {
		e.log.Error("Error encoding currencyChange: " + err.Error())
	}
}

func (e Endpoint) CurrencyMonitor() {
	url := "https://currate.ru/api/?get=rates&pairs=BTCRUB,ETHRUB&key=" + e.config.ApiKey

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		e.log.Error("error creating request: " + err.Error())
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		e.log.Error("error doing request: " + err.Error())
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		e.log.Error("error reading response body: " + err.Error())
	}
	fmt.Println("response Body:", string(body))
}
