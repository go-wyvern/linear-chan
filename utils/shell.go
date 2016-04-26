package utils

import (
	"os"
	"fmt"
	"os/exec"

	"golang.org/x/crypto/ssh"
	"github.com/tmc/scp"
	"github.com/go-wyvern/linear-chan/logger"
)

func Shell(format string, a ...interface{}) error{
	var err error
	shellCmd := fmt.Sprintf(format, a...)
	cmd:=exec.Command("/bin/sh", "-c", shellCmd)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err =cmd.Start()
	if err!=nil{
		return err
	}
	err =cmd.Wait()
	if err!=nil{
		return err
	}
	return nil
}

type ServerConfig struct {
	ServerName string
	UserName string
	Password string
}

func Scp(filePath, destinationPath string,config *ServerConfig)error{
	client, err := ssh.Dial("tcp", config.ServerName+":22", &ssh.ClientConfig{
		User: config.UserName,
		Auth: []ssh.AuthMethod{ssh.Password(config.Password)},
	})
	if err != nil {
		return err
	}

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	session.Stdout = os.Stdout
	err = scp.CopyPath(filePath, destinationPath, session)
	if err != nil {
		return err
	}
	return nil
}

func SshAndDo(server,username,password,cmd string){
	client, err := ssh.Dial("tcp", server+":22", &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{ssh.Password(password)},
	})
	if err!=nil{
		logger.DeployLog.Error(err.Error())
		return
	}

	session, err := client.NewSession()
	if err!=nil{
		logger.DeployLog.Error(err.Error())
		return
	}
	defer session.Close()
	session.Stdout = os.Stdout
	if err := session.Run(cmd); err != nil {
		logger.DeployLog.Error(err.Error())
		panic("Failed to run: " + err.Error())
	}
}