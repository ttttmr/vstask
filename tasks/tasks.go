package tasks

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/tidwall/jsonc"
)

type Tasks struct {
	filename string
	Version  string `json:"version"`
	Tasks    []Task `json:"tasks"`
}

func (t *Tasks) Print() string {
	var s []string
	s = append(s, fmt.Sprintf("File: %s", t.filename))
	for i, t := range t.Tasks {
		s = append(s, fmt.Sprintf("[%d] %s", i, t.Print()))
	}
	return strings.Join(s, "\n")
}

func ParseFromFile(p string) (*Tasks, error) {
	data, err := os.ReadFile(p)
	if err != nil {
		return nil, err
	}

	var tasks Tasks
	err = json.Unmarshal(jsonc.ToJSON(data), &tasks)
	if err != nil {
		return nil, err
	}
	tasks.filename = p

	return &tasks, nil
}

func (ts *Tasks) RunTask(name string) error {
	i, err := strconv.Atoi(name)
	if err == nil && i <= len(ts.Tasks) {
		return ts.Tasks[i].Run()
	}
	for _, t := range ts.Tasks {
		if t.Label == name {
			return t.Run()
		}
	}
	return fmt.Errorf("task %s not found", name)
}
