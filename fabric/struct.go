package fabric

import "text/template"

const openStruct = `
package {{.SchemaName}}

// {{.Name}} is the native representation of the {{.Name}} resource in {{.SchemaName}} schema
type {{.Name}} struct {
	{{$msh := .Marshal}}
	{{range $field := .Fields}}
		{{$field.Name}}   {{$field.Type}}   ` + "`{{$msh}}:\"{{$field.Tag}}\"`" + `
	{{end}}
}
`

var structTemplate = template.Must(template.New("structTemplate").Parse(openStruct))
