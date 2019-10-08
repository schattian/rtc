package git

import "github.com/sebach1/git-crud/schema"

type Team struct {
	Schema  *schema.Schema
	Members []Collaborator
}

// func (t *Team) Discuss(tableName integrity.TableName) Collaborator {
// 	for _, collab := t.Members{
// 	}
// }
