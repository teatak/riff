package cli

import (
	"bytes"
	"fmt"
	"sort"
	"strings"
	"sync"
)

type Commands map[string]Command

type HelpFunc func(Commands) string

type Command interface {
	Run(args []string) int
	Help() string
	Synopsis() string
}

type CLI struct {
	Args           []string
	Commands       Commands
	Name           string
	Version        string
	HelpFunc       HelpFunc
	isHelp         bool
	once           sync.Once
	subCommand     string
	subCommandArgs []string
}

func NewCLI(name, version string) *CLI {
	return &CLI{
		Name:     name,
		Version:  version,
		HelpFunc: BasicHelpFunc(name),
	}
}
func (c *CLI) Run() (int, error) {
	c.once.Do(c.init)
	//get command
	command, ok := c.Commands[c.SubCommand()]
	if !ok {
		fmt.Print(c.HelpFunc(c.Commands))
		return 0, nil
	}
	if c.IsHelp() {
		fmt.Println(command.Help())
		return 0, nil
	}
	exitCode := command.Run(c.SubCommandArgs())
	return exitCode, nil
}

func (c *CLI) init() {
	if c.HelpFunc == nil {
		c.HelpFunc = BasicHelpFunc("app")

		if c.Name != "" {
			c.HelpFunc = BasicHelpFunc(c.Name)
		}
	}
	if len(c.Args) != 0 {
		arg := c.Args[0]
		if c.subCommand == "" {
			c.subCommand = arg
			c.subCommandArgs = c.Args[1:]
		}
	}
	for _, arg := range c.Args {
		if arg == "-h" || arg == "-help" || arg == "--help" {
			c.isHelp = true
			continue
		}
	}
}

func (c *CLI) IsHelp() bool {
	c.once.Do(c.init)
	return c.isHelp
}

func (c *CLI) SubCommand() string {
	c.once.Do(c.init)
	return c.subCommand
}

func (c *CLI) SubCommandArgs() []string {
	c.once.Do(c.init)
	return c.subCommandArgs
}

func BasicHelpFunc(name string) HelpFunc {
	return func(commands Commands) string {
		var buf bytes.Buffer
		buf.WriteString(fmt.Sprintf(
			"Usage: %s [--version] <command> [<args>]\n\n",
			name))
		buf.WriteString("Available commands are:\n\n")

		// Get the list of keys so we can sort them, and also get the maximum
		// key length so they can be aligned properly.
		keys := make([]string, 0, len(commands))
		maxKeyLen := 0
		for key, _ := range commands {
			if len(key) > maxKeyLen {
				maxKeyLen = len(key)
			}

			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			command, ok := commands[key]
			if !ok {
				// This should never happen since we JUST built the list of
				// keys.
				panic("command not found: " + key)
			}

			key = fmt.Sprintf("%s%s", key, strings.Repeat(" ", maxKeyLen-len(key)))
			buf.WriteString(fmt.Sprintf("  %-12s%s\n", key, command.Synopsis()))
		}
		return buf.String()
	}
}
