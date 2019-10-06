package git

import (
	"bytes"
	"encoding/json"
	"errors"

	"github.com/sebach1/git-crud/schema"
)

// A Change represents every purposed difference
type Change struct {
	ID int `json:"id,omitempty"`

	TableName  schema.TableName  `json:"table_name,omitempty"`
	ColumnName schema.ColumnName `json:"column_name,omitempty"`

	StrValue     string          `json:"str_value,omitempty"`
	IntValue     int             `json:"int_value,omitempty"`
	Float32Value float32         `json:"float32_value,omitempty"`
	Float64Value float64         `json:"float64_value,omitempty"`
	JSONValue    json.RawMessage `json:"json_value,omitempty"`
	BytesValue   bytes.Buffer    `json:"bytes_value,omitempty"`

	EntityID schema.ID `json:"entity_id,omitempty"`

	Type string
}

// Value gives an interface handling the real value
// Used to perform comparisons
func (chg *Change) Value() interface{} {
	switch chg.Type {
	case "string":
		return chg.StrValue
	case "int":
		return chg.IntValue
	case "float32":
		return chg.Float32Value
	case "float64":
		return chg.Float64Value
	case "json":
		return chg.JSONValue
	case "bytes":
		return chg.BytesValue
	}
	return nil
}

func (chg *Change) tearDownValue() {
	defer func() { chg.Type = "" }()
	switch chg.Type {
	case "string":
		chg.StrValue = ""

	case "int":
		chg.IntValue = 0
	case "float32":
		chg.Float32Value = 0
	case "float64":
		chg.Float64Value = 0

	case "json":
		chg.JSONValue = json.RawMessage{}
	case "bytes":
		chg.BytesValue = bytes.Buffer{}
	}
}

// SetValue performs type assertion over the given value and sets the value over the given change
// In case of failure on all the possible type assertions, returns an error
// Notice that SetValue will ALWAYS tearDown the value set up before
func (chg *Change) SetValue(val interface{}) (err error) {
	chg.tearDownValue()

	if strVal, ok := val.(string); ok {
		chg.StrValue = strVal
		chg.Type = "string"
		return
	}
	if intVal, ok := val.(int); ok {
		chg.IntValue = intVal
		chg.Type = "int"
		return
	}
	if f32Val, ok := val.(float32); ok {
		chg.Float32Value = f32Val
		chg.Type = "float32"
		return
	}
	if f64Val, ok := val.(float64); ok {
		chg.Float64Value = f64Val
		chg.Type = "float64"
		return
	}
	if jsVal, ok := val.(json.RawMessage); ok {
		chg.JSONValue = jsVal
		chg.Type = "json"
		return
	}
	if byBuff, ok := val.(bytes.Buffer); ok {
		chg.BytesValue = byBuff
		chg.Type = "bytes"
		return
	}

	return errors.New("the given value cannot be safety typed")
}

// Validate self
func (chg *Change) Validate() error {
	if chg.TableName == "" {
		return errZeroTable
	}
	if chg.ColumnName == "" {
		return errZeroColumn
	}
	return nil
}

// IsUntracked retrieves true if the change is a new entity, otherwise returns false
func (chg *Change) IsUntracked() bool {
	if chg.EntityID.IsNil() {
		return true
	}
	return false
}

// Equals checks if the given change will trigger the exactly same action as itself
// That is: the change is not untracked, the column and table affected are the same
// and the value is the same
func (chg *Change) Equals(otherChg *Change) bool {
	if !IsCompatibleWith(chg, otherChg) {
		return false
	}
	if chg.ColumnName != otherChg.ColumnName {
		return false
	}
	if chg.Value() != otherChg.Value() {
		return false
	}
	return true
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
	if chg.IsUntracked() { // In case of both of EntityIDs are nil (see that the above comparison discards 2x checking)
		return false
	}
	return true
}
