package service

import (
	"bytes"
	"github.com/mittacy/go-toy/tools/gotoy/internal/base"
	"text/template"
)

var serviceTemplate = `
{{- /* delete empty line */ -}}
package service

import (
	"context"
	"{{ .AppName }}/{{ .TargetDir }}/data"
	"{{ .AppName }}/{{ .TargetDir }}/service/smodel"
	"github.com/mittacy/go-toy/core/singleton"
	"github.com/pkg/errors"
)

// 一般情况下service应该只引用并控制自己的data模型，需要其他服务的功能请service.Xxx调用服务而不是引入其他data模型

// {{ .Name }} 服务说明注释
var {{ .Name }} {{ .NameLower }}Service

type {{ .NameLower }}Service struct {
	data data.{{ .Name }}
}

func init() {
	singleton.Register(func() {
		{{ .Name }} = {{ .NameLower }}Service{
			data: data.New{{ .Name }}(),
		}
	})
}

func (ctl *{{ .NameLower }}Service) List(c context.Context, page, pageSize int) (*smodel.{{ .Name }}List, error) {
	{{ .NameLower }}, total, err := ctl.data.List(c, page, pageSize)
	if err != nil {
		return nil, errors.WithMessage(err, "查询列表错误")
	}

	return &smodel.{{ .Name }}List{
		{{ .Name }}s: {{ .NameLower }},
		Total: total,
	}, nil
}
`

type Service struct {
	AppName   string
	Name      string
	NameLower string
	TargetDir string
}

func (s *Service) execute() ([]byte, error) {
	s.Name = base.StringFirstUpper(s.Name)
	s.NameLower = base.StringFirstLower(s.Name)

	buf := new(bytes.Buffer)

	tmpl, err := template.New("service").Parse(serviceTemplate)
	if err != nil {
		return nil, err
	}

	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
