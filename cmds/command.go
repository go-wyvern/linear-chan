package cmds

import (
	"strings"
	"fmt"
	"os"
	"flag"
	"html/template"

	"github.com/go-wyvern/linear-chan/tmpl"
)


type Command struct {
	Run         func(cmd *Command, args []string) int
	Flag        flag.FlagSet
	CustomFlags bool
	NeedMysql   bool
	UsageLine   string
	Short       template.HTML
	Long        template.HTML
}

var Deploy struct {
	Version    string
	Configfile string
	LogerLevel int
}

func (c *Command) Name() string {
	name := c.UsageLine
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

func (c *Command) Usage() {
	fmt.Fprintf(os.Stderr, "usage: %s\n\n", c.UsageLine)
	fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(string(c.Long)))
	os.Exit(2)
}

func (c *Command) Help() {
	tmpl.Tmpl(os.Stdout, tmpl.HelpTemplate, c)
	os.Exit(2)
}

var Commands = []*Command{
	cmdVersion,
	cmdInit,
	cmdUp,
	cmdAutoAuth,
	cmdSSH,
	cmdSend,
	cmdCreate,
	cmdDelete,
	cmdUpdate,
	cmdLs,
}
