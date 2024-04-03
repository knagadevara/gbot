package main

import (
	"os"

	"github.com/knagadevara/gbot/listDir"
	utl "github.com/knagadevara/gbot/utl"
	"golang.org/x/crypto/ssh"
)

var (
	hj         utl.HJSShConfig
	host, jump *ssh.Client
	yamlBuf    []byte
)

func init() {
	yamlBuf = utl.LoadFile("files/config.yaml")
	hj = utl.ParseCfg(yamlBuf)
}

func main() {
	host, jump = hj.JumpOrNot()
	defer utl.CloseConn(host, jump)

	argument := os.Args[0]
	switch argument {
	case "ls":
		listDir.Executor(host)
	}
}
