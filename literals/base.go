package literals

import (
	"context"
	"errors"

	"github.com/sebach1/git-crud/git"
)

// The Base builtin literal provides a standard Collaborator agent, and can be used to prevent
// creating too much boilerplates
// Notice that all the Collaborator methods it implements returns errors
type Base struct{}

func (b *Base) Push(ctx context.Context, comm *git.Commit) (*git.Commit, error) {
	return nil, errNotImplemented
}

func (b *Base) Pull(ctx context.Context, comm *git.Commit) (*git.Commit, error) {
	return nil, errNotImplemented
}

func (b *Base) Delete(ctx context.Context, comm *git.Commit) (*git.Commit, error) {
	return nil, errNotImplemented
}

func (b *Base) Init(ctx context.Context) error {
	return errNotImplemented
}

var errNotImplemented = errors.New("not implemented")
