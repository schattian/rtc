package store

import (
	"context"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/sebach1/git-crud/internal/name"
)

type Storable interface {
	SetID(int64)
	Table() string
	Columns() []string
}

func sqlColumnValues(storable Storable) string {
	fValues := []string{}
	for _, s := range storable.Columns() {
		fValues = append(fValues, ":"+s)
	}
	return name.Parenthize(strings.Join(fValues, ","))
}

func sqlColumnNames(storable Storable) string {
	return name.Parenthize(strings.Join(storable.Columns(), ","))
}

func SaveToDB(storable Storable, ctx context.Context, db *sqlx.DB) error {
	res, err := db.NamedExecContext(
		ctx,
		`INSERT INTO`+storable.Table()+` `+sqlColumnNames(storable)+` VALUES `+sqlColumnValues(storable),
		storable,
	)
	if err != nil {
		return err
	}
	possibleID, err := res.LastInsertId()
	if err != nil {
		return err
	}
	storable.SetID(possibleID)
	return nil

}
