package store

import (
	"context"
	"strings"

	"github.com/jmoiron/sqlx"
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

func UpdateIntoDB(ctx context.Context, db *sqlx.DB, storable Storable) error {
	if ctx == nil {
		ctx = context.Background()
	}
	_, err := db.NamedExecContext(ctx, execBoilerplate("UPDATE", storable), storable)
	return err
}

func UpdateBatchIntoDB(ctx context.Context, db *sqlx.DB, storables ...Storable) error {
	if ctx == nil {
		ctx = context.Background()
	}
	ref := storables[0]
	_, err := db.NamedExecContext(ctx, execBoilerplate("UPDATE", ref), storables)
	return err
}

// InsertIntoDB inserts the storable entity to the DB
// Finally, it assigns the inserted Id to the given entity
func InsertIntoDB(ctx context.Context, db *sqlx.DB, storable Storable) error {
	if ctx == nil {
		ctx = context.Background()
	}
	res, err := db.NamedExecContext(ctx, execBoilerplate("INSERT INTO", storable), storable)
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
	return err
}

// func containsColumn(storable Storable, columnName string) bool{
// 	for _, col := range storable.SQLColumns(){
// 		if col == columnName{
// 			return true
// 		}
// 	}
// 			return false
// }

func execBoilerplate(action string, storable Storable) string {
	return action + ` ` + storable.SQLTable() + ` ` + sqlColumnNames(storable) + ` VALUES ` + sqlColumnValues(storable)
}
