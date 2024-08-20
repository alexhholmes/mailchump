package pgdb

import (
	"database/sql"
	"errors"
	"log/slog"
)

// HandleTxError is a helper function to handle transaction commit/rollback;
// it will roll back the transaction if an error is passed. In other cases
// it will do nothing.
func HandleTxError(err error, tx *sql.Tx) func() {
	return func() {
		// Ignore ErrTxDone as just means the transaction is already complete
		if err != nil && !errors.Is(err, sql.ErrTxDone) {
			err = tx.Rollback()
			if err != nil {
				slog.Error("failed to rollback transaction", "error", err)
			}
		}
	}
}
