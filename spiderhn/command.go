package spiderhn

import (
	"fmt"
	"os/exec"
	"strings"
)

type Cmd struct {
	Name string
	Args []string
}

func RunCmds(cmds []Cmd) (logs []string, err error) {
	for _, command := range cmds {
		cmd := exec.Command(command.Name, command.Args...)
		out, err := cmd.Output()
		var logString string
		if err != nil {
			logString = fmt.Sprint("错误:", command)
		} else {
			logString = string(out)
		}
		if strings.TrimSpace(logString) != "" {
			logs = append(logs, logString)
		}
	}
	return logs, nil
}
