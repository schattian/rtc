package git

import (
	"encoding/json"
	"errors"

	"github.com/sebach1/git-crud/internal/integrity"
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

	EntityID integrity.ID `json:"entity_id,omitempty"`

	ValueType string `json:"value_type,omitempty"`

	Type integrity.CRUD `json:"type,omitempty"`
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

// AreCompatible checks if the given changes can be used to merge them into the same action
func AreCompatible(chg, otherChg *Change) bool {
	if chg.TableName != otherChg.TableName {
		return false
	}
	if chg.EntityID != otherChg.EntityID {
		return false
	}
	if chg.Type == "create" { // In case of both of EntityIDs are nil
		//  (see that the above comparison discards 2x checking)
		return false
	}
	if chg.Type != otherChg.Type {
		return false
	}
	return true
}

// Validate self
func (chg *Change) Validate() error {
	if chg.TableName == "" {
		return errNilTable
	}
	if chg.Type == "" {
		newType, err := chg.classifyType()
		if err != nil {
			return err
		}
		chg.Type = newType
	} else {
		err := chg.validateType()
		if err != nil {
			return err
		}
	}

	return nil
}

// Notice that this implementation could be did just because CRUD patterns are mutually exclusive
func (chg *Change) classifyType() (integrity.CRUD, error) {
	if chg.validateCreate() == nil {
		return "create", nil
	}
	if chg.validateRetrieve() == nil {
		return "retrieve", nil
	}
	if chg.validateUpdate() == nil {
		return "update", nil
	}
	if chg.validateDelete() == nil {
		return "delete", nil
	}
	return "", errUnclassifiableChg
}

func (chg *Change) validateType() error {
	err := chg.Type.Validate()
	if err != nil {
		return err
	}
	switch chg.Type {
	case "create":
		err = chg.validateCreate()
	case "retrieve":
		err = chg.validateRetrieve()
	case "update":
		err = chg.validateUpdate()
	case "delete":
		err = chg.validateDelete()
	}
	if err != nil {
		return err
	}
	return nil
}

func (chg *Change) validateCreate() error {
	if chg.EntityID.IsNil() {
		return errNotNilEntityID
	}
	if chg.ValueType == "" {
		return errNilValue
	}
	return nil
}

func (chg *Change) validateRetrieve() (err error) {
	if chg.ValueType != "" {
		return errNotNilValue
	}
	return nil
}

func (chg *Change) validateUpdate() error {
	if chg.EntityID.IsNil() {
		return errNilEntityID
	}
	if chg.ColumnName == "" {
		return errNilColumn
	}
	if chg.ValueType == "" {
		return errNilValue
	}
	return nil
}

func (chg *Change) validateDelete() error {
	if chg.EntityID.IsNil() {
		return errNilEntityID
	}
	if chg.ValueType != "" {
		return errNotNilValue
	}
	if chg.ColumnName != "" {
		return errNilColumn
	}
	return nil
}
