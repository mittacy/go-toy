package method

import (
	"fmt"
	"github.com/mittacy/go-toy/tools/gotoy/internal/base"
	"github.com/spf13/cobra"
	"log"
	"os"
	"strings"
)

// CmdMethod the service command.
var CmdMethod = &cobra.Command{
	Use:   "method",
	Short: "Generate the method template implementations",
	Long:  "Generate the method template implementations. Example: gotoy tpl method xxx(服务名) xxx(方法名) xxx(注释) -t=app",
	Run:   run,
}
var (
	targetDir string
)

func init() {
	CmdMethod.Flags().StringVarP(&targetDir, "target-dir", "t", "app", "generate target directory")
}

func run(cmd *cobra.Command, args []string) {
	if len(args) < 3 {
		fmt.Fprintln(os.Stderr, "Please specify the method template. Example: gotoy tpl method xxx(服务名) xxx(方法名) xxx(注释) -t=app")
		return
	}

	serviceName := args[0]
	methodName := args[1]
	commentName := args[2]

	AddMethod(targetDir, serviceName, methodName, commentName)

	fmt.Fprintln(os.Stderr, "finish")
}

func AddMethod(dir, serviceName, methodName, commentName string) bool {
	method := Method{
		Name:      base.StringFirstUpper(serviceName),
		NameLower: base.StringFirstLower(serviceName),
		Method:    base.StringFirstUpper(methodName),
		Comment:   commentName,
	}

	allTpl := []struct {
		Tpl      string
		FilePath string
	}{
		{
			Tpl:      strings.Replace(validatorTpl, "${backquote}", "`", -1),
			FilePath: dir + "/validator/" + serviceName + "Vdr/" + serviceName + ".go",
		},
		{
			Tpl:      sModelTpl,
			FilePath: dir + "/service/smodel/" + serviceName + ".go",
		},
		{
			Tpl:      dpTpl,
			FilePath: dir + "/dp/" + serviceName + ".go",
		},
		{
			Tpl:      serviceTpl,
			FilePath: dir + "/service/" + serviceName + ".go",
		},
		{
			Tpl:      apiTpl,
			FilePath: dir + "/api/" + serviceName + ".go",
		},
	}

	for _, v := range allTpl {
		// 解析模板
		b, err := method.execute(v.Tpl)
		if err != nil {
			log.Fatal(err)
		}

		// 写入文件
		file, err := os.OpenFile(v.FilePath, os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Fprintf(os.Stdout, "open file %s, err: %s\n", v.FilePath, err)
			continue
		}
		defer file.Close()

		_, err = file.Write(b)
		if err != nil {
			fmt.Printf("writer method to file %s fail, err: %s\n", v.FilePath, err)
			continue
		}
		fmt.Printf("writer method to file %s success\n", v.FilePath)
	}
	return true
}
