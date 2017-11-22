package cmd

import (
	"github.com/spf13/cobra"
)

var consoleMode bool

var RootCmd = &cobra.Command{
	Use:   "eyesctl",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&consoleMode, "console", "c", false, "Run in console mode (for debugging)")
	RootCmd.AddCommand(AgentCmd, ServerCmd)
}
