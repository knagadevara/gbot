package commands

import (
	"github.com/knagadevara/gbot/utl"
	"golang.org/x/crypto/ssh"
)

func GeneralSystemStats(host *ssh.Client) {
	utl.RunCommand(host, "uname -n", "uptime", "free -h")
}
