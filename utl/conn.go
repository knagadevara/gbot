package utl

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net"
	"os"
	"strings"
	"sync"

	ssh "golang.org/x/crypto/ssh"
)

// Load Hostkey, create if known hosts does not exist.
func loadFile(name string, flag int, perm fs.FileMode) *os.File {
	binBf, err := os.OpenFile(name, flag, perm)
	if err != nil {
		log.Fatalf("Failed to parse public key: %v", err)
	}
	return binBf
}

// Scan the known_hosts file create a dictionary of types.
func (sh *SShCfg) scanKnownHostKeys(hostName string) map[string]string {
	var foundHostKeys = make(map[string]string)
	PbBf := loadFile(strings.Join([]string{sh.Path, "known_hosts"}, ""), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0400)
	defer PbBf.Close()
	PbScanr := bufio.NewScanner(PbBf)
	for PbScanr.Scan() {
		AllKnownKeys := strings.Split(PbScanr.Text(), "\r\n")
		for _, keyLine := range AllKnownKeys {
			keyWords := strings.Split(keyLine, " ")
			if hostName == keyWords[0] {
				foundHostKeys[keyWords[1]] = keyWords[2]
			} else {
				continue
			}
		}
	}
	return foundHostKeys
}

// create human-readable SSH-key strings
func keyString(k ssh.PublicKey) string {
	return fmt.Sprintf(k.Type() + " " + base64.StdEncoding.EncodeToString(k.Marshal()))
}

func (sh *SShCfg) WriteKey(hostName string, k ssh.PublicKey) {
	PbBf := loadFile(strings.Join([]string{sh.Path, "known_hosts"}, ""), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	defer PbBf.Close()
	_, err := PbBf.WriteString(fmt.Sprintf(hostName+" "+k.Type()+" "+base64.StdEncoding.EncodeToString(k.Marshal())) + "\n\n")
	if err != nil {
		log.Fatalf("Failed to Write: %v", err)
	}
	fmt.Printf("WARNING: SSH-key verification is *NOT* in effect: to fix, add this trustedKey: %q", keyString(k))
}

func (sh *SShCfg) trustedHostKeyCallback(hkeys map[string]string, hostName string) ssh.HostKeyCallback {
	if len(hkeys) > 0 {
		return func(_ string, _ net.Addr, k ssh.PublicKey) error {
			hKey, ok := hkeys[k.Type()]
			trustedKey := fmt.Sprintf(k.Type() + " " + hKey)
			ks := keyString(k)
			if ok {
				if trustedKey == ks {
					return nil
				}
			} else {
				sh.WriteKey(hostName, k)
				return nil
			}
			return fmt.Errorf("SSH-key verification: expected %q but got %q", trustedKey, ks)
		}
	} else {
		return func(_ string, _ net.Addr, k ssh.PublicKey) error {
			sh.WriteKey(hostName, k)
			return nil
		}
	}
}

// LoadSshConfig loads SSH configuration for a given user
func (sh *SShCfg) LoadSshConfig(userName, hostName string) *ssh.ClientConfig {
	var hostKeyCallBack ssh.HostKeyCallback
	// Load Private Key and make a signer
	PkBf := LoadFile(strings.Join([]string{sh.Path, sh.PK, ".pem"}, ""))
	PkSigner, err := ssh.ParsePrivateKey(PkBf)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}
	hostKeyCallBack = sh.trustedHostKeyCallback(sh.scanKnownHostKeys(hostName), hostName)
	if hostKeyCallBack == nil {
		hostKeyCallBack = ssh.InsecureIgnoreHostKey()
		fmt.Println("Insecure Logging")
	}
	sshClientConfig := &ssh.ClientConfig{
		User: userName,
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(PkSigner),
		},
		HostKeyCallback: hostKeyCallBack,
		HostKeyAlgorithms: []string{
			"ssh-rsa-cert-v01@openssh.com", "ssh-dss-cert-v01@openssh.com",
			"ecdsa-sha2-nistp256-cert-v01@openssh.com", "ecdsa-sha2-nistp384-cert-v01@openssh.com",
			"ecdsa-sha2-nistp521-cert-v01@openssh.com", "sk-ecdsa-sha2-nistp256-cert-v01@openssh.com",
			"ssh-ed25519-cert-v01@openssh.com", "sk-ssh-ed25519-cert-v01@openssh.com",
			"aes256-cbc", "aes128-cbc",
			"3des-cbc", "des-cbc",
			"ssh-rsa", "rsa-sha2-512",
			"rsa-sha2-256", "ecdsa-sha2-nistp256",
			"sk-ecdsa-sha2-nistp256@openssh.com", "sk-ssh-ed25519@openssh.com",
			"ssh-ed25519",
		},
	}
	return sshClientConfig
}

func DialHost(sshConfig *ssh.ClientConfig, sshAddr string, sshPort int8) (*ssh.Client, error) {
	addr := fmt.Sprintf("%s:%d", sshAddr, sshPort)
	fmt.Printf("Conneting to %v\n", addr)
	if cClient, err := ssh.Dial("tcp", addr, sshConfig); err != nil {
		log.Fatal("unable to establish any connection: ", err)
		return nil, errors.New(err.Error())
	} else {
		fmt.Println("Dialed first Host")
		return cClient, nil
	}
}

func ConnHost(jumpCLient *ssh.Client, sshAddr string, sshPort int8) (net.Conn, error) {
	addr := fmt.Sprintf("%s:%d", sshAddr, sshPort)
	fmt.Printf("Conneting to %v\n", addr)
	hostConnection, err := jumpCLient.Dial("tcp", addr)
	if err != nil {
		log.Fatal("unable to Establish Persistant Connection: ", err)
		return nil, errors.New(err.Error())
	}
	fmt.Println("Dialed Connection from first Host")
	return hostConnection, nil
}

func MakeNewClientConn(remoteHostConn net.Conn, remoteAddr string, remoteConfig *ssh.ClientConfig) (rmtHstSshClt *ssh.Client, err error) {
	addr := fmt.Sprintf("%s:22", remoteAddr)
	rmtHstSshConn, remoteChannels, remoteRequests, err := ssh.NewClientConn(remoteHostConn, addr, remoteConfig)
	if err != nil {
		log.Fatal("Failed to establish SSH connection to remote host through jump: ", err)
		return nil, errors.New(err.Error())
	}
	rmtHstSshClient := ssh.NewClient(rmtHstSshConn, remoteChannels, remoteRequests)
	fmt.Println("Created a Connection and forging connection")
	return rmtHstSshClient, nil
}

func CreateStdOutPipedSession(connect *ssh.Client) (*ssh.Session, io.Reader) {
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	session, err := connect.NewSession()
	if err != nil {
		log.Fatal("unable to create session: ", err)
	}
	if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
		log.Fatal("request for pseudo terminal failed: ", err)
	}
	StdOutPipe, err := session.StdoutPipe()
	if err != nil {
		log.Fatalf("Unable to setup stdout for session: %v", err)
	}
	return session, StdOutPipe
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

func CreateSession(connect *ssh.Client) *ssh.Session {
	modes := ssh.TerminalModes{
		ssh.ECHO:          0,
		ssh.TTY_OP_ISPEED: 14400,
		ssh.TTY_OP_OSPEED: 14400,
	}
	session, err := connect.NewSession()
	if err != nil {
		log.Fatal("unable to create session: ", err)
	}
	if err := session.RequestPty("xterm", 40, 80, modes); err != nil {
		log.Fatal("request for pseudo terminal failed: ", err)
	}
	return session
}

func RunCommandStdOut(rmtHstSshClient *ssh.Client, commands ...string) (map[string]string, error) {
	output := make(map[string]string)
	for _, cmd := range commands {
		sshSession, StdOutPipe := CreateSessionNoTrm(rmtHstSshClient)
		defer sshSession.Close()
		cmOut, err := sshSession.CombinedOutput(cmd)
		if err != nil {
			log.Fatal(cmd, err)
			return nil, errors.New(err.Error())
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

func ExecuteCommand(rmtHstSshClient *ssh.Client, commands ...string) chan map[string]string {
	sessionOutChan := make(chan map[string]string, len(commands))
	var wg sync.WaitGroup
	wg.Add(len(commands))
	executeCommand := func(hotcommand string) {
		defer wg.Done()
		session := CreateSession(rmtHstSshClient)
		defer session.Close()
		cmOut, err := session.CombinedOutput(hotcommand)
		if err != nil {
			log.Fatal(hotcommand, err)
		}
		sessionOutChan <- map[string]string{hotcommand: string(cmOut)}
	}
	for _, command := range commands {
		go executeCommand(command)
	}
	go func() {
		defer close(sessionOutChan)
		wg.Wait()
	}()
	return sessionOutChan
}

func (hj *HJSShConfig) CreateSshClientJumpHost() (rmtHstSshClt, JumpSshClient *ssh.Client, err error) {
	sshJumpConfig := hj.SSHConfig.LoadSshConfig(hj.BastionAuth.Uname, hj.BastionAuth.Name)
	sshRmtConfig := hj.SSHConfig.LoadSshConfig(hj.HostAuth.Uname, hj.HostAuth.Name)
	sshJumpClient, err := DialHost(
		sshJumpConfig,
		hj.BastionAuth.Name,
		hj.SSHConfig.Port)
	if err != nil {
		log.Fatalln(err)
		return nil, nil, errors.New(err.Error())
	}
	hostConnection, err := ConnHost(
		sshJumpClient,
		hj.HostAuth.Name,
		hj.SSHConfig.Port)
	if err != nil {
		log.Fatalln(err)
		return nil, nil, errors.New(err.Error())
	}
	rmtHstSshClient, err := MakeNewClientConn(
		hostConnection,
		hj.HostAuth.Name,
		sshRmtConfig)
	if err != nil {
		log.Fatalln(err)
		return nil, nil, errors.New(err.Error())
	}
	return rmtHstSshClient, sshJumpClient, nil
}

func (hj *HJSShConfig) CreateSshClientHost() (*ssh.Client, error) {
	sshConfig := hj.SSHConfig.LoadSshConfig(hj.HostAuth.Uname, hj.HostAuth.Name)
	sshClient, err := DialHost(sshConfig,
		hj.HostAuth.Name,
		hj.SSHConfig.Port)
	if err != nil {
		log.Fatalln(err)
		return nil, errors.New(err.Error())
	}
	return sshClient, nil
}

func (hj *HJSShConfig) JumpOrNot() (host, jump *ssh.Client) {
	if hj.SSHConfig.Jump {
		rmtHstSshClt, jumpSshClient, err := hj.CreateSshClientJumpHost()
		if err != nil {
			log.Fatalln(err)
		}
		return rmtHstSshClt, jumpSshClient
	} else {
		hstSshClt, err := hj.CreateSshClientHost()
		if err != nil {
			log.Fatalln(err)
		}
		return hstSshClt, nil
	}
}

func CloseConn(host, jump *ssh.Client) {
	if jump != nil {
		host.Close()
		jump.Close()
	} else {
		host.Close()
	}
}
