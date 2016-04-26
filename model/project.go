package models

import (
	"reflect"

	"github.com/jinzhu/gorm"
	"github.com/go-wyvern/linear-chan/dbs"
)

type Project struct {
	gorm.Model
	Name        string `sql:"size:25" ini:"name"`    //项目名称

	TempVer     string `sql:"-"`

	ProjectName string `sql:"size:50"  ini:"project_name"`
	GitAddr     string `sql:"size:100" ini:"git_addr"`    //git地址
	GitUser     string `sql:"size:20" ini:"git_user"`   //git用户名
	GitPassword string `sql:"size:20" ini:"git_password"`   //git 密码

	UpdateCmd   string `sql:"size:100" ini:"update_cmd"`   //更新命令
	BuildPath   string `sql:"size:100" ini:"build_path"`   //编译路径
	UpdateFiles string `sql:"size:100" ini:"update_files"`    //需要上传到服务器的文件 ","分隔
	AnsibleYml  string `sql:"size:10000" ini:"ansible_yml"`  //待定

	Server      string `sql:"size:100" ini:"server"`   //服务器名
	Username    string `sql:"size:25" ini:"username"`   //服务器用户名
	Password    string `sql:"size:25" ini:"password"`    //服务器的密码
	SavePath    string `sql:"size:100" ini:"save_path"`  //部署到服务器的目录
}

//创建一个新的Project
func (c *Project) Create() error {
	return dbs.Db.Model(Project{}).Create(c).Error
}

//查询一个Project
func SelectProject(params map[string]interface{}, p *Pagination, attribute *Attribute) (*Project, error) {
	var err error
	var c = new(Project)
	if attribute != nil {
		Models.SetAttribute(Project{}, attribute)
	}
	err = Models.Select(Project{}, params, p).First(&c).Error

	return c, err
}

//根据属性查询多个Project
func SelectProjects(params map[string]interface{}, p *Pagination, attribute *Attribute) ([]Project, error) {
	var err error
	var c []Project
	if attribute != nil {
		Models.SetAttribute(Project{}, attribute)
	}
	err = Models.Select(Project{}, params, p).Find(&c).Error

	return c, err
}

//修改Project
func (c *Project) Update() error {
	return dbs.Db.Model(Project{}).Where("id = ?", c.ID).Update(c).Error
}

func (c *Project) UpdateColumn(name string) error {
	var err error
	v := reflect.ValueOf(c)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	err = dbs.Db.Model(Project{}).Where("id = ?", c.ID).UpdateColumn(name, v.FieldByName(name).Interface()).Error

	return err
}

func (c *Project) UpdateColumns(names ...string) error {
	var kv = make(map[string]interface{})
	v := reflect.ValueOf(c)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for _, name := range names {
		kv[name] = v.FieldByName(name).Interface()
	}
	return dbs.Db.Model(c).UpdateColumns(kv).Error
}

//删除Project
func (c *Project) Delete() error {
	return dbs.Db.Model(Project{}).Delete(c).Error
}

func (c *Project) CreateTable() error {
	return dbs.Db.CreateTable(&Project{}).Error
}

