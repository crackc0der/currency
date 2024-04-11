package currency

import "testing"

func TestSelectAllCurrencies(t *testing.T) {
	// SelectAllCurrencies(ctx context.Context) ([]Currency, error)
}

func TestSelectCurrency(t *testing.T) {
	// SelectCurrency(ctx context.Context, name string) (*Currency, error)
}

func TestInsertCurrency(t *testing.T) {
	// InsertCurrencies(ctx context.Context, currencies []Currency) ([]Currency, error)
}

func TestSelectChangesPerHour(t *testing.T) {
	// SelectChangesPerHour(ctx context.Context, curr string) (float64, error)
}

func TestSetChangesPerHour(t *testing.T) {
	// SetChangesPerHour(ctx context.Context, currencies []Currency) error
}
