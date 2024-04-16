package currency

import (
	"context"

	"github.com/stretchr/testify/mock"
)

type MockRepo struct {
	mock.Mock
}

func (m *MockRepo) SelectAllCurrencies(ctx context.Context) ([]Currency, error) {
	args := m.Called(ctx)
	return args.Get(0).([]Currency), args.Error(1)
}

func (m *MockRepo) SelectCurrency(ctx context.Context, name string) (*Currency, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(*Currency), args.Error(1)
}

func (m *MockRepo) InsertCurrencies(ctx context.Context, currencies []Currency) ([]Currency, error) {
	args := m.Called(ctx, currencies)
	return args.Get(0).([]Currency), args.Error(1)
}

func (m *MockRepo) SelectChangesPerHour(ctx context.Context, name string) (float64, error) {
	args := m.Called(ctx, name)
	return args.Get(0).(float64), args.Error(1)
}

func (m *MockRepo) SetChangesPerHour(ctx context.Context, currencies []Currency) error {
	args := m.Called(ctx, currencies)
	return args.Error(0)
}
