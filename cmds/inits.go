package cmds

import (
	"os"

	"github.com/go-wyvern/linear-chan/logger"
	"github.com/go-wyvern/linear-chan/utils"
)

var cmdInit = &Command{
	UsageLine: "init",
	Short:     "add config to /etc/linear-chan.conf",
	Long: `
Add config to /etc/linear-chan.ini and /etc/linear-chan.d/*.ini
`,
	NeedMysql:false,
}

var configTemplate = `[mysql]
username= root
password =
debug = false
database = deploy
address = ""
parameters = charset=utf8&parseTime=True&loc=Local
maxidle = 10
maxopen = 100
`

var pconfigTemplate = `[example]
project_name=
git_addr=
git_user=
git_password=
update_cmd=
build_path=
update_files=
server=
username=
password=
save_path=
ansible_yml=[{"task_name":"测试服务链接","do_type":"ping","do_cmd":""},{"task_name":"创建最新版本目录","do_type":"shell","do_cmd":"ls -l"}]
`

func init() {
	cmdInit.Run = Init
}

func Init(cmd *Command, args []string) int {
	conf, err := os.Create("/etc/linear-chan.ini")
	if err != nil {
		logger.DeployLog.Error("Create config file error:", err.Error())
		return 2
	}
	defer conf.Close()
	_, err = conf.Write([]byte(configTemplate))
	if err != nil {
		logger.DeployLog.Error("写入错误：", err.Error())
		return 2
	}
	logger.DeployLog.Info("创建deploy.ini成功！")
	utils.Shell("[ ! -d /etc/deploy.d ] && mkdir /etc/deploy.d")
	projectconf, err := os.Create("/etc/deploy.d/projects.ini")
	if err != nil {
		logger.DeployLog.Error("Create deploy.d config file error:", err.Error())
		return 2
	}
	defer projectconf.Close()
	_, err = projectconf.Write([]byte(pconfigTemplate))
	if err != nil {
		logger.DeployLog.Error("写入错误：", err.Error())
		return 2
	}
	logger.DeployLog.Info("创建/etc/deploy.d/projects.ini成功！")
	return 0
}
