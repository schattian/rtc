package schema

import (
	"fmt"

	"github.com/sebach1/git-crud/integrity"
	"github.com/sebach1/git-crud/schema/valide"
)

// A Column is the representation of SQL column which defines the structure of the fields that is contains.
type Column struct {
	Name      integrity.ColumnName `json:"name,omitempty"`
	Validator integrity.Validator
	Type      integrity.ValueType `json:"type,omitempty"`
}

// Validate wraps the column validator func and returns its result
func (c *Column) Validate(val interface{}) error {
	if c.Validator == nil {
		return nil
	}
	err := c.Validator(val)
	if err != nil {
		return err
	}
	return nil
}

func (c *Column) unaliasType() {
	switch c.Type {
	case "json":
		c.Type = "json.RawMessage"
	case "bytes":
		c.Type = "[]byte"
	}
}

// Assigns the appropiated builtin validator (on schema/valide pkg) given the Column.Type
func (c *Column) applyBuiltinValidator() error {
	c.unaliasType()
	switch c.Type {
	case "string":
		c.Validator = valide.String
	case "int":
		c.Validator = valide.Int
	case "float32":
		c.Validator = valide.Float32
	case "float64":
		c.Validator = valide.Float64
	case "json.RawMessage":
		c.Validator = valide.JSON
	case "[]byte":
		c.Validator = valide.Bytes
	case "":
		return fmt.Errorf("the TYPE of COLUMN %v is NIL", c.Name)
	default:
		return fmt.Errorf("the TYPE %v of COLUMN %v is INVALID", c.Type, c.Name)
	}
	return nil
}
