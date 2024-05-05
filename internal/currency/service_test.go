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
	"github.com/stretchr/testify/require"
)

var (
	ErrNoCurrencies              = errors.New("currency not found")
	ErrGetCorrenciesNoCurrencies = errors.New("error in method GetCurrencies currency not found")
	ErrGetCorrencyNoCorrency     = errors.New("error in method GetCurrency: currency not found")
	ErrGetChangesPerHour         = errors.New("error in Service's method GetChangesPerHour: currency not found")
)

func TestSelectAllCurrencies(t *testing.T) {
	// SelectAllCurrencies(ctx context.Context) ([]Currency, error)
	t.Parallel()

	repo := new(MockRepo)

	setList := func(currencies []Currency, err error) {
		repo.On("SelectAllCurrencies", mock.Anything).Return(currencies, err).Once()
	}

	const layout = "2006-01-02T15:04:05Z07:00"
	time, _ := time.Parse(layout, "2024-04-16T12:00:00+03:00")

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
							CurrencyLastUpdate:    time,
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
					CurrencyLastUpdate:    time,
				},
			},
			wantErr: nil,
		},
		{
			name: "some error",
			setup: func() {
				setList(nil, ErrNoCurrencies)
			},
			want:    nil,
			wantErr: ErrGetCorrenciesNoCurrencies,
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

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			defer repo.AssertExpectations(t)

			testCase.setup()

			got, err := svc.GetCurrencies(context.Background())
			if err != nil {
				require.EqualError(t, err, testCase.wantErr.Error())
			} else {
				require.NoError(t, testCase.wantErr)
			}

			assert.Equal(t, testCase.want, got)
		})
	}
}

func TestSelectCurrency(t *testing.T) {
	// SelectCurrency(ctx context.Context, name string) (*Currency, error)
	t.Parallel()

	repo := new(MockRepo)

	type args struct {
		name string
	}

	setList := func(currency *Currency, name string, err error) {
		repo.On("SelectCurrency", mock.Anything, name).Return(currency, err)
	}

	const layout = "2006-01-02T15:04:05Z07:00"
	time, _ := time.Parse(layout, "2024-04-16T12:00:00+03:00")

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
						CurrencyLastUpdate:    time,
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
				CurrencyLastUpdate:    time,
			},
			wantErr: nil,
		},
		{
			name: "some error",
			setup: func() {
				setList(nil, "ETH", ErrNoCurrencies)
			},
			args: args{
				name: "ETH",
			},
			want:    nil,
			wantErr: ErrGetCorrencyNoCorrency,
		},
		{
			name: "empty list",
			setup: func() {
				setList(nil, "", ErrNoCurrencies)
			},
			args: args{
				name: "",
			},
			want:    nil,
			wantErr: ErrGetCorrencyNoCorrency,
		},
	}

	svc := NewService(repo, slog.New(slog.NewTextHandler(os.Stdout, nil)), nil)

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			defer repo.AssertExpectations(t)

			testCase.setup()

			got, err := svc.GetCurrency(context.Background(), testCase.args.name)
			if err != nil {
				require.EqualError(t, err, testCase.wantErr.Error())
			} else {
				require.NoError(t, testCase.wantErr)
			}

			assert.Equal(t, testCase.want, got)
		})
	}
}

func TestSelectChangesPerHour(t *testing.T) {
	// SelectChangesPerHour(ctx context.Context, curr string) (float64, error)
	t.Parallel()

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
				setList("ETH", -1, ErrNoCurrencies)
			},
			args: args{
				name: "ETH",
			},
			want:    -1,
			wantErr: ErrGetChangesPerHour,
		},
		{
			name: "empty list",
			setup: func() {
				setList("", -1, ErrNoCurrencies)
			},
			args: args{
				name: "",
			},
			want:    -1,
			wantErr: ErrGetChangesPerHour,
		},
	}

	svc := NewService(repo, slog.New(slog.NewTextHandler(os.Stdout, nil)), nil)

	for _, testCase := range tests {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			defer repo.AssertExpectations(t)

			testCase.setup()

			got, err := svc.GetChangesPerHour(context.Background(), testCase.args.name)
			if err != nil {
				require.EqualError(t, err, testCase.wantErr.Error())
			} else {
				require.NoError(t, testCase.wantErr)
			}

			assert.InEpsilon(t, testCase.want, got, 0.0001)
		})
	}
}

// func TestSetChangesPerHour(t *testing.T) {
// 	// SetChangesPerHour(ctx context.Context, currencies []Currency) error
// }
