package snaprd

import (
	"os"
	"os/exec"

	log "github.com/sirupsen/logrus"
)

// ExecCmd runs a snapraid command and returns the logfile location and an error
func (s *Snaprd) ExecCmd(command string, args ...string) (*os.File, error) {
	logFile, err := os.CreateTemp("/tmp", "snaprd")
	if err != nil {
		return nil, err
	}
	s.LogFiles = append(s.LogFiles, logFile.Name())
	cmdString := []string{
		command,
		"--conf",
		s.Config.Snapraid.Config,
		"-l",
		logFile.Name(),
	}
	cmdString = append(cmdString, args...)

	cmd := exec.Command(s.Config.Snapraid.Executable, cmdString...)
	output, err := cmd.CombinedOutput()
	log.WithFields(log.Fields{"command": command}).Info(string(output))
	if err != nil {
		// If diff detects differences it returns status code 2
		if command == "diff" {
			if exitError, ok := err.(*exec.ExitError); ok {
				if exitError.ExitCode() == 2 {
					return logFile, nil
				}
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	return logFile, nil
}
