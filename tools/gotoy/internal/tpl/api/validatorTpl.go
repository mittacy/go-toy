package api

import (
	"bytes"
	"github.com/mittacy/go-toy/tools/gotoy/internal/base"
	"html/template"
	"strings"
)

var validatorTemplate = `
{{- /* delete empty line */ -}}
package {{ .NameLower }}Vdr

type ListReq struct {
	Page     int ${backquote}form:"page" json:"page" binding:"required,min=1"${backquote}
	PageSize int ${backquote}form:"page_size" json:"page_size" binding:"required,min=1,max=100"${backquote}
}
type ListReply{{ .Name }} struct {
	Id        int64  ${backquote}json:"id"${backquote}
}
`

type Validator struct {
	Name      string
	NameLower string
}

func (s *Validator) execute() ([]byte, error) {
	s.Name = base.StringFirstUpper(s.Name)
	s.NameLower = base.StringFirstLower(s.Name)

	buf := new(bytes.Buffer)

	validatorTemplate = strings.Replace(validatorTemplate, "${backquote}", "`", -1)
	tmpl, err := template.New("validator").Parse(validatorTemplate)
	if err != nil {
		return nil, err
	}

	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
