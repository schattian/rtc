package git

import (
	"encoding/json"
	"errors"

	"github.com/sebach1/git-crud/internal/integrity"
	"github.com/sebach1/git-crud/schema"
)

// A Change represents every purposed difference
type Change struct {
	TableName  integrity.TableName  `json:"table_name,omitempty"`
	ColumnName integrity.ColumnName `json:"column_name,omitempty"`

	StrValue     string          `json:"str_value,omitempty"`
	IntValue     int             `json:"int_value,omitempty"`
	Float32Value float32         `json:"float32_value,omitempty"`
	Float64Value float64         `json:"float64_value,omitempty"`
	JSONValue    json.RawMessage `json:"json_value,omitempty"`
	BytesValue   []byte          `json:"bytes_value,omitempty"`

	EntityID schema.ID `json:"entity_id,omitempty"`

	ValueType string `json:"value_type,omitempty"`
}

// Value gives an interface handling the real value
// Used to perform comparisons
func (chg *Change) Value() interface{} {
	switch chg.ValueType {
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
	defer func() { chg.ValueType = "" }()
	switch chg.ValueType {
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
		chg.BytesValue = nil
	}
}

// SetValue performs type assertion over the given value and sets the value over the given change
// In case of failure on all the possible type assertions, returns an error
// Notice that SetValue will ALWAYS tearDown the value set up before
func (chg *Change) SetValue(val interface{}) (err error) {
	chg.tearDownValue()

	if strVal, ok := val.(string); ok {
		chg.StrValue = strVal
		chg.ValueType = "string"
		return
	}
	if intVal, ok := val.(int); ok {
		chg.IntValue = intVal
		chg.ValueType = "int"
		return
	}
	if f32Val, ok := val.(float32); ok {
		chg.Float32Value = f32Val
		chg.ValueType = "float32"
		return
	}
	if f64Val, ok := val.(float64); ok {
		chg.Float64Value = f64Val
		chg.ValueType = "float64"
		return
	}
	if jsVal, ok := val.(json.RawMessage); ok {
		chg.JSONValue = jsVal
		chg.ValueType = "json"
		return
	}
	if byVal, ok := val.([]byte); ok {
		chg.BytesValue = byVal
		chg.ValueType = "bytes"
		return
	}

	return errors.New("the given value cannot be safety typed")
}

type changesMatcher func(*Change, *Change) bool

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

// Overrides check if the changes will generate be overrided
func Overrides(chg, otherChg *Change) bool {
	if !AreCompatible(chg, otherChg) {
		return false
	}
	if chg.ColumnName != otherChg.ColumnName {
		return false
	}
	return true
}

// Equals checks if the given change will trigger the exactly same action as itself
// Note that Equals prioritizes false on cryteria
// So, in case needing to check untracked it can use AreCompatibleOrUntracked func
func (chg *Change) Equals(otherChg *Change) bool {
	if !AreCompatible(chg, otherChg) {
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

// AreCompatible uses AreCompatibleOnUntracked, plus discarding untracked changes
func AreCompatible(chg, otherChg *Change) bool {
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
