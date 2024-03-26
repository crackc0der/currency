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

func (s Service) GetCurrency(ctx context.Context, currency string) (Currency, error) {
}

func (s Service) SetCurrencies(ctx context.Context, currencies []Currency) ([]Currency, error) {
}

func (s Service) GetChangesPerHour(ctx context.Context, currency string) (Currency, error) {
}
