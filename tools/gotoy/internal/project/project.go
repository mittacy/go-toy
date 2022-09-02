package project

import (
	"context"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/spf13/cobra"
	"os"
	"path"
	"time"
)

// CmdNew represents the new command.
var CmdNew = &cobra.Command{
	Use:   "new",
	Short: "Create a gin project",
	Long:  "Create a service project using the repository template. Example: gotoy new helloworld",
	Run:   run,
}

var (
	repoURL = "git@github.com:mittacy/go-toy-layout.git"
	branch  = "master"
)

func init() {
	CmdNew.Flags().StringVarP(&repoURL, "repo-url", "r", repoURL, "layout repo")
	CmdNew.Flags().StringVarP(&branch, "branch", "b", branch, "repo branch")
}

func run(cmd *cobra.Command, args []string) {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	// 获取项目名
	name := ""
	if len(args) == 0 {
		prompt := &survey.Input{
			Message: "What is project name ?",
			Help:    "Created project name.",
		}
		survey.AskOne(prompt, &name)
		if name == "" {
			return
		}
	} else {
		name = args[0]
	}

	p := &Project{Name: path.Base(name), Path: name}
	if err := p.New(ctx, wd, repoURL, branch); err != nil {
		fmt.Fprintf(os.Stderr, "\033[31mERROR: %s\033[m\n", err)
		return
	}
}
