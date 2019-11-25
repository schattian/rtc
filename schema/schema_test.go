package schema

import (
	"strconv"
	"sync"
	"testing"

	"github.com/sebach1/git-crud/integrity"
	"github.com/sebach1/git-crud/schema/valide"
)

func TestSchema_preciseColErr(t *testing.T) {
	t.Parallel()
	type args struct {
		sch     *Schema
		colName integrity.ColumnName
	}
	tests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "column is in the schema",
			args:    args{sch: gSchemas.Foo, colName: gColumns.Foo.Name},
			wantErr: errForeignColumn,
		},
		{
			name:    "column doesn't exists in the schema",
			args:    args{sch: gSchemas.Foo, colName: gColumns.Bar.Name},
			wantErr: errNonexistentColumn,
		},
		{
			name:    "schema caller is zero-valued",
			args:    args{sch: gSchemas.Zero, colName: gColumns.Foo.Name},
			wantErr: errNonexistentColumn,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.args.sch.preciseColErr(tt.args.colName); err != tt.wantErr {
				t.Errorf("Schema.preciseColErr() error = %v, want %v", err, tt.wantErr)
			}
		})
	}
}

func TestSchema_ValidateCtx(t *testing.T) {
	t.Parallel()
	type args struct {
		tableName   integrity.TableName
		colName     integrity.ColumnName
		optionKeys  []integrity.OptionKey
		val         interface{}
		helperScope *Planisphere
	}
	tests := []struct {
		name     string
		sch      *Schema
		args     args
		wantsErr bool
	}{
		{
			name: "passes the validations",
			sch:  gSchemas.Foo,
			args: args{
				tableName:   gTables.Foo.Name,
				colName:     gColumns.Foo.Name,
				optionKeys:  []integrity.OptionKey{gTables.Foo.OptionKeys[0]},
				helperScope: &Planisphere{gSchemas.Foo},
			},
			wantsErr: false,
		},
		{
			name: "but NO GIVEN COL passes validation",
			sch:  gSchemas.Foo,
			args: args{
				tableName:   gTables.Foo.Name,
				colName:     "",
				optionKeys:  []integrity.OptionKey{gTables.Foo.OptionKeys[0]},
				helperScope: &Planisphere{gSchemas.Foo},
			},
			wantsErr: false,
		},
		{
			name: "optionKey nonexistant",
			sch:  gSchemas.Foo,
			args: args{
				tableName:   gTables.Foo.Name,
				colName:     gColumns.Foo.Name,
				optionKeys:  []integrity.OptionKey{gTables.Bar.OptionKeys[0]},
				helperScope: &Planisphere{gSchemas.Foo},
			},
			wantsErr: true,
		},
		{
			name: "value doesn't pass the column validator func",
			sch:  gSchemas.Foo.Copy().addColValidator(gColumns.Foo.Name, valide.String),
			args: args{
				tableName:   gTables.Foo.Name,
				colName:     gColumns.Foo.Name,
				val:         3,
				optionKeys:  []integrity.OptionKey{},
				helperScope: &Planisphere{gSchemas.Foo},
			},
			wantsErr: true,
		},
		{
			name: "column not inside any table",
			sch:  gSchemas.Foo,
			args: args{
				tableName:   gTables.Foo.Name,
				colName:     gColumns.Bar.Name,
				optionKeys:  []integrity.OptionKey{gTables.Foo.OptionKeys[0]},
				helperScope: &Planisphere{gSchemas.Foo},
			},
			wantsErr: true,
		},
		{
			name: "table nonexistant",
			sch:  gSchemas.Foo,
			args: args{
				tableName:   gTables.Bar.Name,
				colName:     gColumns.Foo.Name,
				optionKeys:  []integrity.OptionKey{},
				helperScope: &Planisphere{gSchemas.Foo},
			},
			wantsErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			wg := new(sync.WaitGroup)
			errCh := make(chan error, 1)
			wg.Add(1)
			go tt.sch.ValidateCtx(tt.args.tableName, tt.args.colName, tt.args.optionKeys, tt.args.val, tt.args.helperScope, wg, errCh)
			wg.Wait()
			isErrored := len(errCh) == 1
			if isErrored && !tt.wantsErr {
				err := <-errCh
				t.Errorf("Schema.ValidateCtx() error: %v; wantErr %v", err, tt.wantsErr)
			}
			if !isErrored && tt.wantsErr {
				t.Errorf("Schema.ValidateCtx() error: %v; wantErr %v", nil, tt.wantsErr)
			}
		})
	}
}

func TestSchema_WrapValidateSelf(t *testing.T) {
	t.Parallel()
	fuzzyTests := []*schemaVary{
		// Schema
		{
			name:     "sch nil name",
			function: func(sch *Schema) *Schema { sch.Name = ""; return sch },
			err:      errNilSchemaName},
		// Table
		{
			name:     "table nil name",
			function: func(sch *Schema) *Schema { sch.Blueprint[0].Name = ""; return sch },
			err:      errNilTableName},
		// Column
		{
			name:     "col nil type",
			function: func(sch *Schema) *Schema { sch.Blueprint[0].Columns[0].Type = ""; return sch },
			err:      errNilColumnType},
		{
			name:     "col nil name",
			function: func(sch *Schema) *Schema { sch.Blueprint[0].Columns[0].Name = ""; return sch },
			err:      errNilColumnName},
	}
	normalTests := []*schemaVary{
		// Schema
		{
			name:     "sch nil",
			function: func(sch *Schema) *Schema { return nil },
			err:      errNilSchema},
		{
			name:     "sch nil bp",
			function: func(sch *Schema) *Schema { sch.Blueprint = nil; return sch },
			err:      errNilBlueprint},
		// Table
		{
			name:     "table nil",
			function: func(sch *Schema) *Schema { sch.Blueprint[0] = nil; return sch },
			err:      errNilTable},
		{
			name:     "table nil columns",
			function: func(sch *Schema) *Schema { sch.Blueprint[0].Columns = nil; return sch },
			err:      errNilColumns},
		// Column
		{
			name:     "col nil",
			function: func(sch *Schema) *Schema { sch.Blueprint[0].Columns[0] = nil; return sch },
			err:      errNilColumn},
	}
	for _, tt := range normalTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			sch := tt.function(gSchemas.Foo.Copy())
			errs := sch.ValidateSelf().UnwrapAll(integrity.UnwrapValidationError)
			diffGot, diffWant := diffBetweenErrs(errs, []error{tt.err})
			if diffGot != nil {
				t.Errorf("Schema.WrapValidate() GOT more errs than want: %v", diffGot)
			}
			if diffWant != nil {
				t.Errorf("Schema.WrapValidate() WANT more errs than got: %v", diffWant)
			}
		})
	}
	for k, test := range fuzzyFactorial(fuzzyTests) {
		test := test
		var wantErrs []error
		sch := gSchemas.Foo.Copy()
		for _, tt := range test {
			sch = tt.function(sch)
			wantErrs = append(wantErrs, tt.err)
		}
		t.Run(strconv.Itoa(k), func(t *testing.T) {
			t.Parallel()
			errs := sch.ValidateSelf().UnwrapAll(integrity.UnwrapValidationError)
			diffGot, diffWant := diffBetweenErrs(errs, wantErrs)
			if diffGot != nil {
				t.Errorf("Schema.WrapValidate() GOT more errs than want: %v", diffGot)
			}
			if diffWant != nil {
				t.Errorf("Schema.WrapValidate() WANT more errs than got: %v", diffWant)
			}
		})
	}
}

type schemaVary struct {
	name     string
	function func(sch *Schema) *Schema
	err      error
}

func fuzzyFactorial(set []*schemaVary) (subsets [][]*schemaVary) {
	length := uint(len(set))
	for subsetBits := 1; subsetBits < (1 << length); subsetBits++ {
		var subset []*schemaVary
		for object := uint(0); object < length; object++ {
			if (subsetBits>>object)&1 == 1 {
				subset = append(subset, set[object])
			}
		}
		subsets = append(subsets, subset)
	}
	return subsets
}

func strErrIsInArray(err error, errs []error) bool {
	for _, otherErr := range errs {
		if otherErr.Error() == err.Error() {
			return true
		}
	}
	return false
}

func diffBetweenErrs(errs []error, in []error) (errsOnly []error, inOnly []error) {
	for _, err := range errs {
		if !strErrIsInArray(err, in) {
			errsOnly = append(errsOnly, err)
		}
	}
	for _, otherErr := range in {
		if !strErrIsInArray(otherErr, errs) {
			inOnly = append(inOnly, otherErr)
		}
	}
	return
}
