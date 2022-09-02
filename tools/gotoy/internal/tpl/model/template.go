package model

import (
	"bytes"
	"github.com/mittacy/go-toy/tools/gotoy/internal/base"
	"text/template"
)

var modelTemplate = `
{{- /* delete empty line */ -}}
package model

const (
	{{ .Name }}IsDeletedNo  = 0
	{{ .Name }}IsDeletedYes = 1
)

type {{ .Name }} struct {
	Id int64
	IsDeleted int8
}

func (*{{ .Name }}) TableName() string {
	return "table_name"
}

`

type Model struct {
	Name      string
	NameLower string
}

func (s *Model) execute() ([]byte, error) {
	s.Name = base.StringFirstUpper(s.Name)
	s.NameLower = base.StringFirstLower(s.Name)

	buf := new(bytes.Buffer)

	tmpl, err := template.New("model").Parse(modelTemplate)
	if err != nil {
		return nil, err
	}

	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
