package store

import (
	"context"

	"github.com/backersorg/synchronizer/internal/subtypes/schematypes"
)

// The Store abstracts every CRUD app
type Store interface {

	// PushSingle will override a single value on the specified field
	PushSingle(context.Context, schematypes.Value, schematypes.ID, schematypes.TableName, schematypes.ColumnName) error

	// Push will post the given commit. It returns a slice of IDs
	Push(context.Context, *schematypes.Commit) (*schematypes.Summary, error)

	// Pull retrieves the column(s) given from an entity on a table
	Pull(context.Context, schematypes.ID, schematypes.TableName, ...schematypes.ColumnName) (*schematypes.Commit, error)

	// Remove deletes a record given its id and table
	Remove(context.Context, schematypes.ID, schematypes.TableName) error

	// Init initializes every persistant needed action (such as token refreshing)
	Init(context.Context) error
}
