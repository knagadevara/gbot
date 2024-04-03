package main

import (
	utl "github.com/knagadevara/gbot/utl"
)

var (
	hj utl.HJSShConfig
)

func init() {
	yamlBuf := utl.LoadFile("files/config.yaml")
	hj = utl.ParseCfg(yamlBuf)
}

func main() {
	hj.Execute("ls", "pwd", "who")
}
