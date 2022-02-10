package vscode

import (
	"os"
	"path/filepath"

	"github.com/ttttmr/vstask/tasks"
	"github.com/ttttmr/vstask/utils"
)

func GetTasks() (*tasks.Tasks, error) {
	p, err := FindTasksFile()
	if err != nil {
		return nil, err
	}
	ts, err := tasks.ParseFromFile(p)
	if err != nil {
		return nil, err
	}
	// predefined variables
	prevar := NewPreVar()
	prevar["cwd"], err = os.Getwd()
	if err != nil {
		return nil, err
	}
	if len(ts.Tasks) > 0 {
		prevar["defaultBuildTask"] = ts.Tasks[0].Label
	}
	prevar["pathSeparator"] = string(os.PathSeparator)
	prevar["workspaceFolderBasename"] = filepath.Base(p)
	prevar["workspaceFolder"] = filepath.Dir(filepath.Dir(p))
	// replace task predefined variables
	prevar.ReplaceTasks(ts)

	return ts, nil
}

func FindTasksFile() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		p := filepath.Join(dir, ".vscode", "tasks.json")
		if utils.FileExists(p) {
			return p, nil
		}

		next := filepath.Dir(dir)
		if next == dir {
			return "", os.ErrNotExist
		}
		dir = next
	}
}
