package main

import (
	"os"
	"fmt"
	"flag"

	"github.com/go-wyvern/linear-chan/config"
	"github.com/go-wyvern/linear-chan/cmds"
	"github.com/go-wyvern/linear-chan/dbs"
	"github.com/go-wyvern/linear-chan/tmpl"
)

const CONF_VER = "0.0.1"
const DefaultConfigFile = "/etc/deploy.conf"

func main() {
	flag.Usage = usage
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		usage()
	}

	if args[0] == "help" {
		help(args[1:])
		return
	}
	args = LoadConfig(args)
	for _, cmd := range cmds.Commands {
		if cmd.Name() == args[0] && cmd.Run != nil {
			if cmd.NeedMysql {
				err := config.InitConfig()
				if err != nil {
					panic(err)
				}
				dbs.InitMysql()
			}
			cmd.Flag.Usage = func() {
				cmd.Usage()
			}
			if cmd.CustomFlags {
				args = args[1:]
			} else {
				cmd.Flag.Parse(args[1:])
				args = cmd.Flag.Args()
			}
			os.Exit(cmd.Run(cmd, args))
			return
		}
	}

	fmt.Fprintf(os.Stderr, "deploy: unknown subcommand %q\nRun 'deploy help' for usage.\n", args[0])
	os.Exit(2)

}

func usage() {
	tmpl.Tmpl(os.Stdout, tmpl.UsageTemplate, cmds.Commands)
	os.Exit(2)
}

func help(args []string) {
	if len(args) == 0 {
		usage()
		return
	}
	if len(args) != 1 {
		fmt.Fprintf(os.Stdout, "usage: deploy help command\n\nToo many arguments given.\n")
		os.Exit(2)
	}

	arg := args[0]

	for _, cmd := range cmds.Commands {
		if cmd.Name() == arg {
			tmpl.Tmpl(os.Stdout, tmpl.HelpTemplate, cmd)
			return
		}
	}

	fmt.Fprintf(os.Stdout, "Unknown help topic %#q.  Run 'deploy help'.\n", arg)
	os.Exit(2)
}
//读取其他配置：/etc/deploy.conf
func LoadConfig(arguments []string) []string {
	var IsSet = false
	var newargs []string
	newargs = arguments
	cmds.Deploy.Version=CONF_VER
	for i, arg := range arguments {
		if arg == "-f" {
			cmds.Deploy.Configfile = arguments[i + 1]
			newargs = append(arguments[:i], arguments[i + 2:]...)
			IsSet = true
			break
		}
	}
	if !IsSet {
		cmds.Deploy.Configfile = DefaultConfigFile
	}
	return newargs
}


