package utl

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

type hostStr struct {
	hostName   string
	domainName string
	dataCenter string
	jumpBox    string
}

// Opens a file and makes it available in byte array
func ReadFile(flPth string) ([]byte, error) {
	fmt.Println("Loading " + flPth)
	flBf, err := os.ReadFile(flPth)
	if err != nil {
		log.Fatal(err)
		return nil, err
	} else {
		return flBf, nil
	}
}

// takes a host name and gives out full domain name and host name
func GetHostDcJmp(hostname string) (hostStr, error) {
	var infraDetails hostStr
	domainName := ".linode.com"
	if hostname == "" {
		return infraDetails, errors.ErrUnsupported
	} else {
		clusterName := strings.Split(hostname, "-")
		infraDetails.dataCenter = clusterName[1]
		infraDetails.hostName = fmt.Sprintf("%v-%v%v", clusterName[0], clusterName[1], domainName)
		infraDetails.jumpBox = fmt.Sprintf("jumpbox-%v%v", clusterName[1], domainName)
		infraDetails.domainName = domainName
		return infraDetails, nil
	}
}

func DisplayHostDetails(hst hostStr) {
	fmt.Printf("HostName:\t%v\nJumpName:\t%v\n", hst.hostName, hst.jumpBox)
}

func SourceHostName() {

}
