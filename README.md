# omcli
Oh My! CLI is a small package for creating command line apps in go

This package only sports one level of sub commands, and flags may only be set on sub commands.

Commands are looked up by their Name, ie 'cmd run' looks for run a command with the name field being "run".

A help command is automatically added and runing 'cmd [subcommand] help' will print the help for [subcommand].

## Example

```go
package main

import (
	"fmt"
	"strconv"
	"strings"

	"bitbucket.org/gdl_iam/omcli"
)

var root = &omcli.Command{
	Name:  "test_cli",
	Usage: "test_cli",
	Short: "test_cli tests the cli package",
	Long:  "blah blah blah",

	Run: func(cmd *omcli.Command, args []string) {
		fmt.Println("hello world")
	},
}

var concat = &omcli.Command{
	Name:  "concat",
	Usage: "concat [flags]... [string]...",
	Short: "concat concatenates its args and prints them to stdout",
	Long:  `A much longer concat description`,

	Run: func(cmd *omcli.Command, args []string) {
		fmt.Println(strings.Join(args, ""))
	},
}

var add = &omcli.Command{
	Name:  "add",
	Usage: "add [flags]... [int]...",
	Short: "add gives the sum of the args passed",
	Long:  "Add does some hard core addition work",

	Run: func(cmd *omcli.Command, args []string) {
		var total int

		for _, arg := range args {
			val, err := strconv.Atoi(arg)
			if err != nil {
				panic(err)
			}
			total += val
		}
		fmt.Println(total)
	},
}

func init() {
	concat.Flags.String("blah", "test", "a string flag")
	root.AddCommand(concat)
	root.AddCommand(add)
}

func main() {
	root.Execute()
}


```
