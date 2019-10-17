package git

import "github.com/sebach1/git-crud/integrity"

// A Member is a Collaborator which has a table assigned
type Member struct {
	AssignedTable integrity.TableName `json:"assigned_table,omitempty"`
	Collab        Collaborator        `json:"collab,omitempty"`
}
