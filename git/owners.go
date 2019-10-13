package git

import (
	"context"
	"sync"

	"github.com/sebach1/git-crud/schema"

	"github.com/pkg/errors"
	"github.com/sebach1/git-crud/internal/integrity"
)

// Owner is the agent which coordinates any given action
// Notice that an Owner is a Collaborator
type Owner struct {
	Project *schema.Planisphere

	wg      *sync.WaitGroup
	Summary chan *Result
}

func (own *Owner) ReviewPRCommit(sch *schema.Schema, pR *PullRequest, commIdx int) {
	defer own.wg.Done()
	comm := pR.Commits[commIdx]
	defer func() {
		if comm.Reviewer == nil {
			comm.Errored = true
		}
	}()

	schErrCh := make(chan error, len(comm.Changes))
	for _, chg := range comm.Changes {
		own.wg.Add(1)
		go sch.Validate(chg.TableName, chg.ColumnName, chg.Options.Keys(), chg.Value(), own.Project, own.wg, schErrCh)
	}

	tableName, err := comm.TableName()
	if err != nil {
		own.Summary <- &Result{CommitID: comm.ID, Error: errors.Wrap(err, "revieweing pr")}
		return
	}

	_, err = comm.Options()
	if err != nil {
		own.Summary <- &Result{CommitID: comm.ID, Error: errors.Wrap(err, "revieweing pr")}
		return
	}

	if len(schErrCh) > 0 {
		var errs string
		for err := range schErrCh {
			errs += err.Error()
			errs += "; "
		}
		own.Summary <- &Result{CommitID: comm.ID, Error: errors.Wrap(err, "revieweing pr")}
		return
	}

	reviewer, err := pR.Team.Delegate(tableName)
	if err != nil {
		own.Summary <- &Result{CommitID: comm.ID, Error: errors.Wrap(err, "revieweing pr")}
		return
	}
	comm.Reviewer = reviewer
}

// Orchestrate sends the order to all the collaborators available to execute
// the needed actions in order to achieve the commitment, creating a new PullRequest
func (own *Owner) Orchestrate(
	ctx context.Context,
	community *Community,
	schName integrity.SchemaName,
	comm *Commit,
	strategy changesMatcher,
) error {
	sch, err := own.Project.GetSchemaFromName(schName)
	if err != nil {
		return err
	}

	own.wg = new(sync.WaitGroup)
	var pR PullRequest

	for _, changes := range comm.GroupBy(strategy) { // Splits incompatibilities onto the pR
		comm := &Commit{Changes: changes}
		pR.Commits = append(pR.Commits, comm)
	}
	pR.AssignTeam(community, schName)

	own.Summary = make(chan *Result, len(pR.Commits))

	own.wg.Add(len(pR.Commits))
	for commIdx := range pR.Commits {
		go own.ReviewPRCommit(sch, &pR, commIdx)
	}
	own.wg.Wait()

	own.wg.Add(1)
	go own.Merge(ctx, &pR)
	own.wg.Wait()

	return nil
}

// Merge performs the needed actions in order to merge the pullRequest
func (own *Owner) Merge(ctx context.Context, pR *PullRequest) {
	for _, comm := range pR.Commits {
		if comm.Errored {
			continue
		}

		commType, err := comm.Type()
		if err != nil {
			own.Summary <- &Result{CommitID: comm.ID, Error: errors.Wrap(err, "merging")}
			continue
		}

		switch commType {
		case "create", "update":
			own.wg.Add(1)
			go own.Push(ctx, comm)
		case "retrieve":
			own.wg.Add(1)
			go own.Pull(ctx, comm)
		case "delete":
			own.wg.Add(1)
			go own.Delete(ctx, comm)
		}
	}
}

// Push will orchestrate the pushes of any collaborator
func (own *Owner) Push(ctx context.Context, comm *Commit) (*Commit, error) {
	defer own.wg.Done()
	newComm := new(Commit)
	*newComm = *comm
	err := comm.Reviewer.Init(ctx)
	if err != nil {
		own.Summary <- &Result{CommitID: comm.ID, Error: errors.Wrap(err, "pushing from owner")}
		return comm, err
	}
	newComm, err = comm.Reviewer.Push(ctx, newComm)
	if err != nil {
		own.Summary <- &Result{CommitID: comm.ID, Error: errors.Wrap(err, "pushing from owner")}
		return comm, err
	}
	*comm = *newComm
	return comm, nil
}

// Pull will orchestrate the pulls of any collaborator
func (own *Owner) Pull(ctx context.Context, comm *Commit) (*Commit, error) {
	defer own.wg.Done()
	newComm := new(Commit)
	*newComm = *comm
	err := comm.Reviewer.Init(ctx)
	if err != nil {
		own.Summary <- &Result{CommitID: comm.ID, Error: errors.Wrap(err, "pushing from owner")}
		return comm, err
	}
	newComm, err = comm.Reviewer.Pull(ctx, newComm)
	if err != nil {
		own.Summary <- &Result{CommitID: comm.ID, Error: errors.Wrap(err, "pulling from owner")}
		return comm, err
	}
	*comm = *newComm
	return comm, nil
}

// Delete will orchestrate the deletions of any collaborator
func (own *Owner) Delete(ctx context.Context, comm *Commit) (*Commit, error) {
	defer own.wg.Done()
	newComm := new(Commit)
	*newComm = *comm
	err := comm.Reviewer.Init(ctx)
	if err != nil {
		own.Summary <- &Result{CommitID: comm.ID, Error: errors.Wrap(err, "pushing from owner")}
		return comm, err
	}
	newComm, err = comm.Reviewer.Delete(ctx, newComm)
	if err != nil {
		own.Summary <- &Result{CommitID: comm.ID, Error: errors.Wrap(err, "deleting from owner")}
		return comm, err
	}
	*comm = *newComm
	return comm, nil
}

// func (own *Owner) IsAlreadySummarized(commID int) bool {
// 	for summ := range own.Summary {
// 		if summ.CommitID == commID {
// 			return true
// 		}
// 	}
// 	return false
// }
