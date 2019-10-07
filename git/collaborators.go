package git

import "context"

// Collaborator is any agent which performs transactions
type Collaborator interface {
	Push(context.Context, *Commit) (Summary, error)

	Pull(context.Context, *Commit) (*Commit, error)

	Delete(context.Context, *Commit) (Summary, error)

	Init(context.Context) error
}
