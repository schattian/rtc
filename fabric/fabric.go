package fabric

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/sebach1/git-crud/integrity"
	"github.com/sebach1/git-crud/schema"
)

// Fabric creates native go struct types given a schema
type Fabric struct {
	Schema *schema.Schema

	Dir string
	wg  *sync.WaitGroup
}

// Produce is the main Fabric wrapper
func (f *Fabric) Produce(marshal string) error {
	err := f.validate()
	if err != nil {
		return err
	}
	f.Schema.Name = integrity.SchemaName(strings.ToLower(string(f.Schema.Name)))
	f.Dir = fmt.Sprintf("fabric/%v", f.Schema.Name)
	err = os.MkdirAll(f.Dir, os.ModePerm)
	if err != nil {
		return err
	}

	f.wg = new(sync.WaitGroup)
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
	err := structTemplate.Execute(&out, tableStruct)
	filename := fmt.Sprintf("%v/%v.go", f.Dir, strings.ToLower(string(table.Name)))
	generated := out.Bytes()
	generated, err = format.Source(generated)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(filename, generated, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

func (f *Fabric) structFromTable(table *schema.Table, marshal string) *tableData {
	tableStruct := &tableData{
		SchemaName: string(f.Schema.Name),
		Name:       toCamelCase(string(table.Name), true),
		Marshal:    toSnakeCase(marshal, '_'),
	}
	for _, col := range table.Columns {
		tableStruct.Fields = append(tableStruct.Fields, fieldFromColumn(col))
	}
	return tableStruct
}

func fieldFromColumn(col *schema.Column) *columnData {
	return &columnData{
		Name: toCamelCase(string(col.Name), true),
		Type: col.Type,
		Tag:  toSnakeCase(string(col.Name), '_'),
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

type tableData struct {
	SchemaName string
	Name       string
	Fields     []*columnData
	Marshal    string
}

type columnData struct {
	Name string
	Type integrity.ValueType
	Tag  string
}
