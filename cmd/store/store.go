package store

import (
	"context"

	"github.com/sebach1/git-crud/git"
	"github.com/sebach1/git-crud/schema"
)

// The Store abstracts every CRUD app
type Store interface {

	// PushSingle will override a single value on the specified field
	// Its a shortcut for
	PushSingle(context.Context, schema.Value, schema.ID, schema.TableName, schema.ColumnName) error

	// Push will post the given commit. It returns a slice of IDs
	Push(context.Context, *git.Commit) (*git.Summary, error)

	// Pull retrieves the column(s) given from an entity on a table
	Pull(context.Context, schema.ID, schema.TableName, ...schema.ColumnName) (*git.Commit, error)

	// Remove deletes a record given its id and table
	Remove(context.Context, schema.ID, schema.TableName) error

	// Init initializes every persistant needed action (such as token refreshing)
	Init(context.Context) error
}
