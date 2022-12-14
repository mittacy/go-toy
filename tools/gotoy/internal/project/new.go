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
		fmt.Printf("š« %s already exists\n", p.Name)
		override := false
		prompt := &survey.Confirm{
			Message: "š Do you want to override the folder ?",
			Help:    "Delete the existing folder and create the project.",
		}
		survey.AskOne(prompt, &override)
		if !override {
			return err
		}
		os.RemoveAll(to)
	}

	fmt.Printf("š Creating service %s, layout repo is %s, please wait a moment.\n\n", p.Name, layout)

	repo := base2.NewRepo(layout, branch)

	if err := repo.CopyTo(ctx, to, p.Path, []string{".git", ".github"}); err != nil {
		return err
	}

	os.Rename(
		path.Join(to, "cmd", "server"),
		path.Join(to, "cmd", p.Name),
	)
	base2.Tree(to, dir)

	fmt.Printf("\nšŗ Project creation succeeded %s\n", color.GreenString(p.Name))

	fmt.Printf("Wait a moment, the program is in the final configuration work\n")

	// ęæę¢é”¹ē®äø­ēå­ē¬¦äø²
	base2.Replace(to, replaceStr, p.Name)

	// å¤å¶éē½®ęä»¶
	developEnv := fmt.Sprintf("%s/.env.development", to)
	localEnv := fmt.Sprintf("%s/.env", to)
	_ = base2.Copy(localEnv, developEnv)

	fmt.Print("š» Use the following command to start the project š:\n\n")

	fmt.Println(color.WhiteString("$ cd %s", p.Name))
	fmt.Println(color.WhiteString("$ go mod download "))
	fmt.Println(color.WhiteString("edit the .env.* configuration file"))
	fmt.Println(color.WhiteString("$ go run . start http -c=.env.development -e=development -p=8080\n"))
	fmt.Println("		š¤ Thanks for using go-toy")
	fmt.Println("	š Tutorial: https://mittacy.github.io/blog/column/1624512335520")
	return nil
}
