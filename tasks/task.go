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
	var s []string

	if t.Detail == "" {
		s = append(s, t.Label)
	} else {
		s = append(s, fmt.Sprintf("%s(%s)", t.Label, t.Detail))
	}

	if cmd, err := t.BuildCmd(); err == nil {
		s = append(s, fmt.Sprintf("%v", cmd.Args))
	}

	return strings.Join(s, " ")
}

func (t *Task) Run() error {
	cmd, err := t.BuildCmd()
	if err != nil {
		return err
	}
	return cmd.Run()
}

func (t *Task) BuildCmd() (*exec.Cmd, error) {
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
		return nil, fmt.Errorf("unsupported os: %s", runtime.GOOS)
	}
	t.Command = co.Command
	t.Args = co.Args
	t.Options = co.Options
	// build command
	switch t.Type {
	case "shell":
		if t.Options.Shell.Executable != "" {
			cmd = exec.Command(t.Options.Shell.Executable, t.Options.Shell.Args...)
			cmd.Args = append(cmd.Args, t.BuildCommand()...)
		} else {
			// use user default shell
			shell, err := exec.LookPath(os.Getenv("SHELL"))
			if err != nil {
				return nil, err
			}

			cmd = exec.Command(shell, "-c", strings.Join(t.BuildCommand(), " "))
		}
	case "process", "":
		cmd = exec.Command(t.Command, t.Args...)
	default:
		return nil, fmt.Errorf("unsupported task type: %v", t.Type)
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
	return cmd, nil
}

func (t *Task) BuildCommand() []string {
	if len(t.Args) == 0 {
		return []string{t.Command}
	} else {
		return []string{
			fmt.Sprintf("'%s'", t.Command),
			strings.Join(t.Args, " "),
		}
	}
}
