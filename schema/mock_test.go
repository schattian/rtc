package schema

import "github.com/sebach1/git-crud/integrity"

func (c *Column) addValidator(validator integrity.Validator) *Column {
	c.Validator = validator
	return c
}

func (sch *Schema) addColValidator(colName integrity.ColumnName, validator integrity.Validator) *Schema {
	for _, table := range sch.Blueprint {
		for _, col := range table.Columns {
			if col.Name == colName {
				col.addValidator(validator)
			}
		}
	}
	return sch
}
