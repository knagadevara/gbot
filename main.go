package main

import (
	"fmt"

	sshConn "utl/SshConn"
)

func main() {
	sshConfig := sshConn.LoadSshConfig("sdsdsd", "dsds", "dsdsdsd.pem")
	sshClient := sshConn.DialHost(sshConfig, "192.168.MM.6", 22)
	if _, err := sshConn.FireCommand(sshClient, "ls", "whoami", "pwd", "uptime", "who"); err != nil {
		fmt.Println(err)
	}
}
