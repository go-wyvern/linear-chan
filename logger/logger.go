package logger

import "github.com/astaxie/beego/logs"

var DeployLog *logs.BeeLogger

func init(){
	DeployLog= logs.NewLogger(10000)
	DeployLog.SetLogger("console", `{"level":8}`)
	DeployLog.EnableFuncCallDepth(true)
}
