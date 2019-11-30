package thelper

import (
	"database/sql/driver"
	"testing"

	"github.com/sebach1/rtc/internal/store"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/sebach1/rtc/internal/name"
)

func MockDB(t *testing.T) (db *sqlx.DB, mock sqlmock.Sqlmock) {
	t.Helper()
	mockDB, mock, err := sqlmock.New(sqlmock.ValueConverterOption(driver.DefaultParameterConverter))
	if err != nil {
		t.Fatalf("could connect to database: %v", err)
	}
	db = sqlx.NewDb(mockDB, "sqlmock")
	db.MapperFunc(name.ToSnakeCase)
	return db, mock
}

func RowsFor(entity store.Storable) *sqlmock.Rows {
	return sqlmock.NewRows(entity.SQLColumns())
}
