package listDir

import (
	"github.com/knagadevara/gbot/utl"
	"golang.org/x/crypto/ssh"
)

func Executor(host *ssh.Client) {
	utl.FireCommands(host, "ls")
}
