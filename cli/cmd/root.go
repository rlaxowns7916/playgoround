package cmd

import (
	"os"
	"playgoround/cli/cmd/kill"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "myps",
	Short: "My Custom Process Management Tool",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(kill.Cmd())
}
