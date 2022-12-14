package service

import (
	"fmt"
	"github.com/mittacy/go-toy/tools/gotoy/internal/base"
	"github.com/mittacy/go-toy/tools/gotoy/internal/tpl/data"
	"github.com/mittacy/go-toy/tools/gotoy/internal/tpl/model"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
)

// CmdService the service command.
var CmdService = &cobra.Command{
	Use:   "service",
	Short: "Generate the service template implementations",
	Long:  "Generate the service template implementations. Example: gotoy tpl service xxx -t=app",
	Run:   run,
}

var (
	targetDir      string
	databaseHandle []string
)

func init() {
	CmdService.Flags().StringVarP(&targetDir, "target-dir", "t", "app", "generate target directory")
	CmdService.Flags().StringArrayVarP(&databaseHandle, "database", "d", []string{data.InjectMysql, data.InjectRedis}, "inject database handle:null,mysql,redis,mongo,http, example: gotoy tpl data xxx -d mysql -d redis -d mongo -d http")
}

func run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Please specify the service file. Example: gotoy tpl service xxx")
		return
	}

	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		fmt.Printf("Target directory: %s does not exist, example: gotoy tpl service xxx -t=app\n", targetDir)
		return
	}

	modName, err := base.ModulePath("go.mod")
	if modName == "" || err != nil {
		fmt.Printf("go.mod no exist.\nPlease make sure you operate in the go project root directory\n")
		return
	}

	AddService(modName, args[0], targetDir)
	AddSModel(modName, args[0], targetDir)

	data.AddData(modName, args[0], targetDir, databaseHandle)

	model.AddModel(args[0], targetDir)

	fmt.Println("success!")
}

func AddService(appName, name, dir string) bool {
	to := fmt.Sprintf("%s/service/%s.go", dir, name)
	service := Service{
		AppName:   appName,
		Name:      name,
		TargetDir: dir,
	}

	if _, err := os.Stat(to); !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "%s service already exists: %s\n", name, to)
		return false
	}

	b, err := service.execute()
	if err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile(to, b, 0644); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("create file %s\n", to)
	return true
}

func AddSModel(appName, name, dir string) bool {
	to := fmt.Sprintf("%s/service/smodel/%s.go", dir, name)
	smodel := Model{
		AppName:   appName,
		Name:      name,
		TargetDir: dir,
	}

	if _, err := os.Stat(to); !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "%s service model already exists: %s\n", name, to)
		return false
	}

	b, err := smodel.execute()
	if err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile(to, b, 0644); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("create file %s\n", to)
	return true
}
