package cmd

import (
	"strconv"
	"strings"

	"github.com/sh3rp/eyes/controller"
	"github.com/sh3rp/eyes/db"
	"github.com/spf13/cobra"
	"gopkg.in/abiosoft/ishell.v1"
)

var consoleMode bool

var RootCmd = &cobra.Command{
	Use:   "eyesctl",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		if consoleMode {
			shell := ishell.New()

			// display welcome info.
			shell.Println("EyeShell v1.0")

			// run shell
			shell.Start()
		}
	},
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&consoleMode, "console", "c", false, "Run in console mode (for debugging)")
	RootCmd.AddCommand(AgentCmd, ServerCmd)
}

type ShellWrapper struct {
	Controller controller.ControllerServer
	Shell      *ishell.Shell
}

func initShell() *ishell.Shell {
	shell := ishell.New()
	configCmd := ishell.CmdFunc(func(args ...string) (string, error) {
		cfgType, _ := strconv.Atoi(args[0])
		parameters := make(map[string]string)

		for i := 1; i < len(args); i++ {
			tokens := strings.Split(args[i], "=")
			parameters[tokens[0]] = tokens[1]
		}

		cfg := db.Config{
			Action:     cfgType,
			Parameters: parameters,
		}
		return "", nil
	})

	shell.Register("newcfg", configCmd)
	return shell
}
