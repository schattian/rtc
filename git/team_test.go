package git

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sebach1/git-crud/integrity"
)

func TestTeam_Delegate(t *testing.T) {
	t.Parallel()
	type args struct {
		tableName integrity.TableName
	}
	tests := []struct {
		name    string
		team    *Team
		args    args
		want    Collaborator
		wantErr error
	}{
		{
			name:    "a member is assigned to the given table",
			team:    gTeams.ZeroMembers.copy().mock(gChanges.Foo.None.TableName, nil),
			args:    args{tableName: gChanges.Foo.None.TableName},
			want:    &collabMock{},
			wantErr: nil,
		},
		{
			name:    "a member isn't assigned to the given table",
			team:    gTeams.ZeroMembers.copy().mock(gChanges.Foo.None.TableName, nil),
			args:    args{tableName: gChanges.Foo.TableName.TableName},
			want:    &collabMock{},
			wantErr: errNoMembers,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.team.Delegate(tt.args.tableName)
			if err != tt.wantErr {
				t.Errorf("Team.Delegate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				if diff := cmp.Diff(nil, got); diff != "" {
					t.Errorf("Team.Delegate() errored mismatch (-want +got): %s", diff)
				}
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("Team.Delegate() mismatch (-want +got): %s", diff)
			}
		})
	}
}
