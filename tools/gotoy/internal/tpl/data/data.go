package data

import (
	"fmt"
	"github.com/mittacy/go-toy/tools/gotoy/internal/base"
	"github.com/mittacy/go-toy/tools/gotoy/internal/tpl/model"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
)

var (
	targetDir      string
	databaseHandle []string
)

// CmdData the data command.
var CmdData = &cobra.Command{
	Use:   "data",
	Short: "Generate the data template implementations",
	Long:  "Generate the data template implementations. Example: ego tpl data xxx -t=internal",
	Run:   run,
}

const (
	NotInject   = "null"
	InjectMysql = "mysql"
	InjectMongo = "mongo"
	InjectRedis = "redis"
	InjectHttp  = "http"
)

func init() {
	CmdData.Flags().StringVarP(&targetDir, "target-dir", "t", "app", "generate target directory")
	CmdData.Flags().StringArrayVarP(&databaseHandle, "database", "d", []string{InjectMysql, InjectRedis}, "inject database handle:null,mysql,redis,mongo,http, example: gotoy tpl data xxx -d mysql -d redis -d mongo -d http")
}

func run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Please specify the data file. Example: gotoy tpl data xxx")
		return
	}

	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		fmt.Printf("Target directory: %s does not exist, example: gotoy tpl data xxx -t=app\n", targetDir)
		return
	}

	modName, err := base.ModulePath("go.mod")
	if modName == "" || err != nil {
		fmt.Printf("go.mod no exist.\nPlease make sure you operate in the go project root directory\n")
		return
	}

	AddData(modName, args[0], targetDir, databaseHandle)

	model.AddModel(args[0], targetDir)

	fmt.Println("success!")
}

func AddData(appName, name, dir string, databaseHandle []string) bool {
	to := fmt.Sprintf("%s/data/%s.go", dir, name)
	if _, err := os.Stat(to); !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "%s data already exists: %s\n", name, to)
		return false
	}

	data := Data{
		AppName:   appName,
		Name:      name,
		TargetDir: dir,
	}
	data.parseInject(databaseHandle)
	b, err := data.execute()
	if err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile(to, b, 0644); err != nil {
		log.Fatalf("写入data文件错误, err: %+v", err)
	}

	fmt.Printf("create file %s\n", to)
	return true
}
