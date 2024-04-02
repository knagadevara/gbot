package main

import (
	"fmt"
	"log"

	scon "github.com/knagadevara/gbot/connv"
	"github.com/knagadevara/gbot/stypes"
	utl "github.com/knagadevara/gbot/utl"
)

var ymlCfg stypes.SShCfg

func init() {
	ymlCfg = utl.ParseYaml("files/sshConfig.yaml")
}

func main() {
	hstName := utl.SourceHostName()
	if hostDetails, err := utl.GetHostDcJmp(hstName); err != nil {
		log.Fatalf("%v", err)
	} else {
		utl.DisplayHostDetails(hostDetails)
		if ymlCfg.JumpHost {
			sshJumpConfig := scon.LoadSshConfig(ymlCfg.SshBase, ymlCfg.JumpUser, ymlCfg.KeyName)
			sshRmtConfig := scon.LoadSshConfig(ymlCfg.SshBase, ymlCfg.UserName, ymlCfg.KeyName)
			if sshJumpClient, err := scon.DialHost(sshJumpConfig, hostDetails.JumpBox, 22); err != nil {
				log.Fatalln(err)
			} else {
				if hostConnection, jumpCLient, err := scon.ConnHost(sshJumpClient, hostDetails.HostName); err != nil {
					log.Fatalln(err)
				} else {
					if rmtHstSshClient, jmpBxSshClient, err := scon.MakeNewClientConn(hostConnection, jumpCLient, hostDetails.HostName, sshRmtConfig); err != nil {
						log.Fatalln(err)
					} else {
						if _, err := scon.FireCommands(rmtHstSshClient, jmpBxSshClient, "pwd"); err != nil {
							fmt.Println(err)
						}
					}
				}

			}
		}
	}
}
