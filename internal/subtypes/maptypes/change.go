package maptypes

// A Change represents every purposed difference
type Change struct {
	ID int `json:"id"`

	TableName  TableName  `json:"table_name"`
	ColumnName ColumnName `json:"column_name"`

	Value Value `json:"value"`

	EntityID ID `json:"entity_id,omitempty"`
}

// IsUntracked retrieves true if the change is a new entity, otherwise returns false
func (chg *Change) IsUntracked() bool {
	if chg.EntityID == nil {
		return true
	}
	return false
}

// IsCompatibleWith checks if the changes belongs to the same table and the same entity
// Notice that the func is not using syntatic sugar to be a comparator type (see Commit.GroupBy)
func IsCompatibleWith(chg, otherChg *Change) bool {
	if chg.TableName != otherChg.TableName {
		return false
	}
	if chg.EntityID != otherChg.EntityID {
		return false
	}
	if chg.EntityID == nil { // In case of both of them are nil (see that the above comparison discards 2x checking)
		return false
	}
	return true
}
