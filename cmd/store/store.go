package store

import (
	"context"

	"github.com/backersorg/synchronizer/internal/subtypes/maptypes"
)

// The Store abstracts every CRUD app
type Store interface {

	// PushSingle will override a single value on the specified field
	PushSingle(context.Context, maptypes.Value, maptypes.ID, maptypes.TableName, maptypes.ColumnName) error

	// Push will post the given commit. It returns a slice of IDs
	Push(context.Context, *maptypes.Commit) (*maptypes.Summary, error)

	// Pull retrieves the column(s) given from an entity on a table
	Pull(context.Context, maptypes.ID, maptypes.TableName, ...maptypes.ColumnName) (*maptypes.Commit, error)

	// Remove deletes a record given its id and table
	Remove(context.Context, maptypes.ID, maptypes.TableName) error

	// Init initializes every persistant needed action (such as token refreshing)
	Init(context.Context) error
}
