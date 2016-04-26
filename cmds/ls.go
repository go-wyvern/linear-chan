package cmds

import (
	"os"

	"github.com/go-wyvern/linear-chan/model"
	"github.com/go-wyvern/linear-chan/logger"
	"github.com/go-wyvern/linear-chan/dbs"
	"github.com/go-wyvern/linear-chan/tmpl"
)

var cmdLs = &Command{
	UsageLine: "ls [servername]",
	Short:     "list server message",
	Long: `
	list servername message or all server message
`,
	NeedMysql:true,
}

var projectTemplate =
 ` +-----------------------+------------------------------------------------------+----------------------------------+
 |        项目名称       |                     Git服务器地址                    |           部署服务器地址         |
 +=======================+======================================================+==================================+
{{range .}} {{.Name| printf "|   %-20s|"}}{{.GitAddr | printf "   %-51s|"}}{{.Server| printf "   %-31s|"}}{{end}}
 +-----------------------+------------------------------------------------------+----------------------------------+
`

func init() {
	cmdLs.Run = ls
}

func ls(cmd *Command, args []string) int {
	var servername string
	if len(args) == 1 {
		servername = args[0]
	}
	var params = make(map[string]interface{})
	var projects []models.Project
	if servername != "" {
		params["name"] = servername
		err := dbs.Db.Model(models.Project{}).Where(params).Find(&projects).Error
		if err != nil {
			logger.DeployLog.Error(err.Error())
		}
	}else {
		err := dbs.Db.Model(models.Project{}).Find(&projects).Error
		if err != nil {
			logger.DeployLog.Error(err.Error())
		}
	}
	tmpl.Tmpl(os.Stdout, projectTemplate, projects)
	return 0

}
