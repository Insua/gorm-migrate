package gorm_migrate

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

func (m *{{ .StructName }}) Up() {
	
}

func (m *{{ .StructName }}) Down()  {

}

`
