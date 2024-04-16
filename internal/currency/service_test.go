package currency

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSelectAllCurrencies(t *testing.T) {
	// SelectAllCurrencies(ctx context.Context) ([]Currency, error)
	repo := new(MockRepo)

	setList := func(currencies []Currency, err error) {
		repo.On("SelectAllCurrencies", mock.Anything).Return(currencies, err).Once()
	}

	const layout = "2006-01-02T15:04:05Z07:00"
	tm, _ := time.Parse(layout, "2024-04-16T12:00:00+03:00")

	tests := []struct {
		name    string
		setup   func()
		want    []Currency
		wantErr error
	}{
		{
			name: "success",
			setup: func() {
				setList(
					[]Currency{
						{
							CurrencyID:            1,
							CurrencyName:          "USD",
							CurrencyPrice:         90.42,
							CurrencyMinPrice:      27.00,
							CurrencyMaxPrice:      120.00,
							CurrencyChangePerHour: 11.42,
							CurrencyLastUpdate:    tm,
						},
					},
					nil,
				)
			},
			want: []Currency{
				{
					CurrencyID:            1,
					CurrencyName:          "USD",
					CurrencyPrice:         90.42,
					CurrencyMinPrice:      27.00,
					CurrencyMaxPrice:      120.00,
					CurrencyChangePerHour: 11.42,
					CurrencyLastUpdate:    tm,
				},
			},
			wantErr: nil,
		},
		{
			name: "some error",
			setup: func() {
				setList(nil, errors.New("no currencies"))
			},
			want:    nil,
			wantErr: errors.New("error in method GetCurrencies no currencies"),
		},
		{
			name: "empty list",
			setup: func() {
				setList(nil, nil)
			},
			want:    nil,
			wantErr: nil,
		},
	}

	svc := NewService(repo, slog.New(slog.NewTextHandler(os.Stdout, nil)), nil)

	for _, tt := range tests {
		t.Run(tt.name, func(y *testing.T) {
			defer repo.AssertExpectations(t)

			tt.setup()

			got, err := svc.GetCurrencies(context.Background())
			if err != nil && assert.Error(t, tt.wantErr, err.Error()) {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, tt.wantErr)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSelectCurrency(t *testing.T) {
	// SelectCurrency(ctx context.Context, name string) (*Currency, error)
	repo := new(MockRepo)

	type args struct {
		name string
	}

	setList := func(currency *Currency, name string, err error) {
		repo.On("SelectCurrency", mock.Anything, name).Return(currency, err)
	}

	const layout = "2006-01-02T15:04:05Z07:00"
	tm, _ := time.Parse(layout, "2024-04-16T12:00:00+03:00")

	tests := []struct {
		name    string
		setup   func()
		args    args
		want    *Currency
		wantErr error
	}{
		{
			name: "success",
			setup: func() {
				setList(
					&Currency{
						CurrencyID:            1,
						CurrencyName:          "USD",
						CurrencyPrice:         90.42,
						CurrencyMinPrice:      27.00,
						CurrencyMaxPrice:      120.00,
						CurrencyChangePerHour: 11.42,
						CurrencyLastUpdate:    tm,
					},
					"USD",
					nil,
				)
			},
			args: args{
				name: "USD",
			},
			want: &Currency{
				CurrencyID:            1,
				CurrencyName:          "USD",
				CurrencyPrice:         90.42,
				CurrencyMinPrice:      27.00,
				CurrencyMaxPrice:      120.00,
				CurrencyChangePerHour: 11.42,
				CurrencyLastUpdate:    tm,
			},
			wantErr: nil,
		},
		{
			name: "some error",
			setup: func() {
				setList(nil, "ETH", errors.New("currency not found"))
			},
			args: args{
				name: "ETH",
			},
			want:    nil,
			wantErr: errors.New("error in method GetCurrency: currency not found"),
		},
		{
			name: "empty list",
			setup: func() {
				setList(nil, "", errors.New("currency not found"))
			},
			args: args{
				name: "",
			},
			want:    nil,
			wantErr: errors.New("error in method GetCurrency: currency not found"),
		},
	}

	svc := NewService(repo, slog.New(slog.NewTextHandler(os.Stdout, nil)), nil)

	for _, tt := range tests {
		t.Run(tt.name, func(y *testing.T) {
			defer repo.AssertExpectations(t)

			tt.setup()

			got, err := svc.GetCurrency(context.Background(), tt.args.name)
			if err != nil && assert.Error(t, tt.wantErr, err.Error()) {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, tt.wantErr)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestInsertCurrencies(t *testing.T) {
	// InsertCurrencies(ctx context.Context, currencies []Currency) ([]Currency, error)
	repo := new(MockRepo)

	setList := func(currencies []Currency, err error) {
		repo.On("SelectAllCurrencies", mock.Anything).Return(currencies, err).Once()
	}

	const layout = "2006-01-02T15:04:05Z07:00"
	tm, _ := time.Parse(layout, "2024-04-16T12:00:00+03:00")

	tests := []struct {
		name    string
		setup   func()
		want    []Currency
		wantErr error
	}{
		{
			name: "success",
			setup: func() {
				setList(
					[]Currency{
						{
							CurrencyID:            1,
							CurrencyName:          "USD",
							CurrencyPrice:         90.42,
							CurrencyMinPrice:      27.00,
							CurrencyMaxPrice:      120.00,
							CurrencyChangePerHour: 11.42,
							CurrencyLastUpdate:    tm,
						},
					},
					nil,
				)
			},
			want: []Currency{
				{
					CurrencyID:            1,
					CurrencyName:          "USD",
					CurrencyPrice:         90.42,
					CurrencyMinPrice:      27.00,
					CurrencyMaxPrice:      120.00,
					CurrencyChangePerHour: 11.42,
					CurrencyLastUpdate:    tm,
				},
			},
			wantErr: nil,
		},
		{
			name: "some error",
			setup: func() {
				setList(nil, errors.New("no currencies"))
			},
			want:    nil,
			wantErr: errors.New("error in method GetCurrencies no currencies"),
		},
		{
			name: "empty list",
			setup: func() {
				setList(nil, nil)
			},
			want:    nil,
			wantErr: nil,
		},
	}

	svc := NewService(repo, slog.New(slog.NewTextHandler(os.Stdout, nil)), nil)

	for _, tt := range tests {
		t.Run(tt.name, func(y *testing.T) {
			defer repo.AssertExpectations(t)

			tt.setup()

			got, err := svc.GetCurrencies(context.Background())
			if err != nil && assert.Error(t, tt.wantErr, err.Error()) {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, tt.wantErr)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSelectChangesPerHour(t *testing.T) {
	// SelectChangesPerHour(ctx context.Context, curr string) (float64, error)
	repo := new(MockRepo)

	type args struct {
		name string
	}

	setList := func(name string, change float64, err error) {
		repo.On("SelectChangesPerHour", mock.Anything, name).Return(change, err)
	}

	tests := []struct {
		name    string
		setup   func()
		args    args
		want    float64
		wantErr error
	}{
		{
			name: "success",
			setup: func() {
				setList(
					"USD",
					11.42,
					nil,
				)
			},
			args: args{
				name: "USD",
			},
			want:    11.42,
			wantErr: nil,
		},
		{
			name: "some error",
			setup: func() {
				setList("ETH", -1, errors.New("currency not found"))
			},
			args: args{
				name: "ETH",
			},
			want:    -1,
			wantErr: errors.New("error in Service's method GetChangesPerHour: currency not found"),
		},
		{
			name: "empty list",
			setup: func() {
				setList("", -1, errors.New("currency not found"))
			},
			args: args{
				name: "",
			},
			want:    -1,
			wantErr: errors.New("error in Service's method GetChangesPerHour: currency not found"),
		},
	}

	svc := NewService(repo, slog.New(slog.NewTextHandler(os.Stdout, nil)), nil)

	for _, tt := range tests {
		t.Run(tt.name, func(y *testing.T) {
			defer repo.AssertExpectations(t)

			tt.setup()

			got, err := svc.GetChangesPerHour(context.Background(), tt.args.name)
			if err != nil && assert.Error(t, tt.wantErr, err.Error()) {
				assert.EqualError(t, err, tt.wantErr.Error())
			} else {
				assert.NoError(t, tt.wantErr)
			}

			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSetChangesPerHour(t *testing.T) {
	// SetChangesPerHour(ctx context.Context, currencies []Currency) error
}
