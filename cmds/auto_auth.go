package cmds

import (
	"os"

	"github.com/go-wyvern/linear-chan/logger"
	"github.com/go-wyvern/linear-chan/model"
	"github.com/go-wyvern/linear-chan/dbs"
	"github.com/go-wyvern/linear-chan/utils"
)

var cmdAutoAuth = &Command{
	UsageLine: "auth [servername]",
	Short:     "Send public key to server",
	Long: `
1.ssh-keygen -t rsa
2.scp id_rsa.pub server
3.cat id_rsa.pub>>authorized_keys
4.chmod 600 authorized_keys
`,
	NeedMysql:true,
}

func init() {
	cmdAutoAuth.Run = autoAuth
}

func autoAuth(cmd *Command, args []string) int {
	if len(args) == 0 {
		logger.DeployLog.Error("[servername]未设置")
		cmd.Help()
		return 2
	}
	var project = new(models.Project)
	err := dbs.Db.Model(models.Project{}).Select("server,username,password").Where("name = ?", args[0]).Find(&project).Error
	if err != nil {
		logger.DeployLog.Error(err.Error())
		return 2
	}
	logger.DeployLog.Info("检测并创建公钥")
	utils.Shell("[ ! -f $HOME/.ssh/id_rsa.pub ] && ssh-keygen -t rsa")
	logger.DeployLog.Alert("检测成功")
	logger.DeployLog.Info("复制公钥到服务器")
	var sshpath string
	if project.Username == "root" {
		sshpath = "/root/.ssh/"
	}else {
		sshpath = "/home/" + project.Username + "/.ssh/"
	}
	utils.Scp(os.Getenv("HOME") + "/.ssh/id_rsa.pub", sshpath, &utils.ServerConfig{
		ServerName:project.Server,
		UserName:project.Username,
		Password:project.Password,
	})
	logger.DeployLog.Alert("复制成功")
	logger.DeployLog.Info("添加公钥到authorized_keys")
	utils.SshAndDo(project.Server, project.Username, project.Password, "cat " + sshpath + "id_rsa.pub >>" + sshpath + "authorized_keys && chmod 600 " + sshpath + "authorized_keys")

	logger.DeployLog.Alert("认证成功")
	return 0
}