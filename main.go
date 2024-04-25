package main

import (
	"runtime"

	"github.com/knagadevara/gbot/commands"
	"github.com/knagadevara/gbot/utl"
	"golang.org/x/crypto/ssh"
)

var (
	hj         *utl.HJSShConfig
	host, jump *ssh.Client
	yamlBuf    []byte
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	yamlBuf = utl.LoadFile("files/config.yaml")
	hj = utl.ParseCfg(yamlBuf)
	host, jump = hj.JumpOrNot()
}

func main() {
	defer utl.CloseConn(host, jump)
	commands.GeneralSystemStats(host)
}
