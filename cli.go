package omcli

import (
	"flag"
	"fmt"
	"os"
)

// Command is a command to be run from the command line
type Command struct {

	// Name is the name of the command, if empty,
	// fallback to trying first word of usage, then short, then long
	Name string

	// Usage is the single line usage information
	Usage string

	// Short is a short description of the command used when help is called
	Short string

	// Long is the long description of the command used when help [command] is called
	Long string

	// Flags are the command line flags for the command
	// The root command will not have it's flags parsed, flags are reserved
	// for subcommands
	Flags flag.FlagSet

	// Run runs the command
	Run func(*Command, []string)

	SubCommands []*Command

	root *Command
}

// AddCommand ands a subcommand to the command
func (c *Command) AddCommand(command *Command) {
	if c.root != nil {
		panic("commands can only be added to a root")
	}
	command.root = c
	c.SubCommands = append(c.SubCommands, command)
}

// Execute runs the command
func (c *Command) Execute() {

	if c.root != nil {
		panic("Execute should only be called on a root command")
	}

	flag.Usage = c.help
	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		c.help()
		os.Exit(-1)
	}

	cname := args[0]

	if cname == "help" {
		if len(args) < 2 {
			c.help()
			os.Exit(-1)
		}
		c.doHelp(args[1])
		os.Exit(0)
	}

	cmd := c.findCommand(cname)
	if cmd == nil {
		c.unknownExit(cname)
	}

	cmd.Flags.Usage = cmd.usage
	if err := cmd.Flags.Parse(args[1:]); err != nil {
		panic(err)
	}
	args = cmd.Flags.Args()
	cmd.Run(cmd, args)
}

func (c *Command) unknownExit(name string) {
	c.println("unknown command: ", name)
	c.println("Run", "'"+c.Name, "help' for available commands.")
	os.Exit(-1)
}

func (c *Command) doHelp(name string) {
	cmd := c.findCommand(name)
	if cmd == nil {
		c.unknownExit(name)
	}
	cmd.help()
}

func (c *Command) findCommand(name string) *Command {
	for _, cmd := range c.SubCommands {
		if cmd.Name == name && cmd.runnable() {
			return cmd
		}
	}
	return nil
}

func (c *Command) usage() {
	c.println("usage: ", c.Usage)
	c.println()
	c.Flags.SetOutput(os.Stderr)
	c.Flags.PrintDefaults()
	c.println()
	c.println(c.Long)
	os.Exit(-1)
}

func (c *Command) help() {
	c.printf("%s\n\n", c.Short)

	if c.root == nil {
		c.printf("Usage:\n     %s command [arguments]\n\n", c.Name)
		c.printf("Available commands:\n\n")
		for _, cmd := range c.SubCommands {
			if cmd.runnable() {
				c.printf("%+10s     %s\n", cmd.Name, cmd.Short)
			}
		}
		c.printf("\nUse '%s help' [command] to view a command's documentation.\n\n", c.Name)
		return
	}
	c.printf("Usage:\n\n  %s\n\n", c.Usage)
	c.Flags.SetOutput(os.Stderr)
	c.Flags.PrintDefaults()
	c.printf("\n%s\n\n", c.Long)
}

func (c *Command) printf(format string, args ...interface{}) {
	if _, err := fmt.Fprintf(os.Stderr, format, args...); err != nil {
		panic(err)
	}
}

func (c *Command) println(args ...interface{}) {
	if _, err := fmt.Fprintln(os.Stderr, args...); err != nil {
		panic(err)
	}
}

func (c *Command) runnable() bool {
	return c.Run != nil
}
