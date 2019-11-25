package git

import "github.com/sebach1/git-crud/integrity"

// A Member is a Collaborator which has a table assigned
type Member struct {
	AssignedTable integrity.TableName `json:"assigned_table,omitempty"`
	Collab        Collaborator        `json:"collab,omitempty"`
}

func (m *Member) copy() *Member {
	if m == nil {
		return nil
	}
	member := &Member{}
	*member = *m
	return member
}
