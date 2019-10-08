package git

import "github.com/sebach1/git-crud/internal/integrity"

type Community []*Team

func (community *Community) LookFor(schName integrity.SchemaName) (*Team, error) {
	if community == nil {
		return nil, errNilCommunity
	}
	for _, team := range *community {
		if team.Schema.Name == schName {
			return team, nil
		}
	}
	return nil, errNotFoundSchema
}
