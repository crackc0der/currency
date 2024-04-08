package currency

import (
	"context"
	"fmt"
	"strconv"
)

type RepositoryInterface interface {
	SelectAllCurrencies(context.Context) ([]Currency, error)
	SelectCurrency(context.Context, string) (*Currency, error)
	InsertCurrencies(context.Context, []Currency) ([]Currency, error)
	SelectChangesPerHour(context.Context, string) (float64, error)
}

type Service struct {
	repository RepositoryInterface
}

func NewService(repository RepositoryInterface) *Service {
	return &Service{repository: repository}
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

// func (s Service) SetCurrencies(ctx context.Context, currencies []Currency) ([]Currency, error) {
// 	currencies, err := s.repository.InsertCurrencies(ctx, currencies)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to insert currencies in method SetCurrencies: %w", err)
// 	}

// 	return currencies, nil
// }

func (s Service) SetCurrencies(ctx context.Context, data Data) error {
	var currencies []Currency

	var currency Currency

	var minPrice float64

	var maxPrice float64

	BTCRUB, err := strconv.ParseFloat(data.BTCRUB, 64)
	if err != nil {
		return fmt.Errorf("error in method SetCurrencies: %w", err)
	}

	ETHRUB, err := strconv.ParseFloat(data.ETHRUB, 64)

	if err != nil {
		return fmt.Errorf("error in method SetCurrencies: %w", err)
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

	for i := 0; i < len(currencyData); i++ {
		curr, _ := s.repository.SelectCurrency(ctx, currencyData[i].Name)

		if curr != nil {
			if currencyData[i].Price < curr.CurrencyMinPrice {
				minPrice = currencyData[i].Price
			} else {
				minPrice = curr.CurrencyMinPrice

				if currencyData[i].Price > curr.CurrencyMaxPrice {
					maxPrice = currencyData[i].Price
				} else {
					maxPrice = curr.CurrencyMaxPrice
				}
			}
		} else {
			minPrice = currencyData[i].Price
			maxPrice = currencyData[i].Price
		}

		currency.CurrencyName = currencyData[i].Name
		currency.CurrencyPrice = currencyData[i].Price
		currency.CurrencyMinPrice = minPrice
		currency.CurrencyMaxPrice = maxPrice
		currency.CurrencyPercentageChange = 0.0

		currencies = append(currencies, currency)
	}

	_, err = s.repository.InsertCurrencies(ctx, currencies)
	if err != nil {
		return fmt.Errorf("error in method SetCurrency: %w", err)
	}

	return nil
}

func (s Service) GetChangesPerHour(ctx context.Context, currency string) (float64, error) {
	change, err := s.repository.SelectChangesPerHour(ctx, currency)
	if err != nil {
		return -1, fmt.Errorf("error in method GetChangesPerHour: %w", err)
	}

	return change, nil
}
