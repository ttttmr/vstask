package tasks

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Task struct {
	Label  string `json:"label"`
	Detail string `json:"detail"`
	Group  string `json:"group"`

	Type         string   `json:"type"` // process/shell
	DependsOn    []string `json:"dependsOn"`
	DependsOrder string   `json:"dependsOrder"` // parallel/sequence

	CommandOptions
	// os
	Linux   CommandOptions `json:"linux"`
	Osx     CommandOptions `json:"osx"`
	Windows CommandOptions `json:"windows"`
}

func (t Task) Print() string {
	if t.Detail == "" {
		return t.Label
	} else {
		return fmt.Sprintf("%s(%s)", t.Label, t.Detail)
	}
}

func (t *Task) Run() error {
	var cmd *exec.Cmd
	// command options
	var co CommandOptions // tmp
	co.Command = t.Command
	co.Args = t.Args
	co.Options = t.Options
	switch runtime.GOOS {
	case "darwin":
		co.MergeFrom(t.Osx)
	case "linux":
		co.MergeFrom(t.Linux)
	case "windows":
		co.MergeFrom(t.Windows)
	default:
		return fmt.Errorf("unsupported os: %s", runtime.GOOS)
	}
	t.Command = co.Command
	t.Args = co.Args
	t.Options = co.Options
	// build command
	switch t.Type {
	case "shell":
		if t.Options.Shell.Executable != "" {
			cmd = exec.Command(t.Options.Shell.Executable, t.Options.Shell.Args...)
		} else {
			// use user default shell
			shell, err := exec.LookPath(os.Getenv("SHELL"))
			if err != nil {
				return err
			}
			args := []string{
				fmt.Sprintf("'%s'", t.Command),
				strings.Join(t.Args, " "),
			}

			cmd = exec.Command(shell, "-c", strings.Join(args, " "))
		}
	case "process":
		cmd = exec.Command(t.Command, t.Args...)
	default:
		return fmt.Errorf("unsupported task type: %v", t.Type)
	}
	// cwd
	cmd.Dir = t.Options.Cwd
	// env
	for k, v := range t.Options.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
