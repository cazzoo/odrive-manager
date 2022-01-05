package godrive

import (
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

type odriveClientHandler struct {
	cmd  *exec.Cmd
	path string
}

//go:generate stringer -type=OdriveCommand
type OdriveCommand int

const (
	Authenticate OdriveCommand = iota
	Mount
	Unmount
	Backup
	Backupnow
	Removebackup
	Sync
	Placeholderthreshold
	Foldersyncrule
	Unsync
	Autounsyncthreshold
	Stream
	Refresh
	Xlthreshold
	Encpassphrase
	Syncstate
	Status
	Deauthorize
	Diagnostics
	Emptytrash
	Autotrashthreshold
	Restoretrash
	Shutdown
)

type IOdriveClientHandler interface {
	Call(command OdriveCommand) []byte
}

func OdriveClientHandler(path string) IOdriveClientHandler {
	client := &odriveClientHandler{}
	client.path = path

	return client
}

func (client *odriveClientHandler) Call(command OdriveCommand) []byte {
	client.cmd = exec.Command(client.path, strings.ToLower(command.String()))
	client.cmd.Env = os.Environ()
	output, err := client.cmd.CombinedOutput()

	if err != nil {
		log.WithError(err).Error("Error with the client process")
	}

	return output
}
