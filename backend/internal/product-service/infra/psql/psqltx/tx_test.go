package psqltx

import (
	"context"
	"database/sql"
	"errors"
	"testing"
)

func TestRunInTxUsesExistingTransactionWithoutOwningIt(t *testing.T) {
	t.Parallel()

	expected := errors.New("inner error")
	ctx := context.WithValue(context.Background(), txKey{}, (*sql.Tx)(nil))

	err := (&TxRunner{}).RunInTx(ctx, func(got context.Context) error {
		if _, ok := TxFromContext(got); !ok {
			t.Fatal("expected transaction to remain in context")
		}
		return expected
	})

	if !errors.Is(err, expected) {
		t.Fatalf("expected inner error, got %v", err)
	}
}
