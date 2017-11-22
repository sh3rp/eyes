package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/structs"
	"github.com/sh3rp/eyes/agent"
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
			wrapper := NewWrapper(ishell.New(), controller.NewController("/tmp"))
			wrapper.Start()
		}
	},
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&consoleMode, "console", "c", false, "Run in console mode (for debugging)")
	RootCmd.AddCommand(AgentCmd, ServerCmd)
}

type ShellWrapper struct {
	controller controller.Controller
	shell      *ishell.Shell
}

func NewWrapper(shell *ishell.Shell, ctrl controller.Controller) ShellWrapper {
	wrapper := ShellWrapper{
		controller: ctrl,
		shell:      shell,
	}
	wrapper.shell.Register("newcfg", wrapper.newCfgCmd)
	wrapper.shell.Register("lscfg", wrapper.lsCfgCmd)
	return wrapper
}

func (wrapper ShellWrapper) Start() {
	wrapper.shell.Println("EyeShell v1.0")
	wrapper.shell.Start()
}

func (wrapper ShellWrapper) newCfgCmd(args ...string) (string, error) {
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
	newCfg, err := wrapper.controller.NewConfig(cfg)
	return fmt.Sprintf("%+v\n", newCfg), err
}

func (wrapper ShellWrapper) lsCfgCmd(args ...string) (string, error) {
	configs, err := wrapper.controller.GetConfigs()

	str := "\n"
	str += "Configurations"
	str += fmt.Sprintf("ID                         Type           Parameters\n")
	str += fmt.Sprintf("========================== ============== ===============================\n")
	for _, c := range configs {
		str += fmt.Sprintf("%25s %-14s %+v\n", c.Id, structs.Name(agent.ACTIONS[c.Action]), c.Parameters)
	}
	return str, err
}
