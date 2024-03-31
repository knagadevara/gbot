package main

import (
	"log"

	utl "github.com/knagadevara/gbot/utl"
)

func main() {
	// sshConfig := sshConn.LoadSshConfig("sdsdsd", "dsds", "dsdsdsd.pem")
	// sshClient := sshConn.DialHost(sshConfig, "192.168.MM.6", 22)
	// if _, err := sshConn.FireCommand(sshClient, "ls", "whoami", "pwd", "uptime", "who"); err != nil {
	// 	fmt.Println(err)
	// }
	if hostDetails, err := utl.GetHostDcJmp("bs31-cjj1"); err != nil {
		log.Fatalf("%v", err)
	} else {
		utl.DisplayHostDetails(hostDetails)
	}
}
