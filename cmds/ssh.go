package cmds

import (
	"github.com/go-wyvern/linear-chan/logger"
	"github.com/go-wyvern/linear-chan/dbs"
	"github.com/go-wyvern/linear-chan/model"
	"github.com/go-wyvern/linear-chan/utils"
)

var cmdSSH = &Command{
	UsageLine: "ssh [servername]",
	Short:     "SSH to server",
	Long: `
SSH 登录到服务器[servername]
`,
	NeedMysql:true,
}

func init() {
	cmdSSH.Run = SSH
}

func SSH(cmd *Command, args []string) int {
	if len(args) == 0 {
		logger.DeployLog.Error("[servername]未设置")
		cmd.Help()
		return 2
	}
	var project = new(models.Project)
	err := dbs.Db.Model(models.Project{}).Select("server,username,password").Where(" name = ?", args[0]).Find(&project).Error
	if err != nil {
		logger.DeployLog.Error(err.Error())
		return 2
	}
	utils.Shell("ssh %s@%s", project.Username, project.Server)
	return 0
}