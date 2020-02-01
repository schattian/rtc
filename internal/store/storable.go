package store

import (
	"context"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/sebach1/rtc/internal/name"
)

// Storable entity is any entity that can be stored in an SQL database
type Storable interface {
	GetId() int64
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

func UpsertIntoDB(ctx context.Context, db *sqlx.DB, storables ...Storable) error {
	var inserts, updates []Storable
	for _, store := range storables {
		if store.GetId() == 0 {
			inserts = append(inserts, store)
		} else {
			updates = append(updates, store)
		}
	}
	err := UpdateIntoDB(ctx, db, updates...)
	if err != nil {
		return errors.Wrap(err, "update into db")
	}
	err = InsertIntoDB(ctx, db, inserts...)
	if err != nil {
		return errors.Wrap(err, "insert into db")
	}
	return nil
}

func UpdateIntoDB(ctx context.Context, db *sqlx.DB, storables ...Storable) error {
	qtToStore := len(storables)
	if qtToStore == 0 {
		return nil
	}

	ref := storables[0] // takes it as a reference for all entities given
	qr := execBoilerplate("UPDATE", ref)
	rows, err := db.NamedExecContext(ctx, qr, storables)
	if err != nil {
		return errors.Wrap(err, "named exec ctx")
	}
	rowsQt, err := rows.RowsAffected()
	if err != nil {
		return errors.Wrap(err, "rows affected")
	}
	if int(rowsQt) != qtToStore {
		return errMismatchAffectedRows
	}
	return nil
}

var errNilStorableEntity = errors.New("nil storable entity")
var errMismatchAffectedRows = errors.New("the affected rows quantity does not match with the given storables")

// InsertIntoDB inserts the storable entity to the DB
// Finally, it assigns the inserted Id to the given entities
func InsertIntoDB(ctx context.Context, db *sqlx.DB, storables ...Storable) error {
	if len(storables) == 0 {
		return errNilStorableEntity
	}
	ref := storables[0] // takes it as a reference for all entities given
	qr := execBoilerplate("INSERT INTO", ref) + " RETURNING id"
	ids, err := db.NamedQueryContext(ctx, qr, storables)
	if err != nil {
		return errors.Wrap(err, "named query ctx")
	}
	defer ids.Close()

	var i int
	for ids.Next() {
		var id int64
		err := ids.Scan(&id)
		if err != nil {
			return errors.Wrap(err, "id scan")
		}
		storables[i].SetId(id)
		i += 1
	}
	err = ids.Err()
	if err != nil {
		return errors.Wrap(err, "cursor err")
	}

	return nil
}

func DeleteFromDB(ctx context.Context, db *sqlx.DB, storable Storable) error {
	if storable.GetId() == 0 {
		return nil
	}

	if ctx == nil {
		ctx = context.Background()
	}
	_, err := db.NamedExecContext(
		ctx,
		`DELETE FROM `+storable.SQLTable()+` WHERE id=:id`,
		storable,
	)
	if err != nil {
		return errors.Wrap(err, "named exec ctx")
	}
	return nil
}

func execBoilerplate(action string, storable Storable) string {
	return action + ` ` + storable.SQLTable() + ` ` + sqlColumnNames(storable) + ` VALUES ` + sqlColumnValues(storable)
}
