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
// Finally, it assigns the inserted Id to the given entity
func InsertToDB(ctx context.Context, db *sqlx.DB, storable Storable) error {
	if ctx == nil {
		ctx = context.Background()
	}
	res, err := db.NamedExecContext(
		ctx,
		`INSERT INTO `+storable.SQLTable()+` `+sqlColumnNames(storable)+` VALUES `+sqlColumnValues(storable),
		storable,
	)
	if err != nil {
		return err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}
	storable.SetId(id)
	return nil
}
