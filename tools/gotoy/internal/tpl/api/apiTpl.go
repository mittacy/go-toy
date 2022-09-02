package api

import (
	"bytes"
	"github.com/mittacy/go-toy/tools/gotoy/internal/base"
	"html/template"
)

var apiTemplate = `
{{- /* delete empty line */ -}}
package api

import (
	"github.com/gin-gonic/gin"
	"{{ .AppName }}/{{ .TargetDir }}/dp"
	"{{ .AppName }}/{{ .TargetDir }}/service"
	"{{ .AppName }}/{{ .TargetDir }}/validator/{{ .NameLower }}Vdr"
	"github.com/mittacy/go-toy/core/log"
	"github.com/mittacy/go-toy/core/response"
	"github.com/mittacy/go-toy/core/singleton"
)

var {{ .Name }} {{ .NameLower }}Api

type {{ .NameLower }}Api struct {
	dp dp.{{ .Name }}
}

func init() {
	singleton.Register(func() {
		{{ .Name }} = {{ .NameLower }}Api{
			dp: dp.New{{ .Name }}(),
		}
	})
}

func (ctl *{{ .NameLower }}Api) List(c *gin.Context) {
	req := {{ .NameLower }}Vdr.ListReq{}
	if err := c.ShouldBindQuery(&req); err != nil {
		response.ValidateErr(c, err)
		return
	}

	data, err := service.{{ .Name }}.List(c, req.Page, req.PageSize)
	if err != nil {
		response.FailCheckBizErr(c, log.Default(), req, "{{ .Name }} list err", err)
		return
	}

	res := ctl.dp.ListReply(data)
	response.Success(c, res)
}
`

type Api struct {
	AppName   string
	Name      string
	NameLower string
	TargetDir string
}

func (s *Api) execute() ([]byte, error) {
	s.Name = base.StringFirstUpper(s.Name)
	s.NameLower = base.StringFirstLower(s.Name)

	buf := new(bytes.Buffer)

	tmpl, err := template.New("api").Parse(apiTemplate)
	if err != nil {
		return nil, err
	}

	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
