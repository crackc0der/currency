package currency

import (
	"context"
	"testing"

	"time"

	"github.com/chrisyxlee/pgxpoolmock"
	"github.com/golang/mock/gomock"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx"
)

func TestRepositorySelectAllCurrencies(t *testing.T) {
	t.Parallel()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	type PgxPoolInterface interface {
		// Include all methods you need from pgxpool.Pool
		Exec(ctx context.Context, sql string, arguments ...interface{}) (pgconn.CommandTag, error)
		QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
		// Add other methods as needed
	}

	mockPool := pgxpoolmock.NewMockPgxIface(ctrl)
	columns := []string{"id", "currency_name", "price", "price_min", "price_max", "changes_per_hour", "last_update"}
	const layout = "2006-01-02T15:04:05Z07:00"
	tm, _ := time.Parse(layout, "2024-04-16T12:00:00+03:00")
	pgxRows := pgxpoolmock.NewRows(columns).AddRow(1, "USD", 100.42, 36.76, 165.32, 23.54, tm).ToPgxRows()
	mockPool.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Return(pgxRows, nil)

	repository := Repository{mockPool}
	if err != nil {
		t.Fatal(err)
	}

	currencies, err := repository.SelectAllCurrencies(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}
