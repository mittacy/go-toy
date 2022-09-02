package model

import (
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
)

// CmdModel the service command.
var CmdModel = &cobra.Command{
	Use:   "model",
	Short: "Generate the model template implementations",
	Long:  "Generate the model template implementations. Example: gotoy tpl model xxx -t=app",
	Run:   run,
}

var targetDir string

func init() {
	CmdModel.Flags().StringVarP(&targetDir, "target-dir", "t", "app", "generate target directory")
}

func run(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		fmt.Fprintln(os.Stderr, "Please specify the model file. Example: gotoy tpl model xxx")
		return
	}

	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		fmt.Printf("Target directory: %s does not exist, example: gotoy tpl model xxx -t=app\n", targetDir)
		return
	}

	AddModel(args[0], targetDir)

	fmt.Println("success!")
}

func AddModel(name, dir string) bool {
	to := fmt.Sprintf("%s/model/%s.go", dir, name)

	if _, err := os.Stat(to); !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "%s model already exists: %s\n", name, to)
		return false
	}

	model := Model{
		Name: name,
	}
	b, err := model.execute()
	if err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile(to, b, 0644); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("create file %s\n", to)
	return true
}
