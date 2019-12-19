package git

import (
	"testing"

	"github.com/mitchellh/copystructure"
)

func (b *Branch) copy(t *testing.T) *Branch {
	t.Helper()
	new, err := copystructure.Copy(b)
	if err != nil {
		t.Fatalf("could not be able to copy branch: %v", err)
	}
	return new.(*Branch)
}

func (chg *Change) copy(t *testing.T) *Change {
	t.Helper()
	new, err := copystructure.Copy(chg)
	if err != nil {
		t.Fatalf("could not be able to copy change: %v", err)
	}
	return new.(*Change)
}

func (pR *PullRequest) copy(t *testing.T) *PullRequest {
	t.Helper()
	new, err := copystructure.Copy(pR)
	if err != nil {
		t.Fatalf("could not be able to copy pull request: %v", err)
	}
	return new.(*PullRequest)
}

func (t *Team) copy(T *testing.T) *Team {
	T.Helper()
	new, err := copystructure.Copy(t)
	if err != nil {
		T.Fatalf("could not be able to copy pull request: %v", err)
	}
	return new.(*Team)
}
