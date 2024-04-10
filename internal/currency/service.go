package currency

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/crackc0der/currency/config"
)

type RepositoryInterface interface {
	SelectAllCurrencies(context.Context) ([]Currency, error)
	SelectCurrency(context.Context, string) (*Currency, error)
	InsertCurrencies(context.Context, []Currency) ([]Currency, error)
	SelectChangesPerHour(context.Context, string) (float64, error)
	SetChangesPerHour(context.Context, []Currency) error
}

type Service struct {
	repository RepositoryInterface
	log        *slog.Logger
	config     *config.Config
}

func NewService(repository RepositoryInterface, log *slog.Logger, config *config.Config) *Service {
	return &Service{repository: repository, log: log, config: config}
}

func (s Service) GetCurrencies(ctx context.Context) ([]Currency, error) {
	currencies, err := s.repository.SelectAllCurrencies(ctx)
	if err != nil {
		return nil, fmt.Errorf("error in method GetCurrencies %w", err)
	}

	return currencies, nil
}

func (s Service) GetCurrency(ctx context.Context, currencyName string) (*Currency, error) {
	currency, err := s.repository.SelectCurrency(ctx, currencyName)
	if err != nil {
		return nil, fmt.Errorf("error in method GetCurrency: %w", err)
	}

	return currency, nil
}

func (s Service) SetCurrencies(ctx context.Context, data Data) error {
	currencies, err := s.getCurrentPrice(ctx, data)
	if err != nil {
		return fmt.Errorf("error in Service's method SetCurrency: %w", err)
	}

	_, err = s.repository.InsertCurrencies(ctx, currencies)
	if err != nil {
		return fmt.Errorf("error in Service's method SetCurrency: %w", err)
	}

	return nil
}

func (s Service) GetChangesPerHour(ctx context.Context, currency string) (float64, error) {
	change, err := s.repository.SelectChangesPerHour(ctx, currency)
	if err != nil {
		return -1, fmt.Errorf("error in Service's method GetChangesPerHour: %w", err)
	}

	return change, nil
}

func (s Service) getCurrentPrice(ctx context.Context, data Data) ([]Currency, error) {
	var currency Currency

	var minPrice float64

	var maxPrice float64

	currencyAmount := 2

	currencies := make([]Currency, 0, currencyAmount)

	BTCRUB, err := strconv.ParseFloat(data.BTCRUB, 64)
	if err != nil {
		return nil, fmt.Errorf("error in Service's method SetCurrencies: %w", err)
	}

	ETHRUB, err := strconv.ParseFloat(data.ETHRUB, 64)
	if err != nil {
		return nil, fmt.Errorf("error in Service's method SetCurrencies: %w", err)
	}

	currencyData := [2]struct {
		Name  string
		Price float64
	}{
		{
			"BTC",
			BTCRUB,
		},
		{
			"ETH",
			ETHRUB,
		},
	}

	for _, curr := range currencyData {
		currentData, _ := s.repository.SelectCurrency(ctx, curr.Name)

		if currentData == nil {
			minPrice = curr.Price
			maxPrice = curr.Price
		} else {
			minPrice = s.updateMinPrice(curr.Price, currentData.CurrencyMinPrice)
			maxPrice = s.updateMaxPrice(curr.Price, currentData.CurrencyMaxPrice)
		}

		currency.CurrencyName = curr.Name
		currency.CurrencyPrice = curr.Price
		currency.CurrencyMinPrice = minPrice
		currency.CurrencyMaxPrice = maxPrice
		currency.CurrencyChangePerHour = 0.0

		currencies = append(currencies, currency)
	}

	return currencies, nil
}

func (s Service) CurrencyMonitor() {
	data := s.getMonitorData()

	err := s.SetCurrencies(context.Background(), data.Data)
	if err != nil {
		s.log.Error("error in Endpoint's method CurrentMonitor: %w", err)
	}
}

func (s Service) updateMinPrice(currPrice, currentMinPrice float64) float64 {
	if currPrice < currentMinPrice {
		return currPrice
	}

	return currentMinPrice
}

func (s Service) updateMaxPrice(currPrice, currentMaxPrice float64) float64 {
	if currPrice > currentMaxPrice {
		return currPrice
	}

	return currentMaxPrice
}

func (s Service) getMonitorData() *DataCurrencyMonitor {
	var data DataCurrencyMonitor

	timeout := 5

	timeOutClient := 3

	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: time.Duration(timeout) * time.Second,
		}).Dial,
		TLSHandshakeTimeout: time.Duration(timeout) * time.Second,
	}

	url := "https://currate.ru/api/?get=rates&pairs=BTCRUB,ETHRUB&key=" + s.config.APIKey

	req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, url, nil)
	if err != nil {
		s.log.Error("error in Endpoint's method CurrencyMonitor: " + err.Error())
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{
		Timeout:   time.Duration(timeOutClient) * time.Second,
		Transport: transport,
	}

	resp, err := client.Do(req)
	if err != nil {
		s.log.Error("error in Endpoint's method CurrencyMonitor: " + err.Error())
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.log.Error("error in Endpoint's method CurrencyMonitor: " + err.Error())
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		s.log.Error("error in Endpoint's method CurrencyMonitor: ", err)
	}

	return &data
}

func (s Service) SetChangesPerHour() {
	var currency []Currency

	currenciesInDB, err := s.repository.SelectAllCurrencies(context.Background())
	if err != nil {
		s.log.Error("error in method GetCurrency: %w", err)
	}

	currentData := s.getMonitorData()

	if err != nil {
		s.log.Error("error in method SetChangesPerHourn: %w", err)
	}

	for _, curr := range currenciesInDB {
		if curr.CurrencyName == "BTC" {
			data, err := strconv.ParseFloat(currentData.Data.BTCRUB, 64)
			if err != nil {
				s.log.Error("error in method SetChangesPerHourn: %w", err)
			}

			curr.CurrencyChangePerHour = data - curr.CurrencyPrice
			currency = append(currency, curr)
		}

		if curr.CurrencyName == "ETH" {
			data, err := strconv.ParseFloat(currentData.Data.ETHRUB, 64)
			if err != nil {
				s.log.Error("error in method SetChangesPerHourn: %w", err)
			}

			curr.CurrencyChangePerHour = data - curr.CurrencyPrice
			currency = append(currency, curr)
		}
	}

	err = s.repository.SetChangesPerHour(context.Background(), currency)
	if err != nil {
		s.log.Error("error in Service's method SetChangesPerHour: %w", err)
	}
}
