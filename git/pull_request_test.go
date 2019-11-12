package git

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sebach1/git-crud/integrity"
)

func TestPullRequest_AssignTeam(t *testing.T) {
	t.Parallel()
	type args struct {
		community *Community
		schName   integrity.SchemaName
	}
	tests := []struct {
		name     string
		pR       *PullRequest
		args     args
		wantErr  bool
		wantTeam *Team
	}{
		{
			name: "correctly assigns the team",
			pR:   gPullRequests.ZeroTeam.copy(),
			args: args{
				community: &Community{gTeams.Foo},
				schName:   gTeams.Foo.AssignedSchema,
			},
			wantErr:  false,
			wantTeam: gTeams.Foo,
		},
		{
			name: "team not in the community",
			pR:   gPullRequests.ZeroTeam.copy(),
			args: args{
				community: &Community{gTeams.Bar},
				schName:   gTeams.Foo.AssignedSchema,
			},
			wantErr: true,
		},
		{
			name: "given community is nil",
			pR:   gPullRequests.ZeroTeam.copy(),
			args: args{
				community: nil,
				schName:   gTeams.Foo.AssignedSchema,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.pR.AssignTeam(tt.args.community, tt.args.schName)
			if (err != nil) != tt.wantErr {
				t.Errorf("PullRequest.AssignTeam() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil {
				if diff := cmp.Diff(gTeams.Zero, tt.pR.Team); diff != "" {
					t.Errorf("PullRequest.AssignTeam() errored mismatch (-want +got): %s", diff)
				}
				return
			}
			if diff := cmp.Diff(tt.wantTeam, tt.pR.Team); diff != "" {
				t.Errorf("PullRequest.AssignTeam() mismatch (-want +got): %s", diff)
			}
		})
	}
}
