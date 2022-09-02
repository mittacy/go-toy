package service

import (
	"bytes"
	"github.com/mittacy/go-toy/tools/gotoy/internal/base"
	"text/template"
)

var modelTemplate = `
{{- /* delete empty line */ -}}
package smodel

import (
	"{{ .AppName }}/{{ .TargetDir }}/model"
)

type {{ .Name }}List struct {
	{{ .Name }}s []model.{{ .Name }}
	Total int64
}
`

type Model struct {
	AppName   string
	Name      string
	NameLower string
	TargetDir string
}

func (s *Model) execute() ([]byte, error) {
	s.Name = base.StringFirstUpper(s.Name)
	s.NameLower = base.StringFirstLower(s.Name)

	buf := new(bytes.Buffer)

	tmpl, err := template.New("smodel").Parse(modelTemplate)
	if err != nil {
		return nil, err
	}

	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
