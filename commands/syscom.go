package commands

import (
	"github.com/knagadevara/gbot/utl"
	"golang.org/x/crypto/ssh"
)

func GeneralSystemStats(host *ssh.Client) {
	utl.FireCommands(host, "uname -n", "uptime", "free -h")
}
