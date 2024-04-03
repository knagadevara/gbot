package main

import (
	"log"

	utl "github.com/knagadevara/gbot/utl"
)

var (
	dcdn        *utl.DCDN
	sshConfig   *utl.SShCfg
	hostAuth    *utl.HostAuth
	bastionAuth *utl.BastionAuth
	hj          utl.HJSShConfig
)

func init() {
	yamlBuf := utl.LoadFile("files/config.yaml")
	sshConfig, hostAuth, bastionAuth, dcdn = utl.ParseCfg(yamlBuf)
	hj.SSHConfig = sshConfig
	hj.HostAuth = hostAuth
	hj.BastionAuth = bastionAuth
	hj.Dx = dcdn
	hstName := utl.SourceHostName()
	err := hj.MapHostDc(hstName)
	if err != nil {
		log.Fatalf("%v", err)
	}
	hj.DisplayHostDetails()

}

func main() {
	if hj.SSHConfig.Jump {
		rmtHstSshClt, jumpSshClient, err := hj.CreateSshClientJumpHost()
		if err != nil {
			log.Fatalln(err)
		}
		defer rmtHstSshClt.Close()
		defer jumpSshClient.Close()
		utl.FireCommands(rmtHstSshClt, "ls", "pwd", "who")

	} else {
		hstSshClt, err := hj.CreateSshClientHost()
		if err != nil {
			log.Fatalln(err)
		}
		defer hstSshClt.Close()
		utl.FireCommands(hstSshClt, "ls", "pwd", "who")

	}
}
