# vstask

Run vscode task in command line
## install

```
go get "github.com/ttttmr/vstask"
```

## usage

```
NAME:
   vstask - Run vscode task in command line

USAGE:
   c [global options] command [command options] [arguments...]

COMMANDS:
   ls       List all task
   run      Run a task by name or index
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h  show help (default: false)
```

## features

- task type
  - process
  - shell
- command options
  - cwd
  - env
  - shell
- os specific options
  - windows
  - linux
  - mac
- predefined variables
  - `${cwd}`
  - `${defaultBuildTask}`
  - `${pathSeparator}`
  - `${workspaceFolderBasename}`
  - `${workspaceFolder}`

## todo

- more task type
  - npm
  - ...
- task depend
  - dependsOn
  - dependsOrder

## references

https://code.visualstudio.com/docs/editor/tasks
https://code.visualstudio.com/docs/editor/variables-reference#_predefined-variables