package migrate

var fileContent = `package {{ .Package }}

type {{ .StructName }} struct {
	Name string
}

func init()  {
	m := {{ .StructName }}{
		Name: "{{ .FileName }}",
	}
	Migrations = append(Migrations, &m)
}

func (m *{{ .StructName }}) Up() error {

	return nil
}

func (m *{{ .StructName }}) Down() error  {

	return nil
}
`

var migrationContent = `package {{ .Package }}

var Migrations = make([]interface{}, 0)
`
