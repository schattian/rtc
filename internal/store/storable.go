package store

import (
	"context"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sebach1/git-crud/internal/name"
)

type Storable interface {
	setId(int64)
	table() string
	columns() []string
}

func sqlColumnValues(storable Storable) string {
	fValues := []string{}
	for _, s := range storable.columns() {
		fValues = append(fValues, ":"+s)
	}
	return name.Parenthize(strings.Join(fValues, ","))
}

func sqlColumnNames(storable Storable) string {
	return name.Parenthize(strings.Join(storable.columns(), ","))
}

// InsertToDB inserts the storable entity to the DB
func InsertToDB(ctx context.Context, storable Storable, db *sqlx.DB) error {
	res, err := db.NamedExecContext(
		ctx,
		`INSERT INTO`+storable.table()+` `+sqlColumnNames(storable)+` VALUES `+sqlColumnValues(storable),
		storable,
	)
	if err != nil {
		return err
	}
	possibleId, err := res.LastInsertId()
	if err != nil {
		return err
	}
	storable.setId(possibleId)
	return nil
}
