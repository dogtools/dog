package main

import "fmt"

type userArgs struct {
	help      bool
	workdir   string
	directory string
	version   bool
	info      bool
	taskName  string
	taskArgs  map[string][]string
}

var knownFlags = [...]string{
	"-i", "--info",
	"-w", "--workdir",
	"-h", "--help",
	"-v", "--version",
	"-d", "--directory",
}

func printVersion() {
	fmt.Println("Dog version: " + version)
}

func printHelp() {
	fmt.Println(`Usage: dog
       dog [OPTIONS] TASK [ARGS]
       dog [--help] [--version]

Dog is a command line application that executes tasks.
Options:
  -i, --info       Print execution info (duration, exit status) after task execution
  -w, --workdir    Specify the working directory
  -d, --directory  Specify the dogfiles' directory
  -h, --help       Print usage information and help
  -v, --version    Print version information`)
}

func printNoValidDogfile() {
	fmt.Println(`Error: No valid Dogfile in current directory
Need help? --> dog --help
More info  --> https://github.com/dogtools/dog`)
}

func parseArgs(args []string) (a userArgs, err error) {

	// default values
	a = userArgs{
		help:      false,
		workdir:   "",
		directory: "",
		version:   false,
		info:      false,
		taskName:  "",
		taskArgs:  map[string][]string{},
	}

	skipArgument := false

	// iterate over all provided arguments
	for i, arg := range args {

		if skipArgument {
			skipArgument = false
			continue
		}

		if arg == "--help" || arg == "-h" {
			if i == 0 && len(args) == 1 && a.taskName == "" {
				a.help = true
				return a, nil
			}
			return a, fmt.Errorf("Error: %s doesn't accept additional parameters", arg)
		}

		if arg == "--version" || arg == "-v" {
			if i == 0 && len(args) == 1 && a.taskName == "" {
				a.version = true
				return a, nil
			}
			return a, fmt.Errorf("Error: %s doesn't accept additional parameters", arg)
		}

		if arg == "--info" || arg == "-i" {
			if a.taskName == "" {
				a.info = true
			} else {
				return a, fmt.Errorf("Error: %s is not a valid task argument", arg)
			}
		}

		if arg == "--workdir" || arg == "-w" {
			next := i + 1
			a.workdir = args[next]
			skipArgument = true
		}

		if arg == "--directory" || arg == "-d" {
			next := i + 1
			a.directory = args[next]
			skipArgument = true
		}

		if a.taskName == "" && string(arg[0]) != "-" {
			a.taskName = arg
		} else if a.taskName != "" && string(arg[0]) == "-" {
			if _, ok := a.taskArgs[arg]; !ok {
				a.taskArgs[arg] = []string{}
			}
		} else if a.taskName != "" && string(arg[0]) != "-" {
			if _, ok := a.taskArgs[args[i-1]]; !ok {
				return a, fmt.Errorf("Error: only one task can be executed at a time")
			}
			a.taskArgs[args[i-1]] = append(a.taskArgs[args[i-1]], arg)
		} else {
			validArg := false
			for _, f := range knownFlags {
				if arg == f {
					validArg = true
				}
			}
			if !validArg {
				return a, fmt.Errorf("Error: %s is not a valid argument", arg)
			}
		}
	}

	return a, nil
}
