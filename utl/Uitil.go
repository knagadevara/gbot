package utl

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
)

func ReadFile(flPth string) ([]byte, error) {
	fmt.Println("Loading" + flPth)
	flBf, err := os.ReadFile(flPth)
	if err != nil {
		log.Fatal(err)
		return nil, err
	} else {
		return flBf, nil
	}
}

// ssh -J snagadev@jumpbox-fra1.linode.com root@bs21-fra1.linode.com
func GetHostDcJmp(hostname string) (map[string]string, error) {
	var infraDetails map[string]string
	domainName := ".linode.com"
	if hostname == "" {
		return nil, errors.ErrUnsupported
	} else {
		clusterName := strings.Split(hostname, "-")
		infraDetails["dc"] = clusterName[1]
		infraDetails["hostname"] = fmt.Sprintf("%v-%v%v", clusterName[0], clusterName[1], domainName)
		infraDetails["jumpbox"] = fmt.Sprintf("jumpbox-%v%v", clusterName[1], domainName)
		return infraDetails, nil
	}
}
