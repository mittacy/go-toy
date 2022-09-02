package data

import (
	"bytes"
	"github.com/mittacy/go-toy/tools/gotoy/internal/base"
	"text/template"
)

var dataTemplate = `
{{- /* delete empty line */ -}}
package data

import (
	"context"
	"{{ .AppName }}/{{ .TargetDir }}/model"
	{{- if .InjectMongo}}
	"github.com/mittacy/go-toy/core/eMongo"
	{{- end}}
	{{- if .InjectMysql}}
	"github.com/mittacy/go-toy/core/eMysql"
	{{- end}}
	{{- if .InjectRedis}}
	"github.com/mittacy/go-toy/core/eRedis"
	{{- end}}
	{{- if .InjectHttp}}
	"github.com/mittacy/go-toy/core/goHttp"
	{{- end}}
	{{- if or .InjectRedis .InjectHttp}}
	"github.com/spf13/viper"
	{{- end}}
)

type {{ .Name }} struct {
	{{- if .InjectMysql}}
	eMysql.EGorm
	{{- end}}
	{{- if .InjectMongo}}
	eMongo.EMongo
	{{- end}}
	{{- if .InjectRedis}}
	eRedis.ERedis
	{{- end}}
	{{- if .InjectHttp}}
	goHttp.ApiServer
	{{- end}}
}

func New{{ .Name }}() {{ .Name }} {
	return {{ .Name }}{
		{{- if .InjectMysql}}
		EGorm:  eMysql.EGorm{ConfName: "localhost"},
		{{- end}}
		{{- if .InjectMongo}}
		EMongo: eMongo.EMongo{ConfName: "localhost", Collection: "collection_name"},
		{{- end}}
		{{- if .InjectRedis}}
		ERedis: eRedis.ERedis{ConfName: "localhost", DB: 0},
		{{- end}}
		{{- if .InjectHttp}}
		ApiServer: goHttp.NewApiServer("host:port"),
		{{- end}}
	}
}

func (ctl *{{ .Name }}) List(c context.Context, page, pageSize int) ([]model.{{ .Name }}, int64, error) {
	{{ .NameLower }} := []model.{{ .Name }}{
		{
			Id:        1,
			IsDeleted: model.{{ .Name }}IsDeletedNo,
		},
		{
			Id:        2,
			IsDeleted: model.{{ .Name }}IsDeletedNo,
		},
	}

	return {{ .NameLower }}, 2, nil
}

{{if .InjectRedis}}
/*
 * 以下为查询缓存KEY方法
 */
func (ctl *{{ .Name }}) cacheKeyPre() string {
	return viper.GetString("APP_NAME") + ":{{ .NameLower }}"
}
{{- end}}
`

type Data struct {
	AppName     string
	Name        string
	NameLower   string
	TargetDir   string
	InjectMysql bool
	InjectMongo bool
	InjectRedis bool
	InjectHttp  bool
}

func (s *Data) execute() ([]byte, error) {
	s.Name = base.StringFirstUpper(s.Name)
	s.NameLower = base.StringFirstLower(s.Name)

	buf := new(bytes.Buffer)

	tmpl, err := template.New("data").Parse(dataTemplate)
	if err != nil {
		return nil, err
	}

	if err := tmpl.Execute(buf, s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s *Data) parseInject(databaseHandle []string) {
	for _, v := range databaseHandle {
		if v == NotInject {
			s.InjectMysql = false
			s.InjectMongo = false
			s.InjectRedis = false
			s.InjectHttp = false
			return
		} else if v == InjectMysql {
			s.InjectMysql = true
		} else if v == InjectMongo {
			s.InjectMongo = true
		} else if v == InjectRedis {
			s.InjectRedis = true
		} else if v == InjectHttp {
			s.InjectHttp = true
		}
	}
}
