package schema

import (
	"context"

	"github.com/jmoiron/sqlx"

	"github.com/sebach1/git-crud/integrity"
)

// ByName finds a schema from the DB given its name
func ByName(ctx context.Context, db *sqlx.DB, schName integrity.SchemaName) (*Schema, error) {
	sch := Schema{}
	err := db.GetContext(ctx, &sch, `SELECT * FROM schemas WHERE name=?`, schName)
	if err != nil {
		return nil, err
	}
	return &sch, nil
}
