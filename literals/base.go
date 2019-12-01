package literals

import (
	"context"
	"errors"

	"github.com/sebach1/rtc/git"
)

// The BaseCollab builtin literal provides a standard Collaborator agent, and can be used to prevent
// creating too much boilerplates
// Notice that all the Collaborator methods it implements returns errors
type BaseCollab struct{}

// Push is a base not-implemented method to achieve the Collaborator interface
func (b *BaseCollab) Push(ctx context.Context, comm *git.Commit) (*git.Commit, error) {
	return nil, errNotImplemented
}

// Pull is a base not-implemented method to achieve the Collaborator interface
func (b *BaseCollab) Pull(ctx context.Context, comm *git.Commit) (*git.Commit, error) {
	return nil, errNotImplemented
}

// Delete is a base not-implemented method to achieve the Collaborator interface
func (b *BaseCollab) Delete(ctx context.Context, comm *git.Commit) (*git.Commit, error) {
	return nil, errNotImplemented
}

// Init is a base not-implemented method to achieve the Collaborator interface
func (b *BaseCollab) Init(ctx context.Context) error {
	return errNotImplemented
}

var errNotImplemented = errors.New("not implemented")
