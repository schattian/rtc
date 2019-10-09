package git

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/sebach1/git-crud/internal/integrity"
)

func TestTeam_Delegate(t *testing.T) {
	type args struct {
		tableName integrity.TableName
	}
	tests := []struct {
		name    string
		team    *Team
		args    args
		want    Collaborator
		wantErr bool
	}{
		{
			name:    "a member is assigned to the given table",
			team:    gTeams.ZeroMembers.mockedCopy(gChanges.Regular.None.TableName, nil),
			args:    args{tableName: gChanges.Regular.None.TableName},
			want:    &collabMock{},
			wantErr: false,
		},
		{
			name:    "a member isnt assigned to the given table",
			team:    gTeams.ZeroMembers.mockedCopy(gChanges.Regular.None.TableName, nil),
			args:    args{tableName: gChanges.Regular.TableName.TableName},
			want:    &collabMock{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.team.Delegate(tt.args.tableName)
			if (err != nil) != tt.wantErr {
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
