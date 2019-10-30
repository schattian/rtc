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
		name string
		args args
		want error
	}{
		{
			name: "column is in the schema",
			args: args{sch: gSchemas.Basic, colName: gColumns.Basic.Name},
			want: errForeignColumn,
		},
		{
			name: "column doesn't exists in the schema",
			args: args{sch: gSchemas.Basic, colName: gColumns.Rare.Name},
			want: errNonexistentColumn,
		},
		{
			name: "schema caller is zero-valued",
			args: args{sch: gSchemas.Zero, colName: gColumns.Basic.Name},
			want: errNonexistentColumn,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if err := tt.args.sch.preciseColErr(tt.args.colName); err != tt.want {
				t.Errorf("Schema.preciseColErr() error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestSchema_Validate(t *testing.T) {
	t.Parallel()
	type args struct {
		tableName   integrity.TableName
		colName     integrity.ColumnName
		optionKeys  []integrity.OptionKey
		val         interface{}
		helperScope *Planisphere
	}
	tests := []struct {
		name    string
		sch     *Schema
		args    args
		wantErr bool
	}{
		{
			name: "passes the validations",
			sch:  gSchemas.Basic,
			args: args{
				tableName:   gTables.Basic.Name,
				colName:     gColumns.Basic.Name,
				optionKeys:  []integrity.OptionKey{gTables.Basic.OptionKeys[0]},
				helperScope: &Planisphere{gSchemas.Basic},
			},
			wantErr: false,
		},
		{
			name: "optionKey nonexistant",
			sch:  gSchemas.Basic,
			args: args{
				tableName:   gTables.Basic.Name,
				colName:     gColumns.Basic.Name,
				optionKeys:  []integrity.OptionKey{gTables.Rare.OptionKeys[0]},
				helperScope: &Planisphere{gSchemas.Basic},
			},
			wantErr: true,
		},
		{
			name: "value doesn't pass the column validator func",
			sch:  gSchemas.Basic.Copy().addColValidator(gColumns.Basic.Name, valide.String),
			args: args{
				tableName:   gTables.Basic.Name,
				colName:     gColumns.Basic.Name,
				val:         3,
				optionKeys:  []integrity.OptionKey{},
				helperScope: &Planisphere{gSchemas.Basic},
			},
			wantErr: true,
		},
		{
			name: "column not inside any table",
			sch:  gSchemas.Basic,
			args: args{
				tableName:   gTables.Basic.Name,
				colName:     gColumns.Rare.Name,
				optionKeys:  []integrity.OptionKey{gTables.Basic.OptionKeys[0]},
				helperScope: &Planisphere{gSchemas.Basic},
			},
			wantErr: true,
		},
		{
			name: "table nonexistant",
			sch:  gSchemas.Basic,
			args: args{
				tableName:   gTables.Rare.Name,
				colName:     gColumns.Basic.Name,
				optionKeys:  []integrity.OptionKey{},
				helperScope: &Planisphere{gSchemas.Basic},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			wg := new(sync.WaitGroup)
			errCh := make(chan error, 1)
			wg.Add(1)
			go tt.sch.Validate(tt.args.tableName, tt.args.colName, tt.args.optionKeys, tt.args.val, tt.args.helperScope, wg, errCh)
			wg.Wait()
			isErrored := len(errCh) == 1
			if isErrored && !tt.wantErr {
				err := <-errCh
				t.Errorf("Schema.Validate() error: %v; wantErr %v", err, tt.wantErr)
			}
			if !isErrored && tt.wantErr {
				t.Errorf("Schema.Validate() error: %v; wantErr %v", nil, tt.wantErr)
			}
		})
	}
}

func TestSchema_WrapValidateSelf(t *testing.T) {
	t.Parallel()
	fuzzyTests := []*schemaVary{
		// Schema
		&schemaVary{
			name:     "sch nil name",
			function: func(sch *Schema) *Schema { sch.Name = ""; return sch },
			err:      errNilSchemaName},
		// Table
		&schemaVary{
			name:     "table nil name",
			function: func(sch *Schema) *Schema { sch.Blueprint[0].Name = ""; return sch },
			err:      errNilTableName},
		// Column
		&schemaVary{
			name:     "col nil type",
			function: func(sch *Schema) *Schema { sch.Blueprint[0].Columns[0].Type = ""; return sch },
			err:      errNilColumnType},
		&schemaVary{
			name:     "col nil name",
			function: func(sch *Schema) *Schema { sch.Blueprint[0].Columns[0].Name = ""; return sch },
			err:      errNilColumnName},
	}
	normalTests := []*schemaVary{
		// Schema
		&schemaVary{
			name:     "sch nil",
			function: func(sch *Schema) *Schema { return nil },
			err:      errNilSchema},
		&schemaVary{
			name:     "sch nil bp",
			function: func(sch *Schema) *Schema { sch.Blueprint = nil; return sch },
			err:      errNilBlueprint},
		// Table
		&schemaVary{
			name:     "table nil",
			function: func(sch *Schema) *Schema { sch.Blueprint[0] = nil; return sch },
			err:      errNilTable},
		&schemaVary{
			name:     "table nil columns",
			function: func(sch *Schema) *Schema { sch.Blueprint[0].Columns = nil; return sch },
			err:      errNilColumns},
		// Column
		&schemaVary{
			name:     "col nil",
			function: func(sch *Schema) *Schema { sch.Blueprint[0].Columns[0] = nil; return sch },
			err:      errNilColumn},
	}
	for _, tt := range normalTests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			sch := tt.function(gSchemas.Basic.Copy())
			err := sch.WrapValidateSelf()
			unwrappedVErr := integrity.UnwrapValidationError(err, false)
			diffGot, diffWant := diffBetweenErrs(unwrappedVErr, []error{tt.err})
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
		sch := gSchemas.Basic.Copy()
		for _, tt := range test {
			sch = tt.function(sch)
			wantErrs = append(wantErrs, tt.err)
		}
		t.Run(strconv.Itoa(k), func(t *testing.T) {
			t.Parallel()
			err := sch.WrapValidateSelf()
			unwrappedVErr := integrity.UnwrapValidationError(err, false)
			diffGot, diffWant := diffBetweenErrs(unwrappedVErr, wantErrs)
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
