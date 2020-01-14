package git

import "context"

// Collaborator is any agent which performs transactions
type Collaborator interface {
	Init(context.Context) error

	Create(context.Context, *Commit) (*Commit, error)

	Retrieve(context.Context, *Commit) (*Commit, error)

	Update(context.Context, *Commit) (*Commit, error)

	Delete(context.Context, *Commit) (*Commit, error)
}
