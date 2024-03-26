package main

import (
	"fmt"
	"os"

	"github.com/jackgpalfrey/peurl-cli/cmds"
	"github.com/jackgpalfrey/peurl-cli/utils"
)

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		cmds.Help()
		utils.UsageError("peurl <subcommand>")
	}
	subCommand := args[0]
	subArgs := args[1:]

	switch subCommand {
	case "help":
		cmds.Help()
	case "domain":
		cmds.Domain(subArgs)
	case "login":
		cmds.Login(subArgs)
	case "logout":
		cmds.Logout(subArgs)
	case "whoami":
		cmds.Whoami()
	case "user":
		cmds.User(subArgs)
	case "shorten", "s", "new":
		cmds.UrlShorten(subArgs)
	case "go":
		cmds.UrlGo(subArgs)
	case "inspect":
		cmds.UrlInspect(subArgs)
	case "expand":
		cmds.UrlExpand(subArgs)
	case "delete", "rm":
		cmds.UrlDelete(subArgs)
	case "list", "ls":
		cmds.UrlList(subArgs)
	case "update":
		cmds.Update()
	case "uninstall":
		cmds.Uninstall()
	default:
		fmt.Printf("Invalid subcommand '%s'\n", subCommand)
	}
}
