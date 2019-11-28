package store

import (
	"context"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sebach1/rtc/internal/name"
)

type Storable interface {
	SetId(int64)
	SQLTable() string
	SQLColumns() []string
}

func sqlColumnValues(storable Storable) string {
	fValues := []string{}
	for _, s := range storable.SQLColumns() {
		fValues = append(fValues, ":"+s)
	}
	return name.Parenthize(strings.Join(fValues, ","))
}

func sqlColumnNames(storable Storable) string {
	return name.Parenthize(strings.Join(storable.SQLColumns(), ","))
}

// InsertToDB inserts the storable entity to the DB
// Returns the inserted id
func InsertToDB(ctx context.Context, storable Storable, db *sqlx.DB) (int64, error) {
	res, err := db.NamedExecContext(
		ctx,
		`INSERT INTO`+storable.SQLTable()+` `+sqlColumnNames(storable)+` VALUES `+sqlColumnValues(storable),
		storable,
	)
	if err != nil {
		return 0, err
	}
	return res.LastInsertId()
}
