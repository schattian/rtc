package schema

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/sebach1/git-crud/integrity"
)

func ByName(ctx context.Context, DB *sql.DB, schName integrity.SchemaName) (*Schema, error) {
	sch := &Schema{}
	row := DB.QueryRowContext(ctx, fmt.Sprintf("SELECT * FROM schemas WHERE Name = %s", schName))
	err := row.Scan(sch)
	if err != nil {
		return nil, err
	}
	return sch, nil
}
