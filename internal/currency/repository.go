package currency

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func NewRepository(dsn string) (*Repository, error) {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to create new repository in method NewRepository: %w", err)
	}
	return &Repository{conn: conn}, nil
}

type Repository struct {
	conn *pgx.Conn
}

func (r Repository) SelectAllCurrencies(ctx context.Context) ([]Currency, error) {
	var currencies []Currency
	query := "select * from currency"
	rows, err := r.conn.Query(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("error to select all currencies in method repositery's method SelectCurrencies: %w", err)
	}

	for rows.Next() {
		var currency Currency
		err := rows.Scan(&currency.CurrencyID, &currency.CurrencyName, &currency.CurrencyPrice, &currency.CurrencyMinPrice, &currency.CurrencyMaxPrice, &currency.CurrencyPercentageChange)
		if err != nil {
			continue
		}

		currencies = append(currencies, currency)
	}

	return currencies, nil
}

func (r Repository) SelectCurrency(ctx context.Context, name string) (*Currency, error) {
	var currency *Currency
	query := "select * from currency where currency_name = $1"
	err := r.conn.QueryRow(ctx, query, name).Scan(&currency.CurrencyID, &currency.CurrencyName, &currency.CurrencyPrice, &currency.CurrencyMinPrice, &currency.CurrencyMaxPrice, &currency.CurrencyPercentageChange)

	if err != nil {
		return nil, fmt.Errorf("failed to get currency in method SelectCurrency: %w", err)
	}

	return currency, nil
}

func (r Repository) InsertCurrencies(currencies []Currency) bool {
}

func (r Repository) InsertCurrency(currency Currency) bool {
}

func (r Repository) UpdateCurrencies(currency Currency) bool {
}

func (r Repository) SelectChangesPerHour(curr string) float64 {
}
