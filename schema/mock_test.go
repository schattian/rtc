package schema

import "github.com/sebach1/git-crud/internal/integrity"

func (c *Column) copy() *Column {
	newCol := new(Column)
	*newCol = *c
	return newCol
}

func (sch *Schema) copy() *Schema {
	newSch := new(Schema)
	*newSch = *sch
	var newBlueprint []*Table
	for _, table := range newSch.Blueprint {
		newBlueprint = append(newBlueprint, table.copy())
	}
	newSch.Blueprint = newBlueprint
	return newSch
}

func (t *Table) copy() *Table {
	newTab := new(Table)
	*newTab = *t
	var newCols []*Column
	for _, col := range newTab.Columns {
		newCols = append(newCols, col.copy())
	}
	newTab.Columns = newCols
	return newTab
}

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
