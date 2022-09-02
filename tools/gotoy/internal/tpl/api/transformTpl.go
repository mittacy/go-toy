package api

import (
	"bytes"
	"github.com/mittacy/go-toy/tools/gotoy/internal/base"
	"html/template"
)

var transformTemplate = `
{{- /* delete empty line */ -}}
package dp

import (
	"{{ .AppName }}/{{ .TargetDir }}/service/smodel"
	"{{ .AppName }}/{{ .TargetDir }}/validator/{{ .NameLower }}Vdr"
)

type {{ .Name }} struct{}

func New{{ .Name }}() {{ .Name }} {
	return {{ .Name }}{}
}

// 列表数据响应封装
func (ctl *{{ .Name }}) ListReply(data *smodel.{{ .Name }}List) map[string]interface{} {
	{{ .NameLower }}s := make([]{{ .NameLower }}Vdr.ListReply{{ .Name }}, len(data.{{ .Name }}s))
	for i, v := range data.{{ .Name }}s {
		{{ .NameLower }}s[i].Id = v.Id
	}

	return map[string]interface{}{
		"{{ .NameLower }}":  {{ .NameLower }}s,
		"total": data.Total,
	}
}
`

type Transform struct {
	Name      string
	NameLower string
	AppName   string
	TargetDir string
}

func (s *Transform) execute() ([]byte, error) {
	s.Name = base.StringFirstUpper(s.Name)
	s.NameLower = base.StringFirstLower(s.Name)

	buf := new(bytes.Buffer)

	tmpl, err := template.New("validator").Parse(transformTemplate)
	if err != nil {
		return nil, err
	}

	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
