package currency

import (
	"context"
	"fmt"
)

type Service struct {
	repository Repository
}

func (s Service) GetCurrencies(ctx context.Context) ([]Currency, error) {
	currencies, err := s.repository.SelectAllCurrencies(ctx)

	if err != nil {
		return nil, fmt.Errorf("error to get currencies in service's method GetCurrencies %w", err)
	}

	return currencies, nil
}

func (s Service) GetCurrency(ctx context.Context, currencyName string) (*Currency, error) {
	currency, err := s.repository.SelectCurrency(ctx, currencyName)

	if err != nil {
		return nil, fmt.Errorf("failed to get currency in method GetCurrency: %w", err)
	}

	return currency, nil
}

func (s Service) SetCurrencies(ctx context.Context, currencies []Currency) ([]Currency, error) {
	currencies, err := s.repository.InsertCurrencies(ctx, currencies)

	if err != nil {
		return nil, fmt.Errorf("failed to insert currencies in method SetCurrencies: %w", err)
	}

	return currencies, nil
}

func (s Service) GetChangesPerHour(ctx context.Context, currency string) (float64, error) {
	change, err := s.repository.SelectChangesPerHour(ctx, currency)

	if err != nil {
		return -1, fmt.Errorf("error getting changes per hour in method GetChangesPerHour: %w", err)
	}

	return change, nil
}
