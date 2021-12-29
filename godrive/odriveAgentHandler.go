package godrive

import (
	"os/exec"

	"cazzoo.me/godrive/process"
	log "github.com/sirupsen/logrus"
)

type IOdriveAgentHandler interface {
	Start() error
	Stop() error
	HealthCheck() bool
	KillProcess() error
}

type odriveAgentHandler struct {
	cmd  *exec.Cmd
	path string
}

func OdriveAgentHandler(path string) IOdriveAgentHandler {
	agent := &odriveAgentHandler{}
	agent.path = path

	return agent
}

func (agent *odriveAgentHandler) Start() error {
	agent.cmd = exec.Command(agent.path)
	err := agent.cmd.Start()
	if err != nil {
		log.WithError(err).Error("Error with the agent process")
	}

	return err
}

func (agent *odriveAgentHandler) Stop() error {
	var err error
	if processes, err := process.FindProcess("odriveagent"); err == nil {
		if err := process.KillProcesses(processes); err != nil {
			log.Warning("Error stoping agent.")
		}
	} else {
		return err
	}
	return err
}

func (agent *odriveAgentHandler) HealthCheck() bool {
	var pid = 0
	if agent.cmd != nil && agent.cmd.Process != nil {
		pid = agent.cmd.Process.Pid
	}
	return pid != 0
}

func (agent *odriveAgentHandler) KillProcess() error {
	err := agent.cmd.Process.Kill()
	if err != nil {
		log.WithError(err).Fatal("Error killing agent process")
	}

	return err
}
