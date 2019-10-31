package fabric

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"strings"
	"sync"

	"github.com/sebach1/git-crud/integrity"
	"github.com/sebach1/git-crud/schema"
	"github.com/spf13/afero"
)

// Fabric creates native go struct types given a schema
type Fabric struct {
	Schema *schema.Schema
	Dir    string

	wg *sync.WaitGroup
	// FileSystem-related
	fsWg    *sync.WaitGroup
	fsSmp   chan int // Semaphore
	fsErrCh chan error
}

// Produce is the main Fabric wrapper
func (f *Fabric) Produce(marshal string, fs afero.Fs) (err error) {
	err = f.Schema.ValidateSelf()
	if err != nil {
		return
	}

	f.Schema = f.Schema.Copy()
	f.Schema.Name = integrity.SchemaName(strings.ToLower(string(f.Schema.Name)))
	if f.Dir == "" {
		f.Dir = fmt.Sprintf("fabric/%v", f.Schema.Name)
	}

	f.wg = &sync.WaitGroup{}
	f.initFsConcurrency()

	for _, table := range f.Schema.Blueprint {
		f.wg.Add(1)
		go f.mkStructFileFromTable(table, marshal, fs)
	}
	f.wg.Wait()

	if len(f.fsErrCh) > 0 {
		return <-f.fsErrCh
	}

	f.fsWg.Add(1)
	go f.writeAll(fs)
	f.fsWg.Wait()
	if len(f.fsErrCh) > 0 {
		return <-f.fsErrCh
	}
	return nil
}

func (f *Fabric) initFsConcurrency() {
	f.fsWg = &sync.WaitGroup{}
	f.fsSmp = make(chan int, 1)
	f.fsErrCh = make(chan error, 1)
}

func (f *Fabric) mkStructFileFromTable(table *schema.Table, marshal string, fs afero.Fs) (err error) {
	defer f.wg.Done()
	var out bytes.Buffer
	tableStruct := f.structFromTable(table, marshal)
	err = structTemplate.Execute(&out, tableStruct)
	if err != nil {
		return err
	}
	filename := fmt.Sprintf("%v/%v.go", f.Dir, strings.ToLower(string(table.Name)))
	generated := out.Bytes()
	generated, err = format.Source(generated)
	if err != nil {
		return err
	}
	f.fsWg.Add(1)
	go f.writeFile(fs, filename, generated)
	return nil
}

func (f *Fabric) writeAll(fs afero.Fs) {
	defer f.fsWg.Done()
	err := fs.MkdirAll(f.Dir, os.ModePerm)
	if err != nil {
		f.fsErrCh <- err
	}
	f.fsSmp <- 0 // Unblocks any writeFile
}

func (f *Fabric) writeFile(fs afero.Fs, filename string, generated []byte) {
	defer f.fsWg.Done()
	for {
		select {
		case <-f.fsSmp:
			err := afero.WriteFile(fs, filename, generated, os.ModePerm)
			if err != nil {
				f.fsErrCh <- err
			}
			f.fsSmp <- 0 // Unblocks the next any writeFile
			return
		}
	}
}

func (f *Fabric) structFromTable(table *schema.Table, marshal string) *tableData {
	tableStruct := &tableData{
		SchemaName: string(f.Schema.Name),
		Name:       toCamelCase(string(table.Name)),
		Marshal:    toSnakeCase(marshal, '_'),
	}
	for _, col := range table.Columns {
		field := fieldFromColumn(col)
		tableStruct.Fields = append(tableStruct.Fields, field)
	}
	return tableStruct
}

func fieldFromColumn(col *schema.Column) *columnData {
	return &columnData{
		Name: toCamelCase(string(col.Name)),
		Type: col.Type,
		Tag:  toSnakeCase(string(col.Name), '_'),
	}
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
