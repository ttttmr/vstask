package vscode

import (
	"strings"

	"github.com/ttttmr/vstask/tasks"
)

type PreVar map[string]string

// https://code.visualstudio.com/docs/editor/variables-reference#_predefined-variables
func NewPreVar() PreVar {
	return map[string]string{
		"cwd":                     "",
		"defaultBuildTask":        "",
		"pathSeparator":           "",
		"workspaceFolderBasename": "",
		"workspaceFolder":         "",
	}
}

func (p *PreVar) ReplaceString(s string) string {
	for k, v := range *p {
		s = strings.Replace(s, "${"+k+"}", v, -1)
	}
	return s
}

func (p *PreVar) ReplaceTasks(ts *tasks.Tasks) {
	for i := 0; i < len(ts.Tasks); i++ {
		p.ReplaceTask(&ts.Tasks[i])
	}
}

func (p *PreVar) ReplaceTask(t *tasks.Task) {
	p.ReplaceCommandOptions(&t.CommandOptions)
	p.ReplaceCommandOptions(&t.Linux)
	p.ReplaceCommandOptions(&t.Osx)
	p.ReplaceCommandOptions(&t.Windows)
}

func (p *PreVar) ReplaceCommandOptions(co *tasks.CommandOptions) {
	// command
	co.Command = p.ReplaceString(co.Command)
	// args
	for j := 0; j < len(co.Args); j++ {
		co.Args[j] = p.ReplaceString(co.Args[j])
	}
	// options
	co.Options.Cwd = p.ReplaceString(co.Options.Cwd)
	for k, v := range co.Options.Env {
		co.Options.Env[k] = p.ReplaceString(v)
	}
	co.Options.Shell.Executable = p.ReplaceString(co.Options.Shell.Executable)
	for j := 0; j < len(co.Options.Shell.Args); j++ {
		co.Options.Shell.Args[j] = p.ReplaceString(co.Options.Shell.Args[j])
	}
}
