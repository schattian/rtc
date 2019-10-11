package literals

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/sebach1/git-crud/git"
	"github.com/sebach1/git-crud/internal/integrity"
	"github.com/sebach1/git-crud/schema"
	"github.com/sebach1/git-crud/valide"
)

//! TODO: Add aggregations (to sort, pre-fetch, etc)

const baseURL = "https://api.github.com"

var (
	// GitHub is the hub of git
	GitHub = &schema.Schema{
		Name: "github",
		Blueprint: []*schema.Table{
			&schema.Table{
				Name: "repositories",
				Columns: []*schema.Column{
					&schema.Column{Name: "name", Validator: valide.String},
					&schema.Column{Name: "private", Validator: valide.String},
				},
			},
		},
	}

	// OpenSource is the open source code community
	OpenSource = &git.Community{
		&git.Team{
			AssignedSchema: "github",
			Members: []*git.Member{
				&git.Member{AssignedTable: "repositories", Collab: new(repositories)}, //? Take all the members from the same
			},
		},
	}
)

type repositories struct{}

func (r *repositories) Push(ctx context.Context, comm *git.Commit) (*git.Commit, error) {
	commType, _ := comm.Type()

	body, err := integrity.ToJSON(comm)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, commType.ToHTTPVerb(), r.URL(""), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	commit, err := git.CommitFrom(res.Body)
	if err != nil {
		return nil, err
	}
	return commit, nil
}

func (r *repositories) Pull(ctx context.Context, comm *git.Commit) (*git.Commit, error) {
	return nil, nil
}

func (r *repositories) Delete(ctx context.Context, comm *git.Commit) (*git.Commit, error) {
	return nil, nil
}

func (r *repositories) Init(ctx context.Context) error {
	return nil
}

func (r *repositories) URL(username string) string {
	return fmt.Sprintf("%v/user/%v/repos", baseURL, username)
}
