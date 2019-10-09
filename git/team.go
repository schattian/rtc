package git

import (
	"github.com/sebach1/git-crud/internal/integrity"
)

// A Team is a group of Members which work for the same AssignedSchema
type Team struct {
	AssignedSchema integrity.SchemaName `json:"assigned_schema,omitempty"`
	Members        []*Member            `json:"members,omitempty"`
}

// Delegate retrieves the Collaborator which can perform actions over the given tableName
func (t *Team) Delegate(tableName integrity.TableName) (Collaborator, error) {
	for _, member := range t.Members {
		if member.AssignedTable == tableName {
			return member.Collab, nil
		}
	}
	return nil, errNoCollaborators
}
