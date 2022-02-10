package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ttttmr/vstask/vscode"
	"github.com/urfave/cli/v2" // imports as package "cli"
)

func main() {
	tasks, err := vscode.GetTasks()
	if err != nil {
		fmt.Println(err)
		return
	}
	app := &cli.App{
		Name:  "vstask",
		Usage: "Run vscode task in command line",
		Commands: []*cli.Command{
			{
				Name:  "ls",
				Usage: "List all task",
				Action: func(c *cli.Context) error {
					fmt.Println(tasks.Print())
					return nil
				},
			},
			{
				Name:  "run",
				Usage: "Run a task by name or index",
				Action: func(c *cli.Context) error {
					return tasks.RunTask(c.Args().First())
				},
			},
		}}
	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
