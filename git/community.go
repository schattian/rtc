package git

import (
	"github.com/sebach1/rtc/integrity"
)

// A Community delimits the teams whose can take a work
type Community []*Team

// LookFor searches for a team given the correspondent schema (id est: topic abstraction)
func (community *Community) LookFor(schName integrity.SchemaName) (*Team, error) {
	if community == nil {
		return nil, errNilCommunity
	}
	for _, team := range *community {
		if team.AssignedSchema == schName {
			return team, nil
		}
	}
	return nil, errSchemaNotFoundInCommunity
}
