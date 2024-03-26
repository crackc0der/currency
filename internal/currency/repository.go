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
		return nil, fmt.Errorf("error to select all currencies in repositery's method SelectCurrencies: %w", err)
	}

	for rows.Next() {
		var currency Currency
		err := rows.Scan(&currency.CurrencyID, &currency.CurrencyName, &currency.CurrencyPrice, &currency.CurrencyMinPrice, &currency.CurrencyMaxPrice, &currency.CurrencyPercentageChange)
		if err != nil {
			return nil, fmt.Errorf("failed to get currency in repositery's method SelectAllCurrencies, Scan: %w", err)
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
		return nil, fmt.Errorf("failed to get currency in repositery's method SelectCurrency: %w", err)
	}

	return currency, nil
}

func (r Repository) InsertCurrencies(ctx context.Context, currencies []Currency) ([]Currency, error) {
	query := `insert into currency (currency_name, price_min, price_max, changes_per_hour) 
				values (@currencyName, @priceMin, @priceMax, @changesPerHour) on conflict (currency_name) do update set
				currency_name=currencyName, price_min=priceMin, price_max=priceMax, changes_per_hour=changesPerHour`
	batch := &pgx.Batch{}

	for _, currency := range currencies {
		args := pgx.NamedArgs{
			"currenncyName":  currency.CurrencyName,
			"priceMin":       currency.CurrencyMinPrice,
			"priceMax":       currency.CurrencyMaxPrice,
			"changesPerHour": currency.CurrencyPercentageChange,
		}

		batch.Queue(query, args)
	}

	results := r.conn.SendBatch(ctx, batch)
	defer results.Close()

	for _, currency := range currencies {
		_, err := results.Exec()
		if err != nil {
			return nil, fmt.Errorf("error to insert %v in repositery's method InsertCurrencies %w", currency, err)
		}
	}

	return currencies, nil
}

func (r Repository) SelectChangesPerHour(ctx context.Context, curr string) (float64, error) {
	var currencyPerHour float64
	query := "select changes_per_hour from currency where currency_name = $1"
	err := r.conn.QueryRow(ctx, query, &currencyPerHour).Scan(&currencyPerHour)

	if err != nil {
		return -1, fmt.Errorf("failed to get changes per hour in repositery's method SelectChangesPerHour %w", err)
	}

	return currencyPerHour, nil
}
