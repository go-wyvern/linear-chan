package cmds

import (
	"time"
	"fmt"

	"github.com/go-wyvern/linear-chan/logger"
	"github.com/go-wyvern/linear-chan/model"
	"github.com/go-wyvern/linear-chan/dbs"
	"github.com/go-wyvern/linear-chan/utils"
)

var cmdSend = &Command{
	UsageLine: "send [filenames] [servername]",
	Short:     "send files to server",
	Long: `
将文件传到指定服务器[servername]
`,
	NeedMysql:true,
}

func init() {
	cmdSend.Run = Send
}

func Send(cmd *Command, args []string) int {
	if len(args) <= 1 {
		logger.DeployLog.Error("[filenames]或者[servername]未设置")
		cmd.Help()
		return 2
	}
	servername := args[len(args) - 1]
	files := args[:len(args) - 1]
	var project = new(models.Project)
	err := dbs.Db.Model(models.Project{}).Select("server,username,password,save_path").Where("name = ?", servername).Find(&project).Error
	if err != nil {
		logger.DeployLog.Error(err.Error())
		return 2
	}
	var tar_cmd string
	tar_name := fmt.Sprintf("%d.tar.gz", time.Now().Unix())
	tar_cmd = fmt.Sprintf("tar -zcf %s ", tar_name)
	for _, file := range files {
		tar_cmd = tar_cmd + file + " "
	}
	logger.DeployLog.Info("开始打包")
	err = utils.Shell(tar_cmd)
	if err != nil {
		logger.DeployLog.Error("打包失败：%s", err.Error())
		return 2
	}else {
		logger.DeployLog.Info("打包成功")
	}
	logger.DeployLog.Info("检测远程目录是否创建")
	err = utils.Shell("ssh %s@%s \"[ ! -d %s ] && mkdir %s\"", project.Username, project.Server, project.SavePath, project.SavePath)
	if err != nil {
		logger.DeployLog.Warning("检测目录出错，或创建远程目录失败：%s", err.Error())
	}else {
		logger.DeployLog.Info("检测成功")
	}
	logger.DeployLog.Info("开始发送TAR包到服务器:%s", project.Server)
	logger.DeployLog.Info("CMD : scp %s %s@%s:%s", tar_name, project.Username, project.Server, project.SavePath)
	err = utils.Shell("scp %s %s@%s:%s", tar_name, project.Username, project.Server, project.SavePath)
	if err != nil {
		logger.DeployLog.Error("发送TAR包失败：%s", err.Error())
		return 2
	}else {
		logger.DeployLog.Info("发送成功")
		utils.Shell("rm %s", tar_name)
	}
	return 0
}
