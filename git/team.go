package git

import (
	"github.com/sebach1/git-crud/internal/integrity"
)

type Team struct {
	AssignedSchema integrity.SchemaName
	Members        []*Member
}

func (t *Team) Delegate(tableName integrity.TableName) (Collaborator, error) {
	for _, member := range t.Members {
		// member.
		if member.AssignedTable == tableName {
			return member.Collab, nil
		}
	}
	return nil, errNoCollaborators
}
