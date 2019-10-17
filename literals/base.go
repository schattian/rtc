package literals

import (
	"context"
	"errors"

	"github.com/sebach1/git-crud/git"
)

// The BaseCollab builtin literal provides a standard Collaborator agent, and can be used to prevent
// creating too much boilerplates
// Notice that all the Collaborator methods it implements returns errors
type BaseCollab struct{}

func (b *BaseCollab) Push(ctx context.Context, comm *git.Commit) (*git.Commit, error) {
	return nil, errNotImplemented
}

func (b *BaseCollab) Pull(ctx context.Context, comm *git.Commit) (*git.Commit, error) {
	return nil, errNotImplemented
}

func (b *BaseCollab) Delete(ctx context.Context, comm *git.Commit) (*git.Commit, error) {
	return nil, errNotImplemented
}

func (b *BaseCollab) Init(ctx context.Context) error {
	return errNotImplemented
}

var errNotImplemented = errors.New("not implemented")
