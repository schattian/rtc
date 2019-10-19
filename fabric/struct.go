package fabric

import "text/template"

const openStruct = `
package {{.SchemaName}}

// {{.Name}} is the native representation of the {{.Name}} resource in {{.SchemaName}} schema
type {{.Name}} struct {
	{{range $field := .Fields}}
		{{$field.Name}}   {{$field.Type}}   ` + "`{{.Marshal}}:\"{{$field.Tag}}\"`" + `
	{{end}}
}
`

var structTemplate = template.Must(template.New("structTemplate").Parse(openStruct))
