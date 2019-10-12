package github

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/sebach1/git-crud/git"
	"github.com/sebach1/git-crud/internal/integrity"
	"github.com/sebach1/git-crud/literals"
	"github.com/sebach1/git-crud/schema"
	"github.com/sebach1/git-crud/valide"
)

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
				OptionKeys: []integrity.OptionKey{"username"},
			},

			&schema.Table{
				Name: "organizations",
				Columns: []*schema.Column{
					&schema.Column{Name: "name", Validator: valide.String},
					&schema.Column{Name: "projects", Validator: valide.Bytes},
				},
				OptionKeys: []integrity.OptionKey{"owner"},
			},
		},
	}

	// OpenSource is the open source code community
	OpenSource = &git.Community{
		&git.Team{
			AssignedSchema: "github",
			Members: []*git.Member{
				&git.Member{AssignedTable: "repositories", Collab: new(repositories)},
				&git.Member{AssignedTable: "organizations", Collab: new(organizations)},
			},
		},
	}
)

type repositories struct {
	literals.Base
}

type organizations struct {
	literals.Base
}

func (orgs *organizations) Push(ctx context.Context, comm *git.Commit) (*git.Commit, error) {
	commType, _ := comm.Type()

	body, err := integrity.ToJSON(comm)
	if err != nil {
		return nil, err
	}

	opts, err := comm.Options()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		commType.ToHTTPVerb(),
		orgs.URL(opts["owner"].(string)),
		bytes.NewReader(body),
	)
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

func (r *repositories) Push(ctx context.Context, comm *git.Commit) (*git.Commit, error) {
	commType, _ := comm.Type()

	body, err := integrity.ToJSON(comm)
	if err != nil {
		return nil, err
	}

	opts, err := comm.Options()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		commType.ToHTTPVerb(),
		r.URL(opts["username"].(string)),
		bytes.NewReader(body),
	)
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

func (r *repositories) URL(username string) string {
	return fmt.Sprintf("%v/user/%v/repos", baseURL, username)
}

func (orgs *organizations) URL(owner string) string {
	return fmt.Sprintf("%v/orgs/%v", baseURL, owner)
}
