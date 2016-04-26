package cmds

import (
	"github.com/jinzhu/gorm"

	"github.com/go-wyvern/linear-chan/logger"
	"github.com/go-wyvern/linear-chan/config"
	"github.com/go-wyvern/linear-chan/model"
	"github.com/go-wyvern/linear-chan/dbs"
)

var cmdCreate = &Command{
	UsageLine: "create [projectname]",
	Short:     "read /etc/deploy.conf and add tag(project) to mysql",
	Long: `
	从/etc/deploy.conf读取tag名字为项目名称，并将配置写入数据库
`,
	NeedMysql:true,
}

func init() {
	cmdCreate.Run = create
}

func create(cmd *Command, args []string) int {
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
	p.CreateTable()
	err = dbs.Db.Model(models.Project{}).Where("name = ?", args[0]).First(&p).Error
	if err != nil&&err != gorm.ErrRecordNotFound {
		logger.DeployLog.Error(err.Error())
	}else if err == nil {
		logger.DeployLog.Error("你要创建的项目，mysql中已存在,请使用deploy ls查看！")
	}
	p.Server = section.Key("server").Value()
	p.Username = section.Key("username").Value()
	p.Password = section.Key("password").Value()
	p.SavePath = section.Key("save_path").Value()
	p.GitAddr = section.Key("git_addr").Value()
	p.GitPassword = section.Key("git_password").Value()
	p.GitUser = section.Key("git_user").Value()
	p.UpdateCmd = section.Key("update_cmd").Value()
	p.UpdateFiles = section.Key("update_files").Value()
	p.BuildPath = section.Key("build_path").Value()
	p.AnsibleYml = section.Key("ansible_yml").Value()
	p.Name = args[0]
	err = dbs.Db.Model(models.Project{}).Create(p).Error
	if err != nil {
		logger.DeployLog.Error("创建失败：%v", err)
	}
	logger.DeployLog.Info("创建成功")
	return 0
}

