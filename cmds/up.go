package cmds

import (
	"os"
	"time"
	"fmt"
	"bytes"
	"encoding/json"

	"github.com/go-wyvern/linear-chan/logger"
	"github.com/go-wyvern/linear-chan/model"
	"github.com/go-wyvern/linear-chan/dbs"
	"github.com/go-wyvern/linear-chan/utils"
	"github.com/go-wyvern/linear-chan/ansible"
	"github.com/go-wyvern/linear-chan/tmpl"
)

var cmdUp = &Command{
	UsageLine: "up [projectname]",
	Short:     "Automated Deployment server",
	Long: `
	1.git clone or git pull from git server
	2.go build project
	3.tar -zcvf project.tar.gz
	4.send project.tar.gz to Deployment server
	5.restart server
`,
	NeedMysql:true,
}

func init() {
	cmdUp.Run = up
}

func up(cmd *Command, args []string) int {
	if len(args) == 0 {
		logger.DeployLog.Error("[projectname]未设置")
		cmd.Help()
		return 2
	}
	p := new(models.Project)
	err := dbs.Db.Model(models.Project{}).Where("name = ?", args[0]).First(&p).Error
	if err != nil {
		logger.DeployLog.Error("查找数据库中的项目出错: %v", err.Error())
		return 2
	}
	version := p.Name + time.Now().Format("2006-01-02-15-04")
	GOPATH := os.Getenv("GOPATH")
	logger.DeployLog.Info("获取源码中...")
	if isDirExists(GOPATH + "/src/" + p.ProjectName) {
		logger.DeployLog.Info("源码目录已存在，更新源码")
		err = utils.Shell("cd %s/src/%s && %s", GOPATH, p.ProjectName, p.UpdateCmd)
	}else {
		logger.DeployLog.Info("源码目录不存在，下载源码")
		err = utils.Shell("cd %s/src && git clone %s", GOPATH, p.GitAddr)
	}
	if err != nil {
		logger.DeployLog.Error("获取失败: %v", err.Error())
		return 2
	}else {
		logger.DeployLog.Alert("获取源码成功")
	}
	logger.DeployLog.Info("开始编译源码")
	err = utils.Shell("cd %s && go build", p.BuildPath)
	if err != nil {
		logger.DeployLog.Error("编译失败: %v", err.Error())
		return 2
	}
	logger.DeployLog.Alert("编译成功")

	var tar_cmd string
	tar_name := fmt.Sprintf("%s.tar.gz", version)
	logger.DeployLog.Debug("tar name: %s", tar_name)
	tar_cmd = fmt.Sprintf("cd %s && tar -zcf %s %s", p.BuildPath, tar_name, p.UpdateFiles)
	logger.DeployLog.Debug("tar cmd: %s", tar_cmd)
	logger.DeployLog.Info("开始打包")
	err = utils.Shell(tar_cmd)
	if err != nil {
		logger.DeployLog.Error("打包失败：%s", err.Error())
		return 2
	}else {
		logger.DeployLog.Info("打包成功")
	}
	logger.DeployLog.Info("检测远程目录是否创建")
	err = utils.Shell("ssh %s@%s \"[ ! -d %s ] && mkdir %s\"", p.Username, p.Server, p.SavePath, p.SavePath)
	if err != nil {
		logger.DeployLog.Warning("检测目录出错，或创建远程目录失败：%s", err.Error())
	}else {
		logger.DeployLog.Info("检测成功")
	}
	logger.DeployLog.Info("开始发送TAR包到服务器:%s", p.Server)
	logger.DeployLog.Info("CMD : scp %s/%s %s@%s:%s", p.BuildPath, tar_name, p.Username, p.Server, p.SavePath)
	err = utils.Shell("scp %s/%s %s@%s:%s", p.BuildPath, tar_name, p.Username, p.Server, p.SavePath)
	if err != nil {
		logger.DeployLog.Error("发送TAR包失败：%s", err.Error())
		return 2
	}else {
		logger.DeployLog.Info("发送成功")
		utils.Shell("rm %s/%s", p.BuildPath, tar_name)
	}

	logger.DeployLog.Info("自动化配置ansible.yml")
	ansible := absible.NewAnsible()
	var tasks []absible.AnsibleTask
	var buf = make([]byte, 0)
	w := bytes.NewBuffer(buf)
	p.TempVer = version
	tmpl.Tmpl(w, p.AnsibleYml, p)
	p.AnsibleYml = w.String()
	logger.DeployLog.Debug(p.AnsibleYml)
	json.Unmarshal([]byte(p.AnsibleYml), &tasks)
	logger.DeployLog.Debug("tasks:%v", tasks)
	logger.DeployLog.Info("ansible配置生成中...")
	ansible = ansible.SetHost(p.Server).SetRemoteUser("root").SetVersion(version).SetRootDir(p.SavePath)

	for _, task := range tasks {
		logger.DeployLog.Debug(task.DoCmd)
		ansible = ansible.AddTask(task.TaskName, task.DoType, task.DoCmd)
	}
	err = ansible.GenAnsibleYml(p.Name)
	if err != nil {
		logger.DeployLog.Error("配置失败：%v", err.Error())
		return 2
	}
	logger.DeployLog.Alert("配置成功")

	logger.DeployLog.Info("ansible部署开始")
	err = utils.Shell("ansible-playbook %s.yml", p.Name)
	if err != nil {
		logger.DeployLog.Error("部署失败：%v", err.Error())
		return 2
	}
	logger.DeployLog.Alert("部署成功！")
	utils.Shell("rm %s.yml", p.Name)
	return 0
}

func isDirExists(path string) bool {
	fi, err := os.Stat(path)

	if err != nil {
		return os.IsExist(err)
	} else {
		return fi.IsDir()
	}

	panic("not reached")
}

