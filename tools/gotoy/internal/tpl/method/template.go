package method

import (
	"bytes"
	"text/template"
)

const (
	validatorTpl = `
type {{ .Method }}Req struct {}
type {{ .Method }}Reply struct {}
`
	dpTpl = `
// {{ .Comment }}
func (ctl *{{ .Name }}) {{ .Method }}(data *smodel.{{ .Name }}{{ .Method }}) map[string]interface{} {
	return map[string]interface{}{
		"hi": "hello",
	}
}
`
	apiTpl = `
// {{ .Comment }}
func (ctl *{{ .NameLower }}Api) {{ .Method }}(c *gin.Context) {
	req := {{ .NameLower }}Vdr.{{ .Method }}Req{}
	if err := c.ShouldBind(&req); err != nil {
		response.ValidateErr(c, err)
		return
	}

	serviceRes, err := service.{{ .Name }}.{{ .Method }}(c, req)
	if err != nil {
		response.FailCheckBizErr(c, log.Default(), req, "{{ .Name }} {{ .Method }} err", err)
		return
	}

	res := ctl.dp.{{ .Method }}(serviceRes)
	response.Success(c, res)
}
`
	serviceTpl = `
// {{ .Comment }}
func (ctl *{{ .NameLower }}Service) {{ .Method }}(c context.Context, req {{ .NameLower }}Vdr.{{ .Method }}Req) (*smodel.{{ .Name }}{{ .Method }}, error) {
	// 并发查询
	var (
		errMsg string
	)
	eg := errgroup.WithCancel(context.Background())
	eg.Go(func(ctx context.Context) error {
		var err error
		if err != nil {
			errMsg = "查询XXX错误"
		}
		return err
	})
	if err := eg.Wait(); err != nil {
		return nil, errors.WithMessage(err, errMsg)
	}

	// 串行查询

	// 返回数据库数据

	return &smodel.{{ .Name }}{{ .Method }}{}, nil
}
`
	sModelTpl = `
type {{ .Name }}{{ .Method }} struct {
}
`
)

type Method struct {
	Name      string
	NameLower string
	Method    string
	Comment   string
}

func (s *Method) execute(tpl string) ([]byte, error) {
	buf := new(bytes.Buffer)

	tmpl, err := template.New(tpl).Parse(tpl)
	if err != nil {
		return nil, err
	}

	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
