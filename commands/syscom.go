package commands

import (
	"github.com/knagadevara/gbot/utl"
	"golang.org/x/crypto/ssh"
)

func GeneralSystemStats(host *ssh.Client) {
	sessionOutChan := utl.ExecuteCommand(host, "uname -n", "uptime", "free -h")
	utl.PrintHashMap(sessionOutChan)

}
