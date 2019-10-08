package git

import (
	"context"
	"sync"

	"github.com/sebach1/git-crud/internal/integrity"
)

type Owner struct {
	wg         *sync.WaitGroup
	mergeConfs chan error
	summ       chan *Result
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
	own.wg = new(sync.WaitGroup)
	var pR PullRequest

	for _, changes := range comm.GroupBy(strategy) { // Splits incompatibilities onto the pR
		comm := &Commit{Changes: changes}
		pR.Commits = append(pR.Commits, comm)
	}
	pR.AssignTeam(community, schName)
	go own.Merge(ctx, &pR)
	return
}

func (own *Owner) Merge(ctx context.Context, pR *PullRequest) {
	for _, comm := range pR.Commits {

		tableName, err := comm.TableName()
		if err != nil {
			own.mergeConfs <- err
		}

		commType, err := comm.Type()
		if err != nil {
			own.mergeConfs <- err
			continue // Notice that this action discards the commit
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
	own.wg.Wait()
}

func (own *Owner) Pull(ctx context.Context, comm *Commit) (*Commit, error) {
	defer own.wg.Done()

	pR.Team.Delegate(tableName)
	return comm, err
}

func (own *Owner) Push(ctx context.Context, comm *Commit) (Summary, error) {
	defer own.wg.Done()
	return own.summ, err
}

func (own *Owner) Delete(ctx context.Context, comm *Commit) (Summary, error) {
	defer own.wg.Done()

	return own.summ, err
}
