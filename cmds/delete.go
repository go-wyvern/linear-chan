package cmds

import (
	"github.com/go-wyvern/linear-chan/logger"
	"github.com/go-wyvern/linear-chan/dbs"
	"github.com/go-wyvern/linear-chan/model"
)

var cmdDelete = &Command{
	UsageLine: "delete [projectname]",
	Short:     "delete tag(project) from mysql",
	Long: `
	删除数据库中指定项目
`,
	NeedMysql:true,
}

func init() {
	cmdDelete.Run = delete
}

func delete(cmd *Command, args []string) int {
	if len(args) == 0 {
		logger.DeployLog.Error("[projectname]未设置")
		cmd.Help()
		return 2
	}
	p := new(models.Project)
	err := dbs.Db.Model(models.Project{}).Where("name = ?", args[0]).First(&p).Error
	if err != nil {
		logger.DeployLog.Error("查找数据库中的项目出错:", err.Error())
	}
	err = dbs.Db.Unscoped().Delete(&p).Error
	if err != nil {
		logger.DeployLog.Error("删除失败：", err.Error())
	}
	logger.DeployLog.Info("删除成功")
	return 0
}