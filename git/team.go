package git

import (
	"github.com/sebach1/rtc/integrity"
)

// A Team is a group of Members which work for the same AssignedSchema
type Team struct {
	AssignedSchema integrity.SchemaName `json:"assigned_schema,omitempty"`
	Members        []*Member            `json:"members,omitempty"`
}

// AddMember validates if a member with the provided args can be created and then adds it to the team
func (t *Team) AddMember(tableName integrity.TableName, collab Collaborator, force bool) error {
	newMember := &Member{}
	for _, member := range t.Members {
		if member.AssignedTable == tableName {
			if !force {
				return errTableInUse
			}
			newMember = member
		}
	}
	newMember.Collab = collab
	if newMember.AssignedTable == "" { // Note: if already is a non-zero assigned table, then it must be a forced
		// type of addition, so instead of deleting the old member it is reused modifying its pointer (then, omit the append)
		newMember.AssignedTable = tableName
		t.Members = append(t.Members, newMember)
	}
	return nil
}

// Delegate retrieves the Collaborator which can perform actions over the given tableName
func (t *Team) Delegate(tableName integrity.TableName) (Collaborator, error) {
	for _, member := range t.Members {
		if member.AssignedTable == tableName {
			if member.Collab == nil {
				return nil, errNoCollaborators
			}
			return member.Collab, nil
		}
	}
	return nil, errNoMembers
}
