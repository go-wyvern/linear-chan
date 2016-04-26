package tmpl

import (
	"io"
	"html/template"
	"strings"
)

var UsageTemplate = `Deploy is a tool for managing company Automated Deployment.

Usage:

    deploy command [arguments]

The commands are:
{{range .}}
    {{.UsageLine | printf "%-30s"}} {{.Short}}{{end}}

If need config file pls use -f config file path in the End

Use "deploy help [command]" for more information about a command.

Additional help topics:
{{range .}}
    {{.UsageLine | printf "%-30s"}} {{.Short}}{{end}}

Use "deploy help [topic]" for more information about that topic.

`
var HelpTemplate = `usage: deploy {{.UsageLine}}

{{.Long | trim}}
`


func Tmpl(w io.Writer, text string, data interface{})error {
	t := template.New("top")
	t.Funcs(template.FuncMap{"trim": func(s template.HTML) template.HTML {
		return template.HTML(strings.TrimSpace(string(s)))
	}})
	template.Must(t.Parse(text))
	if err := t.Execute(w, data); err != nil {
		return err
	}

	return nil
}
