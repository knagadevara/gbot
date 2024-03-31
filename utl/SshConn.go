package utl

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	utl . "Uitil"

	ssh "golang.org/x/crypto/ssh"
	kh "golang.org/x/crypto/ssh/knownhosts"
)

func LoadSshConfig(sshBasePath string, sshUserName string, sshPkName string) *ssh.ClientConfig {
	knownHosts := sshBasePath + "known_hosts"
	privateKeyFile := sshBasePath + sshPkName
	pkBuf, err := utl.ReadFile(privateKeyFile)
	if err != nil {
		log.Fatal(err)
	}

	signer, err := ssh.ParsePrivateKey(pkBuf)
	if err != nil {
		log.Fatal(err)
	}

	knHtsFile, err := kh.New(knownHosts)
	if err != nil {
		log.Fatal(err)
	}

	sshClientConfig := &ssh.ClientConfig{
		User: sshUserName,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: knHtsFile,
	}
	return sshClientConfig
}

func DialHost(sshConfig *ssh.ClientConfig, sshAddr string, sshPort int8) *ssh.Client {
	cClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", sshAddr, sshPort), sshConfig)
	if err != nil {
		log.Fatal("unable to establish any connection: ", err)
	}
	return cClient
}

func ConnHost(connect *ssh.Client, sshAddr string) net.Conn {
	conn, err := connect.Dial("tcp", sshAddr)
	if err != nil {
		log.Fatal("unable to connect: ", err)
	}
	return conn
}

func MakeNewClientConn(remoteConn net.Conn, remoteAddr string, remoteConfig *ssh.ClientConfig) (*ssh.Client, error) {
	// Establish SSH connection to remote host through jump host connection
	remoteClient, remoteChannels, remoteRequests, err := ssh.NewClientConn(remoteConn, remoteAddr, remoteConfig)
	if err != nil {
		fmt.Printf("Failed to establish SSH connection to remote host through jump host: %v\n", err)
		return nil, err
	}
	// Create SSH client for remote host
	remoteSSHClient := ssh.NewClient(remoteClient, remoteChannels, remoteRequests)
	return remoteSSHClient, nil
}

func CreateSession(connect *ssh.Client) (*ssh.Session, io.Reader) {
	session, err := connect.NewSession()
	if err != nil {
		log.Fatal("unable to create session: ", err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	// Request pseudo terminal
	if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
		log.Fatal("request for pseudo terminal failed: ", err)
	}

	StdOutPipe, err := session.StdoutPipe()
	if err != nil {
		log.Fatalf("Unable to setup stdout for session: %v", err)
	}
	return session, StdOutPipe
}

func FireCommand(client *ssh.Client, commands ...string) (map[string]string, error) {

	output := make(map[string]string)
	defer client.Close()

	for _, cmd := range commands {
		sshSession, StdOutPipe := CreateSession(client)
		defer sshSession.Close()
		if cmOut, err := sshSession.CombinedOutput(cmd); err != nil {
			log.Fatal(cmd, err)
			return nil, err
		} else {
			go func() {
				if _, err := io.Copy(os.Stdout, StdOutPipe); err != nil {
					log.Fatalf("Error copying stdout: %v", err)
				}
			}()
			output[cmd] = string(cmOut)
		}
	}
	return output, nil
}
