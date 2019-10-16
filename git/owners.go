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
	Summary chan *Result

	wg *sync.WaitGroup
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
	err := own.Validate()
	if err != nil {
		return err
	}

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
	close(own.Summary)
	return nil
}

// Merge performs the needed actions in order to merge the pullRequest
func (own *Owner) Merge(ctx context.Context, pR *PullRequest) {
	defer own.wg.Done()
	for _, comm := range pR.Commits {
		if comm.Errored {
			continue // Skips validation errs
		}

		commType, err := comm.Type()
		if err != nil {
			own.Summary <- &Result{CommitID: comm.ID, Error: errors.Wrap(err, "merging")}
			continue
		}

		own.wg.Add(1)
		switch commType {
		case "create", "update":
			go own.Push(ctx, comm)
		case "retrieve":
			go own.Pull(ctx, comm)
		case "delete":
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

// Validate validates itself integrity to be able to perform orchestration & reviewing (owner)
func (own *Owner) Validate() error {
	if own.Project == nil {
		return errNilProject
	}
	if len(*own.Project) == 0 {
		return errEmptyProject
	}
	return nil
}

// ReviewPRCommit wraps schema validations to a specified commit of the given PullRequest
func (own *Owner) ReviewPRCommit(sch *schema.Schema, pR *PullRequest, commIdx int) {
	var err error
	defer own.wg.Done()
	var reviewWg sync.WaitGroup

	comm := pR.Commits[commIdx]
	defer func() {
		if err != nil {
			own.Summary <- &Result{CommitID: comm.ID, Error: errors.Wrap(err, "reviewing merge")}
			comm.Errored = true
		}
	}()

	schErrCh := make(chan error, len(comm.Changes))
	reviewWg.Add(len(comm.Changes))
	for _, chg := range comm.Changes {
		err = chg.Validate() // Performed sync to be strictly before any type assertion of the entire commit
		if err != nil {
			return
		}
		go sch.Validate(chg.TableName, chg.ColumnName, chg.Options.Keys(), chg.Value(), own.Project, &reviewWg, schErrCh)
	}

	tableName, err := comm.TableName()
	if err != nil {
		return
	}

	_, err = comm.Options()
	if err != nil {
		return
	}

	_, err = comm.Type()
	if err != nil {
		return
	}

	reviewWg.Wait()
	close(schErrCh)
	if len(schErrCh) > 0 {
		var errs string
		for err := range schErrCh {
			errs += err.Error()
			errs += "; "
		}
		err = errors.New(errs)
		return
	}

	reviewer, err := pR.Team.Delegate(tableName)
	if err != nil {
		return
	}
	comm.Reviewer = reviewer
}
