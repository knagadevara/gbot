package main

import (
	"runtime"

	"github.com/knagadevara/gbot/commands"
	"github.com/knagadevara/gbot/utl"
	"golang.org/x/crypto/ssh"
)

var (
	host, jump *ssh.Client
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	yamlBuf := utl.LoadFile("files/config.yaml")
	hj := utl.ParseCfg(yamlBuf)
	host, jump = hj.JumpOrNot()
}

func main() {
	commands.GeneralSystemStats(host)
	defer utl.CloseConn(host, jump)
}
