package cmds

import (
	"fmt"
	"os/exec"
	"html/template"
)

var cmdVersion = &Command{
	UsageLine: "version",
	Short:     "show the deploy version",
	Long: `
show the Deploy version

deploy version
    deploy   :`+template.HTML(Deploy.Version),
	NeedMysql:false,
}

func init() {
	cmdVersion.Run = versionCmd
}

func versionCmd(cmd *Command, args []string) int {
	fmt.Println("Deploy :deploy version " + Deploy.Version)
	goversion, err := exec.Command("go", "version").Output()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Go     :" + string(goversion))
	return 0
}