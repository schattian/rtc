package git

import "github.com/sebach1/git-crud/internal/integrity"

type Member struct {
	AssignedTable integrity.TableName `json:"assigned_table,omitempty"`
	Collab        Collaborator        `json:"collab,omitempty"`
}
