package utl

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	yl "gopkg.in/yaml.v3"
)

// Opens a file and makes it available in byte array
func LoadFile(flPth string) []byte {
	fmt.Println("Loading " + flPth)
	flBf, err := os.ReadFile(flPth)
	if err != nil {
		log.Fatal(err)
		return nil
	} else {
		return flBf
	}
}

// takes a host name and gives out full domain name and host name
func (hj *HJSShConfig) MapHostDc(hostname string) error {
	if hostname == "" {
		return errors.ErrUnsupported
	} else {

		tmpVar := strings.Split(hostname, "-")

		hj.Dx.DataCenter = tmpVar[1]

		hj.HostAuth.Name =
			fmt.Sprintf("%v-%v%v",
				tmpVar[0],
				tmpVar[1],
				hj.Dx.DomainName)

		hj.BastionAuth.Name =
			fmt.Sprintf("%v-%v%v",
				hj.BastionAuth.Prefix,
				hj.Dx.DataCenter,
				hj.Dx.DomainName)
		return nil
	}
}

func (hj *HJSShConfig) DisplayHostDetails() {
	fmt.Printf(
		"HostName:\t%v\nBastionName:\t%v\n",
		hj.HostAuth.Name,
		hj.BastionAuth.Name)
}

func SourceHostName() string {
	var hostName, confirm string
	reader := bufio.NewReader(os.Stdin)
	for {
		if confirm != "y" {
			reader.Reset(reader)
			fmt.Print("Enter hostname: ")
			name, _ := reader.ReadString('\n')
			hostName = strings.TrimSpace(name)
			fmt.Printf("Confirm if hostName is %v (y/n): ", hostName)
			cfm, _ := reader.ReadString('\n')
			confirm = strings.TrimSpace(cfm)
		} else {
			break
		}
	}
	return hostName
}

func ParseCfg(yamlBuf []byte) HJSShConfig {
	var (
		sshConfig   SShCfg
		hostAuth    HostAuth
		bastionAuth BastionAuth
		dcdn        DCDN
		hj          HJSShConfig
	)
	err := yl.Unmarshal(yamlBuf, &sshConfig)
	if err != nil {
		log.Fatalln(err)
	}
	err = yl.Unmarshal(yamlBuf, &hostAuth)
	if err != nil {
		log.Fatalln(err)
	}
	err = yl.Unmarshal(yamlBuf, &bastionAuth)
	if err != nil {
		log.Fatalln(err)
	}
	err = yl.Unmarshal(yamlBuf, &dcdn)
	if err != nil {
		log.Fatalln(err)
	}
	hj.SSHConfig = &sshConfig
	hj.HostAuth = &hostAuth
	hj.BastionAuth = &bastionAuth
	hj.Dx = &dcdn
	hstName := SourceHostName()
	err = hj.MapHostDc(hstName)
	if err != nil {
		log.Fatalf("%v", err)
	}
	hj.DisplayHostDetails()
	return hj
}
