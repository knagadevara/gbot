package utl

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	ssh "golang.org/x/crypto/ssh"
)

func (sh *SShCfg) LoadSshConfig(userName string) *ssh.ClientConfig {

	// knownHosts := sshBasePath + "known_hosts"
	// knHtsFile, err := kh.New(knownHosts)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	pkBuf := LoadFile(sh.Path + sh.PK)
	signer, err := ssh.ParsePrivateKey(pkBuf)
	if err != nil {
		log.Fatal(err)
	}

	sshClientConfig := &ssh.ClientConfig{
		User: userName,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		// HostKeyCallback: knHtsFile,
	}
	return sshClientConfig
}

func DialHost(sshConfig *ssh.ClientConfig, sshAddr string, sshPort int8) (*ssh.Client, error) {
	if cClient, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", sshAddr, sshPort), sshConfig); err != nil {
		log.Fatal("unable to establish any connection: ", err)
		return nil, err
	} else {
		fmt.Println("Connection Established..")
		return cClient, nil
	}
}

func ConnHost(jumpCLient *ssh.Client, sshAddr string) (net.Conn, *ssh.Client, error) {
	if hostConnection, err := jumpCLient.Dial("tcp", fmt.Sprintf("%s:22", sshAddr)); err != nil {
		log.Fatal("unable to Establish Persistant Connection: ", err)
		return nil, nil, err
	} else {
		fmt.Println("Persistant Connection Established to Jump box...")
		return hostConnection, jumpCLient, nil
	}
}

func MakeNewClientConn(remoteHostConn net.Conn, jmpBxSshClient *ssh.Client, remoteAddr string, remoteConfig *ssh.ClientConfig) (rmtHstSshClt, JumpSshClient *ssh.Client, err error) {
	// Establish SSH connection to remote host through jump host connection
	rmtHstSshConn, remoteChannels, remoteRequests, err := ssh.NewClientConn(remoteHostConn, remoteAddr, remoteConfig)
	if err != nil {
		fmt.Printf("Failed to establish SSH connection to remote host through jump host: %v\n", err)
		return nil, nil, err
	}
	// Create SSH client for remote host
	rmtHstSshClient := ssh.NewClient(rmtHstSshConn, remoteChannels, remoteRequests)
	fmt.Println("Connection Established to Host Established...")
	return rmtHstSshClient, jmpBxSshClient, nil
}

func CreateSession(connect *ssh.Client) (*ssh.Session, io.Reader) {
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}
	if session, err := connect.NewSession(); err != nil {
		log.Fatal("unable to create session: ", err)
	} else {
		if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
			log.Fatal("request for pseudo terminal failed: ", err)
		} else {
			if StdOutPipe, err := session.StdoutPipe(); err != nil {
				log.Fatalf("Unable to setup stdout for session: %v", err)
			} else {
				fmt.Println("Session Established...")
				return session, StdOutPipe
			}
		}
	}
	return nil, nil
}

func CreateSessionNoTrm(connect *ssh.Client) (*ssh.Session, io.Reader) {
	session, err := connect.NewSession()
	if err != nil {
		log.Fatal("unable to create session: ", err)
	}
	StdOutPipe, err := session.StdoutPipe()
	if err != nil {
		log.Fatalf("Unable to setup stdout for session: %v", err)
	}
	return session, StdOutPipe
}

func FireCommands(rmtHstSshClient *ssh.Client, commands ...string) (map[string]string, error) {

	output := make(map[string]string)

	for _, cmd := range commands {
		sshSession, StdOutPipe := CreateSessionNoTrm(rmtHstSshClient)
		defer sshSession.Close()
		cmOut, err := sshSession.CombinedOutput(cmd)
		if err != nil {
			log.Fatal(cmd, err)
			return nil, err
		}
		go func() {
			if _, err := io.Copy(os.Stdout, StdOutPipe); err != nil {
				log.Fatalf("Error copying stdout: %v", err)
			}
			fmt.Printf("%v->%v\n", cmd, string(cmOut))
		}()
		output[cmd] = string(cmOut)

	}
	return output, nil
}

func (hj *HJSShConfig) CreateSshClientJumpHost() (rmtHstSshClt, JumpSshClient *ssh.Client, err error) {

	sshJumpConfig := hj.SSHConfig.LoadSshConfig(hj.BastionAuth.Uname)
	sshRmtConfig := hj.SSHConfig.LoadSshConfig(hj.HostAuth.Uname)

	sshJumpClient, err := DialHost(
		sshJumpConfig,
		hj.BastionAuth.Name,
		hj.SSHConfig.Port)
	if err != nil {
		log.Fatalln(err)
		return nil, nil, err
	}

	hostConnection, jumpCLient, err := ConnHost(sshJumpClient, hj.HostAuth.Uname)
	if err != nil {
		log.Fatalln(err)
		return nil, nil, err
	}

	rmtHstSshClient, jmpBxSshClient, err := MakeNewClientConn(
		hostConnection,
		jumpCLient,
		hj.HostAuth.Uname,
		sshRmtConfig)

	if err != nil {
		log.Fatalln(err)
		return nil, nil, err
	}
	return rmtHstSshClient, jmpBxSshClient, nil

}

func (hj *HJSShConfig) CreateSshClientHost() (*ssh.Client, error) {
	sshRmtConfig := hj.SSHConfig.LoadSshConfig(hj.HostAuth.Uname)
	sshClient, err := DialHost(sshRmtConfig,
		hj.HostAuth.Name,
		hj.SSHConfig.Port)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	return sshClient, nil
}
