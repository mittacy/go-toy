package main

import (
	"github.com/mittacy/go-toy/tools/gotoy/internal/project"
	"github.com/mittacy/go-toy/tools/gotoy/internal/tpl"
	"github.com/spf13/cobra"
	"log"
)

const version = "v0.1.0"

var rootCmd = &cobra.Command{
	Use:     "gotoy",
	Short:   "gotoy: An elegant toolkit for Gin.",
	Long:    `gotoy: An elegant toolkit for Gin.`,
	Version: version,
}

func init() {
	rootCmd.AddCommand(project.CmdNew)
	rootCmd.AddCommand(tpl.CmdTpl)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
