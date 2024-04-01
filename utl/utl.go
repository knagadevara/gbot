package utl

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/knagadevara/gbot/stypes"
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
func GetHostDcJmp(hostname string) (stypes.HostStr, error) {
	var infraDetails stypes.HostStr
	domainName := ".linode.com"
	if hostname == "" {
		return infraDetails, errors.ErrUnsupported
	} else {
		clusterName := strings.Split(hostname, "-")
		infraDetails.DataCenter = clusterName[1]
		infraDetails.HostName = fmt.Sprintf("%v-%v%v", clusterName[0], clusterName[1], domainName)
		infraDetails.JumpBox = fmt.Sprintf("jumpbox-%v%v", clusterName[1], domainName)
		infraDetails.DomainName = domainName
		return infraDetails, nil
	}
}

func DisplayHostDetails(hst stypes.HostStr) {
	fmt.Printf("HostName:\t%v\nJumpName:\t%v\n", hst.HostName, hst.JumpBox)
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

func ParseYaml(YmlPth string) stypes.SShCfg {

	yamlBuf := LoadFile(YmlPth)
	var sshConfig stypes.SShCfg
	err := yl.Unmarshal(yamlBuf, &sshConfig)
	if err != nil {
		log.Fatalln(err)
	}
	return sshConfig
}
