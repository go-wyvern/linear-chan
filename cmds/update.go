package cmds

import (
	"github.com/go-wyvern/linear-chan/logger"
	"github.com/go-wyvern/linear-chan/config"
	"github.com/go-wyvern/linear-chan/model"
	"github.com/go-wyvern/linear-chan/dbs"
)

var cmdUpdate = &Command{
	UsageLine: "update [projectname]",
	Short:     "read /etc/deploy.conf and update tag(project) to mysql",
	Long: `
	从/etc/deploy.conf读取tag名字为项目名称，并修改数据库中指定项目的参数
`,
	NeedMysql:true,
}

func init() {
	cmdUpdate.Run = update
}

func update(cmd *Command, args []string) int {
	if len(args) == 0 {
		logger.DeployLog.Error("[projectname]未设置")
		cmd.Help()
		return 2
	}
	section, err := config.ProjectConfig.GetSection(args[0])
	if err != nil {
		logger.DeployLog.Error("读取配置出错，请检查你的配置文件里是否添加了你的项目:", err.Error())
		return 2
	}
	p := new(models.Project)
	err = dbs.Db.Model(models.Project{}).Where("name = ?", args[0]).First(&p).Error
	if err != nil {
		logger.DeployLog.Error("查找数据库中的项目出错:", err.Error())
	}
	err = section.MapTo(p)
	if err != nil {
		logger.DeployLog.Error("MAPTO ERROR:", err.Error())
		return 2
	}
	err = dbs.Db.Model(&p).UpdateColumns(p).Error
	if err != nil {
		logger.DeployLog.Error("更新失败：", err.Error())
	}
	logger.DeployLog.Info("更新成功")
	return 0
}


