package spiderhn

import (
	"testing"
	"time"
)

func TestRunCmds(t *testing.T) {
	cmds := []Cmd{
		Cmd{"git", []string{"stash"}},
		Cmd{"git", []string{"pull", "origin", "master"}},
		Cmd{"git", []string{"stash", "apply"}},
		Cmd{"git", []string{"add", "."}},
		Cmd{"git", []string{"stash"}},
		Cmd{"git", []string{"commit", "-am", time.Now().Format("2006-01-02T15:04:05")}},
		Cmd{"ps", []string{"ps", "-ef", "|", "grep", "md-genie"}},
		Cmd{"netstat", []string{"-lntp"}},
		Cmd{"free", []string{"-m"}},
		Cmd{"ps", []string{"aux"}},
	}
	if logs, err := RunCmds(cmds); err != nil {
		t.Fatal(err)
	} else {
		t.Log(logs)
	}
}
