package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/structs"
	"github.com/sh3rp/eyes/agent"
	"github.com/sh3rp/eyes/controller"
	"github.com/sh3rp/eyes/db"
	"github.com/sh3rp/eyes/util"
	"github.com/spf13/cobra"
	"gopkg.in/abiosoft/ishell.v1"
)

var consoleMode bool

var handler = ResultHandler{}

var RootCmd = &cobra.Command{
	Use:   "eyesctl",
	Short: "",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		if consoleMode {
			wrapper := NewWrapper(ishell.New(), controller.NewController("/tmp", handler.handleResult))
			wrapper.Start()
		}
	},
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&consoleMode, "console", "c", false, "Run in console mode (for debugging)")
	RootCmd.AddCommand(AgentCmd, ServerCmd)
}

type ResultHandler struct {
	results []agent.Result
}

func (handler ResultHandler) handleResult(r agent.Result) {
	handler.results = append(handler.results, r)
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
	wrapper.shell.Register("newsched", wrapper.newSchedCmd)
	wrapper.shell.Register("lssched", wrapper.lsSchedCmd)
	wrapper.shell.Register("lsagent", wrapper.lsAgentCmd)
	wrapper.shell.Register("newagent", wrapper.newAgentCmd)
	wrapper.shell.Register("deploy", nil)
	wrapper.shell.Register("undeploy", nil)
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
	str += "Configurations\n"
	str += fmt.Sprintf("ID                         Type           Parameters\n")
	str += fmt.Sprintf("========================== ============== ===============================\n")
	for _, c := range configs {
		str += fmt.Sprintf("%25s %-14s %+v\n", c.Id, structs.Name(agent.ACTIONS[c.Action]), c.Parameters)
	}
	return str, err
}

func (wrapper ShellWrapper) newSchedCmd(args ...string) (string, error) {
	configId := args[0]
	cronString := args[1]

	schedule := db.Schedule{
		ConfigId: util.ID(configId),
		Schedule: cronString,
	}

	newSchedule, err := wrapper.controller.NewSchedule(schedule)
	return fmt.Sprintf("%+v\n", newSchedule), err
}

func (wrapper ShellWrapper) lsSchedCmd(args ...string) (string, error) {
	schedules, err := wrapper.controller.GetSchedules()

	str := "\n"
	str += "Schedules\n"
	str += fmt.Sprintf("ID                         ConfigId                   Cron\n")
	str += fmt.Sprintf("========================== ========================== ===============================\n")
	for _, s := range schedules {
		str += fmt.Sprintf("%25s %25s %s\n", s.Id, s.ConfigId, s.Schedule)
	}
	return str, err
}

func (wrapper ShellWrapper) newAgentCmd(args ...string) (string, error) {
	agent, err := wrapper.controller.NewAgentLocal()
	return fmt.Sprintf("%+v\n", agent), err
}

func (wrapper ShellWrapper) lsAgentCmd(args ...string) (string, error) {
	agents, err := wrapper.controller.GetAgents()

	str := "\n"
	str += "Agents\n"
	str += fmt.Sprintf("ID                         Type\n")
	str += fmt.Sprintf("========================== ====\n")
	for _, a := range agents {
		str += fmt.Sprintf("%25s %d\n", a.Id, a.AgentType)
	}
	return str, err
}
