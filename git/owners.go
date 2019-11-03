package git

import (
	"context"
	"errors"
	"sync"

	"github.com/sebach1/git-crud/integrity"
	"github.com/sebach1/git-crud/schema"
)

// Owner is the agent which coordinates any given action
// Notice that an Owner is a Collaborator
type Owner struct {
	Project *schema.Planisphere
	Summary chan *Result

	Waiter *sync.WaitGroup
	err    error
}

// NewOwner returns a new instance of Owner, with needed initialization and validation
func NewOwner(project *schema.Planisphere) (*Owner, error) {
	if project == nil || len(*project) == 0 {
		return nil, errors.New("The PROJECT cannot be NIL")
	}
	return newOwnerUnsafe(project), nil
}

// newOwnerUnsafe returns a new instance of Owner, with needed initialization
func newOwnerUnsafe(project *schema.Planisphere) *Owner {
	own := &Owner{Project: project}
	own.Waiter = &sync.WaitGroup{}
	return own
}

// Orchestrate sends the order to all the collaborators available to execute
// the needed actions in order to achieve the commitment, creating a new PullRequest
func (own *Owner) Orchestrate(
	ctx context.Context,
	community *Community,
	schName integrity.SchemaName,
	comm *Commit,
	strategy changesMatcher,
) {
	defer own.Waiter.Done()
	pR, err := own.Delegate(ctx, community, schName, comm, strategy)
	if err != nil {
		own.err = err
		return
	}
	own.Waiter.Add(1)
	go own.Merge(ctx, pR)
}

// Delegate creates a PullRequest and assigns reviewers from a given commit
func (own *Owner) Delegate(
	ctx context.Context,
	community *Community,
	schName integrity.SchemaName,
	comm *Commit,
	strategy changesMatcher,
) (*PullRequest, error) {

	err := own.Validate()
	if err != nil {
		return nil, err
	}

	sch, err := own.Project.GetSchemaFromName(schName)
	if err != nil {
		return nil, err
	}

	var pR PullRequest
	var wg sync.WaitGroup

	for _, changes := range comm.GroupBy(strategy) { // Splits incompatibilities onto the pR
		comm := &Commit{Changes: changes}
		pR.Commits = append(pR.Commits, comm)
	}
	pR.AssignTeam(community, schName)

	own.Summary = make(chan *Result, len(pR.Commits))

	wg.Add(len(pR.Commits))
	for commIdx := range pR.Commits {
		go own.ReviewPRCommit(sch, &pR, commIdx, &wg)
	}
	wg.Wait()

	return &pR, nil
}

// Close will wait for the Owner WaitGroup to be done and close the Owner.Summary
// It closes an orchestration (Owner.Orchestrate())
func (own *Owner) Close() error {
	own.Waiter.Wait()
	if own.Summary != nil {
		// The channel can be nil if the owner was errored before
		// a merge / review
		close(own.Summary)
	}
	return own.err
}

// Merge performs the needed actions in order to merge the pullRequest
func (own *Owner) Merge(ctx context.Context, pR *PullRequest) {
	defer own.Waiter.Done()
	for _, comm := range pR.Commits {
		if comm.Errored {
			continue // Skips validation errs
		}

		commType, err := comm.Type()
		if err != nil {
			own.Summary <- &Result{CommitID: comm.ID, Error: err}
			continue
		}

		own.Waiter.Add(1)
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
	defer own.Waiter.Done()
	newComm := new(Commit)
	*newComm = *comm
	err := comm.Reviewer.Init(ctx)
	if err != nil {
		own.Summary <- &Result{CommitID: comm.ID, Error: err}
		return comm, err
	}
	newComm, err = comm.Reviewer.Push(ctx, newComm)
	if err != nil {
		own.Summary <- &Result{CommitID: comm.ID, Error: err}
		return comm, err
	}
	*comm = *newComm
	return comm, nil
}

// Pull will orchestrate the pulls of any collaborator
func (own *Owner) Pull(ctx context.Context, comm *Commit) (*Commit, error) {
	defer own.Waiter.Done()
	newComm := new(Commit)
	*newComm = *comm
	err := comm.Reviewer.Init(ctx)
	if err != nil {
		own.Summary <- &Result{CommitID: comm.ID, Error: err}
		return comm, err
	}
	newComm, err = comm.Reviewer.Pull(ctx, newComm)
	if err != nil {
		own.Summary <- &Result{CommitID: comm.ID, Error: err}
		return comm, err
	}
	*comm = *newComm
	return comm, nil
}

// Delete will orchestrate the deletions of any collaborator
func (own *Owner) Delete(ctx context.Context, comm *Commit) (*Commit, error) {
	defer own.Waiter.Done()
	newComm := new(Commit)
	*newComm = *comm
	err := comm.Reviewer.Init(ctx)
	if err != nil {
		own.Summary <- &Result{CommitID: comm.ID, Error: err}
		return comm, err
	}
	newComm, err = comm.Reviewer.Delete(ctx, newComm)
	if err != nil {
		own.Summary <- &Result{CommitID: comm.ID, Error: err}
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
func (own *Owner) ReviewPRCommit(sch *schema.Schema, pR *PullRequest, commIdx int, delegationWg *sync.WaitGroup) {
	var err error
	defer delegationWg.Done()
	var reviewWg sync.WaitGroup

	comm := pR.Commits[commIdx]
	defer func() {
		if err != nil {
			own.Summary <- &Result{CommitID: comm.ID, Error: err}
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
		go sch.ValidateCtx(chg.TableName, chg.ColumnName, chg.Options.Keys(), chg.Value(),
			own.Project, &reviewWg, schErrCh)
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
			errs += integrity.ErrorsSeparator
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
