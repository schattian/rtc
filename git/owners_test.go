package git

import (
	"context"
	"errors"
	"sync"
	"testing"
)

func TestOwner_Merge(t *testing.T) {
	t.Parallel()
	type args struct {
		ctx context.Context
		pR  *PullRequest
	}
	tests := []struct {
		name      string
		own       *Owner
		args      args
		wantQtErr int
	}{
		{
			name:      "successfull FULL CRUD",
			own:       new(Owner),
			args:      args{pR: gPullRequests.Full.copy().mock(gChanges.Regular.None.TableName, nil), ctx: context.Background()},
			wantQtErr: 0,
		},
		{
			name:      "successfull ONLY one CREATE",
			own:       new(Owner),
			args:      args{pR: gPullRequests.Create.copy().mock(gChanges.Regular.None.TableName, nil), ctx: context.Background()},
			wantQtErr: 0,
		},
		{
			name:      "successfull ONLY one RETRIEVE",
			own:       new(Owner),
			args:      args{pR: gPullRequests.Retrieve.copy().mock(gChanges.Regular.None.TableName, nil), ctx: context.Background()},
			wantQtErr: 0,
		},
		{
			name:      "successfull ONLY one UPDATE",
			own:       new(Owner),
			args:      args{pR: gPullRequests.Update.copy().mock(gChanges.Regular.None.TableName, nil), ctx: context.Background()},
			wantQtErr: 0,
		},
		{
			name:      "successfull ONLY one DELETE",
			own:       new(Owner),
			args:      args{pR: gPullRequests.Delete.copy().mock(gChanges.Regular.None.TableName, nil), ctx: context.Background()},
			wantQtErr: 0,
		},

		{
			name:      "merge with ALL CRUD operations but NO COLLABORATORS",
			own:       new(Owner),
			args:      args{pR: gPullRequests.Full.copy(), ctx: context.Background()},
			wantQtErr: len(gPullRequests.Full.Commits),
		},
		{
			name:      "ERRORED COLLABORATORS",
			own:       new(Owner),
			args:      args{pR: gPullRequests.Full.copy().mock(gChanges.Regular.None.TableName, errors.New("mock")), ctx: context.Background()},
			wantQtErr: len(gPullRequests.Full.Commits),
		},
		{
			name:      "ERRORED COLLABORATORS",
			own:       new(Owner),
			args:      args{pR: gPullRequests.Full.copy().mock(gChanges.Regular.None.TableName, errors.New("mock")), ctx: context.Background()},
			wantQtErr: len(gPullRequests.Full.Commits),
		},
		{
			name: "one of bunch commits is MIXED TABLES",
			own:  new(Owner),
			args: args{
				pR: gPullRequests.Delete.copy().addCommit(
					&Commit{Changes: []*Change{gChanges.Regular.None, gChanges.Regular.TableName}},
				).mock(
					gChanges.Regular.None.TableName, nil,
				),
				ctx: context.Background(),
			},
			wantQtErr: 1,
		},
		{
			name: "one of bunch commits is MIXED TYPES",
			own:  new(Owner),
			args: args{
				pR: gPullRequests.Delete.copy().addCommit(
					&Commit{Changes: []*Change{gChanges.Regular.Create, gChanges.Regular.Update}},
				).mock(
					gChanges.Regular.None.TableName, nil,
				),
				ctx: context.Background(),
			},
			wantQtErr: 1,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			tt.own.wg = new(sync.WaitGroup)
			tt.own.Summary = make(chan *Result, len(tt.args.pR.Commits))
			tt.own.Merge(tt.args.ctx, tt.args.pR)
			tt.own.wg.Wait()
			gotQtErr := len(tt.own.Summary)
			if gotQtErr != tt.wantQtErr {
				t.Errorf("Owner.Merge() errorQt mismatch; got: %v wantQtErr %v", gotQtErr, tt.wantQtErr)
			}
		})
	}
}
