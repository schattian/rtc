package fabric

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"sync"

	"github.com/sebach1/git-crud/integrity"
	"github.com/sebach1/git-crud/schema"
)

// Fabric creates native go struct types given a schema
type Fabric struct {
	Schema *schema.Schema

	dir string
	wg  *sync.WaitGroup
}

// Produce is the main Fabric wrapper
func (f *Fabric) Produce(marshal string) error {
	err := f.validate()
	if err != nil {
		panic(err)
	}
	f.dir = fmt.Sprintf("fabric/%v", f.Schema.Name)
	err = os.MkdirAll(f.dir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	for _, table := range f.Schema.Blueprint {
		f.wg.Add(1)
		go f.writeStructFromTable(table, marshal)
	}
	f.wg.Wait()
	return nil
}

func (f *Fabric) writeStructFromTable(table *schema.Table, marshal string) {
	defer f.wg.Done()
	var out bytes.Buffer
	tableStruct := f.structFromTable(table, marshal)
	structTemplate.Execute(bufio.NewWriter(&out), tableStruct)
	filename := fmt.Sprintf("%v/%v.go", f.dir, strings.ToLower(string(table.Name)))
	generated := out.Bytes()
	generated, err := format.Source(generated)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(filename, generated, os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func (f *Fabric) structFromTable(table *schema.Table, marshal string) *tableData {
	tableStruct := &tableData{
		SchemaName: f.Schema.Name,
		Name:       table.Name,
		Marshal:    marshal,
	}
	for _, col := range table.Columns {
		tableStruct.Fields = append(tableStruct.Fields, fieldFromColumn(col))
	}
	return tableStruct
}

func fieldFromColumn(col *schema.Column) *columnData {
	return &columnData{
		Name: col.Name,
		Type: col.Validator.NativeType(),
		Tag:  toSnakeCase(string(col.Name)),
	}
}

func (f *Fabric) validate() error {
	if f.Schema == nil {
		return errors.New("the SCHEMA cant be NIL")
	}
	if f.Schema.Name == "" {
		return errors.New("the SCHEMA NAME cant be NIL")
	}
	return nil
}

func toSnakeCase(str string) string {
	matchFirstCap := regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

type tableData struct {
	SchemaName integrity.SchemaName
	Name       integrity.TableName
	Fields     []*columnData
	Marshal    string
}

type columnData struct {
	Name integrity.ColumnName
	Type integrity.ValueType
	Tag  string
}
