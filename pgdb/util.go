package pgdb

import (
	"database/sql"
	"errors"
	"log"
	"log/slog"
	"os"
	"reflect"
	"strings"
)

// TODO write a unit tests and runtime check during init() to ensure that the
// `model` pkg structs all adhere to to the `tables.sql` schema.

// scan is an interface for the Scan method of sql.Rows
// and sql.Row.
type scan interface {
	Scan(dest ...any) error
}

// MapStruct uses reflection to map the columns of a sql.Rows to
// a generic struct.
func MapStruct[T any, R scan](row R) (T, error) {
	var fields T

	if &row == nil {
		return fields, sql.ErrNoRows
	}

	val := reflect.ValueOf(&fields).Elem()
	numCols := val.NumField()
	columns := make([]interface{}, numCols)
	for i := 0; i < numCols; i++ {
		columns[i] = val.
			Field(i).
			Addr().
			Interface()
	}

	err := row.Scan(columns...)
	if err != nil {
		log.Fatal(err)
	}

	return fields, nil
}

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

// HandleCloseResult is a helper function for the `defer` to handle closing
// sql.Rows; log and ignore any error that occurs.
func HandleCloseResult(res *sql.Rows) func() {
	return func() {
		if err := res.Close(); err != nil {
			slog.Warn("failed to close rows", "error", err)
		}
	}
}

// InitializeLocalDB initializes the pgdb tables and adds test data. This will log fatal
// any error that occurs.
func InitializeLocalDB() {
	env := strings.ToLower(os.Getenv("ENVIRONMENT"))
	if env != "local" {
		log.Fatalf("Cannot initialize database in non-local environment")
	}
	if os.Getenv("MIGRATIONS") == "" {
		log.Fatalf("MIGRATIONS environment variable not set")
	}
	if os.Getenv("MIGRATIONS_DIR") == "" {
		log.Fatalf("MIGRATIONS_DIR environment variable not set")
	}

	db, err := Init()
	if err != nil {
		log.Fatalf("Could not make db connection for dev environment initialization: %s", err.Error())
	}
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	// Order matters in the MIGRATIONS environment variable, as some sql files are
	// dependent on others.
	migrations := strings.Split(os.Getenv("MIGRATIONS"), ",")
	var files []string
	for _, f := range migrations {
		f = strings.TrimSpace(f)
		file, err := os.ReadFile(os.Getenv("MIGRATIONS_DIR") + "/" + f)
		if err != nil {
			log.Fatalf("Could not read migration file: %s", err.Error())
		}
		files = append(files, string(file))
	}

	// Execute the SQL migrations to postgres
	for _, query := range files {
		// Transaction to migrate entire tables.sql files the fresh database
		tx, err := db.Begin()
		if err != nil {
			log.Fatalf("Could not start migration transaction: %s", err.Error())
		}

		// Execute the tables.sql files statements to create the empty tables
		_, err = tx.Exec(query)
		if err != nil {
			roll := tx.Rollback()
			if roll != nil {
				log.Fatalf("Could not rollback migration transaction: %s", err.Error())
			}
			log.Fatalf("Could not execute migration statement: %s", err.Error())
		}

		// Commit the transaction
		err = tx.Commit()
		if err != nil {
			log.Fatalf("Could not commit migration transaction: %s", err.Error())
		}
	}
}
