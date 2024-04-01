package main

import (
	"fmt"
	"log"

	scon "github.com/knagadevara/gbot/connv"
	utl "github.com/knagadevara/gbot/utl"
)

func main() {
	ymlCfg := utl.ParseYaml("files/sshConfig.yaml")
	sshConfig := scon.LoadSshConfig(ymlCfg.SshBase, ymlCfg.UserName, ymlCfg.KeyName)
	hstName := utl.SourceHostName()
	if hostDetails, err := utl.GetHostDcJmp(hstName); err != nil {
		log.Fatalf("%v", err)
	} else {
		utl.DisplayHostDetails(hostDetails)
		sshClient := scon.DialHost(sshConfig, hostDetails.HostName, 22)
		if _, err := scon.FireCommand(sshClient, "ls", "whoami", "pwd", "uptime", "who"); err != nil {
			fmt.Println(err)
		}
	}
}
