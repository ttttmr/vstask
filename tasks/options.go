package tasks

type CommandOptions struct {
	Command string   `json:"command"`
	Args    []string `json:"args"`
	Options Options  `json:"options"`
}

type Options struct {
	Cwd   string            `json:"cwd"`
	Env   map[string]string `json:"env"`
	Shell Shell             `json:"shell"`
}

type Shell struct {
	Args       []string `json:"args"`
	Executable string   `json:"executable"`
}

func (co *CommandOptions) MergeFrom(other CommandOptions) {
	if len(other.Args) > 0 {
		co.Args = other.Args
	}
	if other.Command != "" {
		co.Command = other.Command
	}
	if other.Options.Cwd != "" {
		co.Options.Cwd = other.Options.Cwd
	}
	if other.Options.Env != nil {
		for k, v := range other.Options.Env {
			co.Options.Env[k] = v
		}
	}
	if other.Options.Shell.Executable != "" {
		co.Options.Shell.Executable = other.Options.Shell.Executable
	}
	if len(other.Options.Shell.Args) > 0 {
		co.Options.Shell.Args = other.Options.Shell.Args
	}
}
