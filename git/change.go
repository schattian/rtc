package git

import (
	"github.com/sebach1/rtc/integrity"
)

// A Change represents every purposed/lookup for difference
type Change struct {
	Id int64 `json:"id,omitempty"`

	TableName  integrity.TableName  `json:"table_name,omitempty"`
	ColumnName integrity.ColumnName `json:"column_name,omitempty"`

	StringValue  string  `json:"string_value,omitempty"`
	IntValue     int     `json:"int_value,omitempty"`
	Float32Value float32 `json:"float_32_value,omitempty"`
	Float64Value float64 `json:"float_64_value,omitempty"`
	BytesValue   []byte  `json:"bytes_value,omitempty"`

	EntityId integrity.Id `json:"entity_id,omitempty"`

	ValueType integrity.ValueType `json:"value_type,omitempty"`

	Type integrity.CRUD `json:"type,omitempty"`

	Options Options

	IndexId int64 `json:"index_id,omitempty"`

	CommitId int64 `json:"commit_id,omitempty"`
}

// NewChange safety creates a new Change entity
// Notice it doesn't saves it on the db
func NewChange(
	entityId integrity.Id,
	tableName integrity.TableName,
	columnName integrity.ColumnName,
	val interface{},
	Type integrity.CRUD,
	opts Options,
) (*Change, error) {
	chg := &Change{EntityId: entityId, TableName: tableName, ColumnName: columnName, Type: Type, Options: opts}
	err := chg.SetValue(val)
	if err != nil {
		return nil, err
	}

	err = chg.Validate()
	if err != nil {
		return nil, err
	}
	return chg, nil
}

// SetOption assigns the given key to the given value. Returns an error if the key is not allowed for any option
func (chg *Change) SetOption(key integrity.OptionKey, val interface{}) error {
	if key == "" {
		return errNilOptionKey
	}
	if chg.Options == nil {
		chg.Options = make(Options)
	}
	chg.Options[key] = val
	return nil
}

// Value gives an interface handling the real value
// Used to perform comparisons
func (chg *Change) Value() interface{} {
	switch chg.ValueType {
	case "string":
		return chg.StringValue
	case "int":
		return chg.IntValue
	case "float_32":
		return chg.Float32Value
	case "float_64":
		return chg.Float64Value
	case "json", "bytes":
		return chg.BytesValue
	}
	return nil
}

// SetValue performs type assertion over the given value and sets the value over the given change
// In case of failure on all the possible type assertions, returns an error
// Notice that SetValue will ALWAYS tearDown the value set up before
func (chg *Change) SetValue(val interface{}) (err error) {
	chg.tearDownValue()

	if strVal, ok := val.(string); ok {
		chg.StringValue = strVal
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
		chg.ValueType = "float_32"
		return
	}
	if f64Val, ok := val.(float64); ok {
		chg.Float64Value = f64Val
		chg.ValueType = "float_64"
		return
	}
	if byVal, ok := val.([]byte); ok {
		chg.BytesValue = byVal
		chg.ValueType = "bytes"
		return
	}
	return errUnsafeValueType
}

// Overrides check if changes are overridable by each other
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
// skips comparison in case of Id matching (and Id exists for both of 'em)
func (chg *Change) Equals(otherChg *Change) bool {
	if chg.Id == otherChg.Id && (chg.Id != 0) {
		return true
	}

	if !AreCompatible(chg, otherChg) {
		return false
	}
	if chg.ColumnName != otherChg.ColumnName {
		return false
	}
	if chg.Value() != otherChg.Value() {
		return false
	}
	if len(chg.Options) != len(otherChg.Options) {
		return false
	}
	for k, v := range chg.Options {
		if otherChg.Options[k] != v {
			return false
		}
	}

	return true
}

// Validate self, wrapping up type validations and table assertion
func (chg *Change) Validate() (err error) {
	if chg.TableName == "" {
		err = errNilTable
		return
	}
	if chg.Type != "" {
		err = chg.validateType()
		if err != nil {
			return
		}
	}

	newType, err := chg.classifyType()
	if err != nil {
		return
	}
	chg.Type = newType
	return nil
}

// FromMap decodes the commit from its map version
// Notice that FromMap() is reciprocal to ToMap(), so it doesn't assign a table
func (chg *Change) FromMap(Map map[string]interface{}) error {
	for col, val := range Map {
		if col == "id" {
			realVal, ok := val.(integrity.Id)
			if !ok {
				return errInvalidChangeId
			}
			chg.EntityId = realVal
			continue
		}
		chg.ColumnName = integrity.ColumnName(col)
		err := chg.SetValue(val)
		if err != nil {
			return err
		}

	}
	return nil
}

// ToMap retrieves a map with the minimum required -not validable- data
// id est: {column_name: value}
func (chg *Change) ToMap() map[string]interface{} {
	chgMap := make(map[string]interface{})
	if chg.ColumnName != "" {
		chgMap[string(chg.ColumnName)] = chg.Value()
	}
	if !chg.EntityId.IsNil() {
		chgMap["id"] = chg.EntityId
	}
	return chgMap
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

// ClassifyType will auto-bind the change to the accurate type
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

func (chg *Change) validateCreate() error {
	if !chg.EntityId.IsNil() {
		return errNotNilEntityId
	}
	if chg.ColumnName == "" {
		return errNilColumn
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
	if chg.ColumnName == "" {
		return errNilColumn
	}
	return nil
}

func (chg *Change) validateUpdate() error {
	if chg.EntityId.IsNil() {
		return errNilEntityId
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
	if chg.EntityId.IsNil() {
		return errNilEntityId
	}
	if chg.ValueType != "" {
		return errNotNilValue
	}
	if chg.ColumnName != "" {
		return errNotNilColumn
	}
	return nil
}

func (chg *Change) tearDownValue() {
	defer func() { chg.ValueType = "" }()
	switch chg.ValueType {
	case "string":
		chg.StringValue = ""
	case "int":
		chg.IntValue = 0
	case "float32":
		chg.Float32Value = 0
	case "float64":
		chg.Float64Value = 0
	case "bytes":
		chg.BytesValue = nil
	}
}
