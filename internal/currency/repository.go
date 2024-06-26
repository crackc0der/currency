package currency

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRepository(dsn string) (*Repository, error) {
	conn, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		return nil, fmt.Errorf("error in Repository's method NewRepository: %w", err)
	}

	return &Repository{conn: conn}, nil
}

type Repository struct {
	conn *pgxpool.Pool
}

func (r Repository) SelectAllCurrencies(ctx context.Context) ([]Currency, error) {
	var currencies []Currency

	query := "select * from currency"

	rows, err := r.conn.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("error in Repository's method SelectCurrencies: %w", err)
	}

	for rows.Next() {
		var currency Currency

		err := rows.Scan(&currency.CurrencyID, &currency.CurrencyName, &currency.CurrencyPrice, &currency.CurrencyMinPrice,
			&currency.CurrencyMaxPrice, &currency.CurrencyChangePerHour, &currency.CurrencyLastUpdate)
		if err != nil {
			return nil, fmt.Errorf("error in Repository's method SelectAllCurrensies: %w", err)
		}

		currencies = append(currencies, currency)
	}

	return currencies, nil
}

func (r Repository) SelectCurrency(ctx context.Context, name string) (*Currency, error) {
	var currency Currency

	query := "select * from currency where currency_name = $1"

	err := r.conn.QueryRow(ctx, query, name).Scan(&currency.CurrencyID, &currency.CurrencyName, &currency.CurrencyPrice,
		&currency.CurrencyMinPrice, &currency.CurrencyMaxPrice, &currency.CurrencyChangePerHour,
		&currency.CurrencyLastUpdate)
	if err != nil {
		return nil, fmt.Errorf("error in Repository's method SelectCurrency: %w", err)
	}

	return &currency, nil
}

func (r Repository) InsertCurrencies(ctx context.Context, currencies []Currency) ([]Currency, error) {
	query := `insert into currency (currency_name, price, price_min, price_max, changes_per_hour) 
				values (@currencyName, @price, @priceMin, @priceMax, @changesPerHour) on conflict (currency_name) do update set
				currency_name=@currencyName, price=@price, price_min=@priceMin, price_max=@priceMax, 
				changes_per_hour=@changesPerHour, last_update=now()`
	batch := &pgx.Batch{}

	for _, currency := range currencies {
		args := pgx.NamedArgs{
			"currencyName":   currency.CurrencyName,
			"price":          currency.CurrencyPrice,
			"priceMin":       currency.CurrencyMinPrice,
			"priceMax":       currency.CurrencyMaxPrice,
			"changesPerHour": currency.CurrencyChangePerHour,
		}

		batch.Queue(query, args)
	}

	results := r.conn.SendBatch(ctx, batch)
	defer results.Close()

	for _, currency := range currencies {
		_, err := results.Exec()
		if err != nil {
			return nil, fmt.Errorf("error to add %s in Repository's method InsertCurrencies %w", currency.CurrencyName, err)
		}
	}

	return currencies, nil
}

func (r Repository) SelectChangesPerHour(ctx context.Context, curr string) (float64, error) {
	var currencyPerHour float64

	query := "select changes_per_hour from currency where currency_name = $1"

	err := r.conn.QueryRow(ctx, query, curr).Scan(&currencyPerHour)
	if err != nil {
		return -1, fmt.Errorf("error in Repository's method SelectChangesPerHour %w", err)
	}

	return currencyPerHour, nil
}

func (r Repository) SetChangesPerHour(ctx context.Context, currencies []Currency) error {
	query := `update currency set changes_per_hour=@changesPerHour where currency_name=@currencyName`

	batch := &pgx.Batch{}

	for _, currency := range currencies {
		args := pgx.NamedArgs{
			"currencyName":   currency.CurrencyName,
			"changesPerHour": currency.CurrencyChangePerHour,
		}
		batch.Queue(query, args)
	}

	results := r.conn.SendBatch(ctx, batch)
	defer results.Close()

	for _, currency := range currencies {
		_, err := results.Exec()
		if err != nil {
			return fmt.Errorf("error to add %s in Repository's method SetChangesPerHour: %w", currency.CurrencyName, err)
		}
	}

	return nil
}
