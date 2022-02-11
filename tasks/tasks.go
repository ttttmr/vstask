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

func (ts *Tasks) FindTaskByName(name string) (*Task, error) {
	for _, t := range ts.Tasks {
		if t.Label == name {
			return &t, nil
		}
	}
	return nil, fmt.Errorf("task %s not found", name)
}

func (ts *Tasks) FindTaskByIndex(i int) (*Task, error) {
	if i < 0 || i >= len(ts.Tasks) {
		return nil, fmt.Errorf("task %d not found", i)
	}
	return &ts.Tasks[i], nil
}

func (ts *Tasks) RunTask(name string) error {
	// find task
	task, err := ts.FindTaskByName(name)
	if task == nil {
		// try index
		i, ierr := strconv.Atoi(name)
		if ierr != nil {
			return err
		}
		task, err = ts.FindTaskByIndex(i)
		if err != nil {
			return err
		}
	}
	fmt.Printf("> Executing task: %s <\n\n", task.Print())
	// depends
	if task.DependsOrder == "" || task.DependsOrder == "parallel" {
		rt := make(chan error, len(task.DependsOn))
		// run depends
		for _, depName := range task.DependsOn {
			go func(depName string, rt chan error) {
				err := ts.RunTask(depName)
				if err == nil {
					rt <- err
				} else {
					rt <- fmt.Errorf("%s %s", depName, err)
				}
			}(depName, rt)
		}
		// check rt
		var finErr error
		for range task.DependsOn {
			e := <-rt
			if e != nil {
				fmt.Printf("\n> Task depend error: %s <\n", e)
				finErr = e
			}
		}
		if finErr != nil {
			return finErr
		}
	} else {
		for _, depName := range task.DependsOn {
			err = ts.RunTask(depName)
			if err != nil {
				return err
			}
		}
	}

	return task.Run()
}
