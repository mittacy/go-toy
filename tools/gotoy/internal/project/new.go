package project

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	base2 "github.com/mittacy/go-toy/tools/gotoy/internal/base"
	"os"
	"path"
)

// Project is a project template.
type Project struct {
	Name string
	Path string
}

const (
	replaceStr = "go-toy-layout"
)

// New new a project from remote repo.
func (p *Project) New(ctx context.Context, dir string, layout string, branch string) error {
	to := path.Join(dir, p.Name)

	if _, err := os.Stat(to); !os.IsNotExist(err) {
		fmt.Printf("🚫 %s already exists\n", p.Name)
		override := false
		prompt := &survey.Confirm{
			Message: "📂 Do you want to override the folder ?",
			Help:    "Delete the existing folder and create the project.",
		}
		survey.AskOne(prompt, &override)
		if !override {
			return err
		}
		os.RemoveAll(to)
	}

	fmt.Printf("🚀 Creating service %s, layout repo is %s, please wait a moment.\n\n", p.Name, layout)

	repo := base2.NewRepo(layout, branch)

	if err := repo.CopyTo(ctx, to, p.Path, []string{".git", ".github"}); err != nil {
		return err
	}

	os.Rename(
		path.Join(to, "cmd", "server"),
		path.Join(to, "cmd", p.Name),
	)
	base2.Tree(to, dir)

	fmt.Printf("\n🍺 Project creation succeeded %s\n", color.GreenString(p.Name))

	fmt.Printf("Wait a moment, the program is in the final configuration work\n")

	// 替换项目中的字符串
	base2.Replace(to, replaceStr, p.Name)

	// 复制配置文件
	developEnv := fmt.Sprintf("%s/.env.development", to)
	localEnv := fmt.Sprintf("%s/.env", to)
	_ = base2.Copy(localEnv, developEnv)

	fmt.Print("💻 Use the following command to start the project 👇:\n\n")

	fmt.Println(color.WhiteString("$ cd %s", p.Name))
	fmt.Println(color.WhiteString("$ go mod download "))
	fmt.Println(color.WhiteString("edit the .env.* configuration file"))
	fmt.Println(color.WhiteString("$ go build -o ./server ."))
	fmt.Println(color.WhiteString("$ ./server start http -c=.env.development -e=development -p=8080\n"))
	fmt.Println("		🤝 Thanks for using go-toy")
	//fmt.Println("	📚 Tutorial: http://www.mittacy.com/column/1633512445750")
	return nil
}
