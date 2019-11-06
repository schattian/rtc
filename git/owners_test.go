package git

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/sebach1/git-crud/integrity"
	"github.com/sebach1/git-crud/schema"
)

// func TestOwner_Merge(t *testing.T) {
// 	t.Parallel()
// 	type args struct {
// 		ctx context.Context
// 		pR  *PullRequest
// 	}
// 	tests := []struct {
// 		name      string
// 		own       *Owner
// 		args      args
// 		wantQtErr int
// 	}{
// 		{
// 			name:      "successfull FULL CRUD",
// 			own:       new(Owner),
// 			args:      args{pR: gPullRequests.Full.copy().mock(gChanges.Basic.None.TableName, nil), ctx: context.Background()},
// 			wantQtErr: 0,
// 		},
// 		{
// 			name:      "successfull ONLY one CREATE",
// 			own:       new(Owner),
// 			args:      args{pR: gPullRequests.Create.copy().mock(gChanges.Basic.None.TableName, nil), ctx: context.Background()},
// 			wantQtErr: 0,
// 		},
// 		{
// 			name:      "successfull ONLY one RETRIEVE",
// 			own:       new(Owner),
// 			args:      args{pR: gPullRequests.Retrieve.copy().mock(gChanges.Basic.None.TableName, nil), ctx: context.Background()},
// 			wantQtErr: 0,
// 		},
// 		{
// 			name:      "successfull ONLY one UPDATE",
// 			own:       new(Owner),
// 			args:      args{pR: gPullRequests.Update.copy().mock(gChanges.Basic.None.TableName, nil), ctx: context.Background()},
// 			wantQtErr: 0,
// 		},
// 		{
// 			name:      "successfull ONLY one DELETE",
// 			own:       new(Owner),
// 			args:      args{pR: gPullRequests.Delete.copy().mock(gChanges.Basic.None.TableName, nil), ctx: context.Background()},
// 			wantQtErr: 0,
// 		},

// 		{
// 			name:      "merge with ALL CRUD operations but NO COLLABORATORS",
// 			own:       new(Owner),
// 			args:      args{pR: gPullRequests.Full.copy(), ctx: context.Background()},
// 			wantQtErr: len(gPullRequests.Full.Commits),
// 		},
// 		{
// 			name:      "ERRORED COLLABORATORS",
// 			own:       new(Owner),
// 			args:      args{pR: gPullRequests.Full.copy().mock(gChanges.Basic.None.TableName, errors.New("mock")), ctx: context.Background()},
// 			wantQtErr: len(gPullRequests.Full.Commits),
// 		},
// 		{
// 			name:      "ERRORED COLLABORATORS",
// 			own:       new(Owner),
// 			args:      args{pR: gPullRequests.Full.copy().mock(gChanges.Basic.None.TableName, errors.New("mock")), ctx: context.Background()},
// 			wantQtErr: len(gPullRequests.Full.Commits),
// 		},
// 		{
// 			name: "one of bunch commits is MIXED TABLES",
// 			own:  new(Owner),
// 			args: args{
// 				pR: gPullRequests.Delete.copy().addCommit(
// 					&Commit{Changes: []*Change{gChanges.Basic.None, gChanges.Basic.TableName}},
// 				).mock(
// 					gChanges.Basic.None.TableName, nil,
// 				),
// 				ctx: context.Background(),
// 			},
// 			wantQtErr: 1,
// 		},
// 		{
// 			name: "one of bunch commits is MIXED TYPES",
// 			own:  new(Owner),
// 			args: args{
// 				pR: gPullRequests.Delete.copy().addCommit(
// 					&Commit{Changes: []*Change{gChanges.Basic.Create, gChanges.Basic.Update}},
// 				).mock(
// 					gChanges.Basic.None.TableName, nil,
// 				),
// 				ctx: context.Background(),
// 			},
// 			wantQtErr: 1,
// 		},
// 	}
// 	for _, tt := range tests {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()
// 			tt.own.wg = new(sync.WaitGroup)
// 			tt.own.Summary = make(chan *Result, len(tt.args.pR.Commits))
// 			tt.own.Merge(tt.args.ctx, tt.args.pR)
// 			tt.own.wg.Wait()
// 			gotQtErr := len(tt.own.Summary)
// 			if gotQtErr != tt.wantQtErr {
// 				t.Errorf("Owner.Merge() errorQt mismatch; got: %v wantQtErr %v", gotQtErr, tt.wantQtErr)
// 			}
// 		})
// 	}
// }

func TestOwner_ReviewPRCommit(t *testing.T) {
	t.Parallel()
	type args struct {
		sch *schema.Schema
		pR  *PullRequest
	}
	tests := []struct {
		name      string
		own       *Owner
		args      args
		wantQtErr int
	}{
		{
			name: "successfull FULL CRUD",
			own:  &Owner{Project: &schema.Planisphere{gSchemas.Basic}},
			args: args{
				sch: gSchemas.Basic,
				pR:  gPullRequests.Full.copy().mock(gTables.Basic.Name, nil),
			},
			wantQtErr: 0,
		},
		{
			name: "NO COLLABORATORS",
			own:  &Owner{Project: &schema.Planisphere{gSchemas.Basic}},
			args: args{
				sch: gSchemas.Basic,
				pR:  gPullRequests.Full.copy(),
			},
			wantQtErr: len(gPullRequests.Full.Commits),
		},
		{
			name: "commit is MIXED TABLES",
			own:  &Owner{Project: &schema.Planisphere{gSchemas.Basic}},
			args: args{
				sch: gSchemas.Basic,
				pR: gPullRequests.ZeroCommits.copy().addCommit(
					&Commit{Changes: []*Change{gChanges.Basic.None.copy(), gChanges.Basic.TableName.copy()}},
				).mock(gTables.Basic.Name, nil),
			},
			wantQtErr: 1,
		},
		{
			name: "commit CHANGE IS INCONSISTENT",
			own:  &Owner{Project: &schema.Planisphere{gSchemas.Basic}},
			args: args{
				sch: gSchemas.Basic,
				pR: gPullRequests.ZeroCommits.copy().addCommit(
					&Commit{Changes: []*Change{gChanges.Inconsistent.Delete}},
				).mock(gTables.Basic.Name, nil),
			},
			wantQtErr: 1,
		},
		{
			name: "commit is MIXED OPTIONS",
			own:  &Owner{Project: &schema.Planisphere{gSchemas.Basic}},
			args: args{
				sch: gSchemas.Basic,
				pR: gPullRequests.ZeroCommits.copy().addCommit(
					&Commit{Changes: []*Change{gChanges.Basic.None.copy(), gChanges.Rare.TableName.copy()}},
				).mock(gTables.Basic.Name, nil),
			},
			wantQtErr: 1,
		},
		{
			name: "commit is MIXED TYPES",
			own:  &Owner{Project: &schema.Planisphere{gSchemas.Basic}},
			args: args{
				sch: gSchemas.Basic,
				pR: gPullRequests.ZeroCommits.copy().addCommit(
					&Commit{Changes: []*Change{gChanges.Basic.Create.copy(), gChanges.Basic.Update.copy()}},
				).mock(gTables.Basic.Name, nil),
			},
			wantQtErr: 1,
		},
		{
			name: "commit does NOT PASSES the SCHEMA VALIDATION",
			own:  &Owner{Project: &schema.Planisphere{gSchemas.Basic, gSchemas.Rare}},
			args: args{
				sch: gSchemas.Rare,
				pR:  gPullRequests.Basic.copy().mock(gTables.Basic.Name, nil),
			},
			wantQtErr: 1,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.own.Summary = make(chan *Result, len(tt.args.pR.Commits))
			wg := &sync.WaitGroup{}
			wg.Add(len(tt.args.pR.Commits))
			for commIdx := range tt.args.pR.Commits {
				go tt.own.ReviewPRCommit(tt.args.sch, tt.args.pR, commIdx, wg)
			}
			wg.Wait()
			var gotQtErr int
			for _, comm := range tt.args.pR.Commits {
				if comm.Errored {
					gotQtErr++
				}
			}
			if gotQtErr != tt.wantQtErr {
				t.Errorf("Owner.ReviewPRCommit() errorQt mismatch; got: %v wantQtErr %v", gotQtErr, tt.wantQtErr)
			}
		})
	}
}

func TestOwner_Orchestrate(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx       context.Context
		community *Community
		schName   integrity.SchemaName
		comm      *Commit
		strategy  changesMatcher
	}
	tests := []struct {
		name          string
		own           *Owner
		args          args
		wantErr       bool
		wantQtResErrs int // Quantity of results in summary that are errored
	}{
		{
			name: "fully successful",
			own:  newOwnerUnsafe(&schema.Planisphere{gSchemas.Basic}),
			args: args{
				ctx:       context.Background(),
				community: &Community{gTeams.Basic.copy().mock(gChanges.Basic.None.TableName, nil)},
				schName:   gSchemas.Basic.Name,
				comm: &Commit{Changes: []*Change{
					gChanges.Basic.Create.copy(),
					gChanges.Basic.Retrieve.copy(),
					gChanges.Basic.Update.copy(),
					gChanges.Basic.Delete.copy(),
				}},
				strategy: AreCompatible,
			},
			wantErr:       false,
			wantQtResErrs: 0,
		},
		{
			name: "but NIL PROJECT",
			own:  newOwnerUnsafe(nil),
			args: args{
				ctx:       context.Background(),
				community: &Community{gTeams.Basic},
				schName:   gSchemas.Basic.Name,
				comm:      &Commit{Changes: []*Change{gChanges.Basic.None.copy()}},
				strategy:  AreCompatible,
			},
			wantErr: true,
		},
		{
			name: "but NO COLLABORATORS",
			own:  newOwnerUnsafe(&schema.Planisphere{gSchemas.Basic}),
			args: args{
				ctx:       context.Background(),
				community: &Community{gTeams.Basic.copy()},
				schName:   gSchemas.Basic.Name,
				comm:      &Commit{Changes: []*Change{gChanges.Basic.None.copy()}},
				strategy:  AreCompatible,
			},
			wantErr:       false,
			wantQtResErrs: 1,
		},
		{
			name: "but COLLABORATORS MOCK RETURNS ERRS",
			own:  newOwnerUnsafe(&schema.Planisphere{gSchemas.Basic}),
			args: args{
				ctx:       context.Background(),
				community: &Community{gTeams.Basic.copy().mock(gChanges.Basic.None.TableName, errors.New("test"))},
				schName:   gSchemas.Basic.Name,
				comm: &Commit{Changes: []*Change{
					gChanges.Basic.Create.copy(),
					gChanges.Basic.Retrieve.copy(),
					gChanges.Basic.Update.copy(),
					gChanges.Basic.Delete.copy(),
				}},
				strategy: AreCompatible,
			},
			wantErr:       false,
			wantQtResErrs: 4,
		},
		{
			name: "given SCHEMA NOT IN PLANISPHERE",
			own:  newOwnerUnsafe(&schema.Planisphere{gSchemas.Rare}),
			args: args{
				ctx:       context.Background(),
				community: &Community{gTeams.Basic.copy().mock(gChanges.Basic.None.TableName, nil)},
				schName:   gSchemas.Basic.Name,
				comm:      &Commit{Changes: []*Change{gChanges.Basic.None.copy()}},
				strategy:  AreCompatible,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.own.Waiter.Add(1) // Prevent a line of the boilerplate when calling orchestra
			go tt.own.Orchestrate(tt.args.ctx, tt.args.community, tt.args.schName, tt.args.comm, tt.args.strategy)
			err := tt.own.Close()
			if (err != nil) != tt.wantErr {
				t.Errorf("Owner.Orchestrate() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err != nil {
				return
			}

			var gotQtResErrs int
			var gotErrs string
			for result := range tt.own.Summary {
				if result.Error != nil {
					gotErrs += result.Error.Error()
					gotErrs += integrity.ErrorsSeparator
					gotQtResErrs++
				}
			}

			if gotQtResErrs != tt.wantQtResErrs {
				t.Errorf("Owner.Orchestrate() gotQtResErrs = %v, wantQtResErrs %v", gotQtResErrs, tt.wantQtResErrs)
				t.Logf("HINT: Owner.Orchestrate() gotErrs = %v", gotErrs)
			}
		})
	}
}
