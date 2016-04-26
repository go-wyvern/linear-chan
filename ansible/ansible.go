package absible

import (
	"os"

	"github.com/go-wyvern/linear-chan/logger"
	"github.com/go-wyvern/linear-chan/tmpl"
)

var ansibleTemplate = `---

  - hosts: {{.Host}}
    remote_user: {{.RemoteUser}}
    vars:
         version: {{.Version}}  #本次部署版本名
         rootdir: {{.RootDir}}
    tasks:{{range .Tasks}}
        - name: {{.TaskName}}
          {{.DoType}}: {{.DoCmd}}{{end}}
`
type Ansible struct {
	Host       string
	RemoteUser string
	Version    string
	RootDir    string
	Tasks      []AnsibleTask
}

type AnsibleTask struct {
	TaskName string   `json:"task_name"`
	DoType   string   `json:"do_type"`
	DoCmd    string   `json:"do_cmd"`
}

func NewAnsible() *Ansible {
	return new(Ansible)
}

func (a *Ansible) SetHost(host string) *Ansible {
	a.Host = host
	return a
}

func (a *Ansible) SetRemoteUser(remote_user string) *Ansible {
	a.RemoteUser = remote_user
	return a
}

func (a *Ansible) SetVersion(version string) *Ansible {
	a.Version = version
	return a
}

func (a *Ansible) SetRootDir(root_dir string) *Ansible {
	a.RootDir = root_dir
	return a
}

func (a *Ansible) AddTask(task_name, do_type, do_cmd string) *Ansible {
	ansible_task := AnsibleTask{
		TaskName:task_name,
		DoType:do_type,
		DoCmd:do_cmd,
	}
	a.Tasks = append(a.Tasks, ansible_task)
	return a
}

func (a *Ansible) GenAnsibleYml(name string) error {
	f, err := os.Create(name + ".yml")
	if err != nil {
		logger.DeployLog.Error("Create ansible config file error:", err.Error())
		return err
	}
	defer f.Close()
	err = tmpl.Tmpl(f, ansibleTemplate, a)
	if err != nil {
		logger.DeployLog.Error("Parse ansible config file error:", err.Error())
		return err
	}
	return nil
}


