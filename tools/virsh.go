package tools

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"

	"github.com/libvirt/libvirt-go"
	libvirtxml "github.com/libvirt/libvirt-go-xml"
	"github.com/suseclee/microCaaSP/configs/constants"
)

func ActivateNetwork(networkFilePath string, networkName string) {
	debug := constants.DEBUGMODE
	virshCmd := []string{"virsh", "net-define", networkFilePath}
	if _, err := Shell(virshCmd, debug); err != nil {
		log.Fatal(err)
	}
	virshCmd = []string{"virsh", "net-start", networkName}
	if _, err := Shell(virshCmd, debug); err != nil {
		log.Fatal(err)
	}

	conn := GetConnection()
	if net, errp := conn.LookupNetworkByName(constants.VIRSHNETWORK); errp == nil {
		if err := net.SetAutostart(true); err != nil {
			log.Fatal(err)
		}
	}
}

func GetMacAddress() string {
	xmlFile, err := os.Open(constants.GetNetworkXMLPath())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Successfully Opened %s \n", constants.GetNetworkXMLPath())
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)
	domcfg := &libvirtxml.Network{}
	err = xml.Unmarshal(byteValue, domcfg)
	if err != nil {
		log.Fatal(err)
	}
	return domcfg.IPs[0].DHCP.Hosts[0].MAC
}

func InstallDomain(imagePath string) {
	debug := constants.DEBUGMODE
	virshCmd := []string{"virt-install", "--connect", "qemu:///system",
		"--virt-type", "kvm", "--name", "microCaaSP", "--ram", "4096", "--vcpus=4",
		"--os-type", "linux", "--os-variant", "sle15", "--disk", "path=" + imagePath + ",format=qcow2",
		"--import", "--network", "network=" + constants.VIRSHNETWORK + ",mac=" + GetMacAddress(), "--noautoconsole"}

	out, err := Shell(virshCmd, debug)
	if err != nil {
		log.Printf("Err on '%s'", virshCmd)
		log.Fatal(err)
	}
	log.Printf("%s", out)
}

func TerminateDomain(domain string) {
	conn := GetConnection()
	if domain, errp := conn.LookupDomainByName(constants.VIRSHDOMAIN); errp == nil {
		if err := domain.Shutdown(); err != nil {
			log.Fatal(err)
		}
		if err := domain.Destroy(); err != nil {
			log.Fatal(err)
		}
		if err := domain.Undefine(); err != nil {
			log.Fatal(err)
		}
	}
}

func TerminateNetwork(network string) {
	conn := GetConnection()
	if net, errp := conn.LookupNetworkByName(constants.VIRSHNETWORK); errp == nil {
		if err := net.Destroy(); err != nil {
			log.Fatal(err)
		}
		if err := net.Undefine(); err != nil {
			log.Fatal(err)
		}
	}
}

func TerminatePool(pool string) {
	conn := GetConnection()
	if pool, errp := conn.LookupStoragePoolByName(constants.VIRSHPOOL); errp == nil {
		if err := pool.Destroy(); err != nil {
			log.Fatal(err)
		}
		if err := pool.Undefine(); err != nil {
			log.Fatal(err)
		}
	}
}

func GetConnection() *libvirt.Connect {
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		log.Fatal("Error: connect to qemu:///system")
	}
	return conn
}

func MicroCaaSPDomainExist() bool {
	conn := GetConnection()
	if _, err := conn.LookupDomainByName(constants.VIRSHDOMAIN); err == nil {
		return true
	}
	return false
}

func GetMicroCaaSPDomain() (*libvirt.Domain, error) {
	conn := GetConnection()
	return conn.LookupDomainByName(constants.VIRSHDOMAIN)
}

func MicroCaaSPNetworkExist() bool {
	conn := GetConnection()
	if _, err := conn.LookupNetworkByName(constants.VIRSHNETWORK); err == nil {
		return true
	}
	return false
}

func GetMicroCaaSPNetwork() (*libvirt.Network, error) {
	conn := GetConnection()
	return conn.LookupNetworkByName(constants.VIRSHNETWORK)
}

func MicroCaaSPStoragePoolExist() bool {
	conn := GetConnection()
	if _, err := conn.LookupStoragePoolByName(constants.VIRSHPOOL); err == nil {
		return true
	}
	return false
}

func WaitForMicroCaaSPNetworkReady() error {
	net, err := GetMicroCaaSPNetwork()
	if err != nil {
		log.Fatal(err)
	}

	for i := 0; i < 20; i++ {
		leases, _ := net.GetDHCPLeases()
		if leases != nil && len(leases) > 0 {
			return nil
		}
		ShellSpin([]string{"sleep", strconv.Itoa(2)})
	}
	return fmt.Errorf("Error: Network is not ready")
}

func GetIP() (string, error) {
	net, err := GetMicroCaaSPNetwork()
	if err != nil {
		log.Fatal(err)
	}

	if err := WaitForMicroCaaSPNetworkReady(); err != nil {
		log.Fatal(err)
	}

	networkDHCPLeases, err := net.GetDHCPLeases()
	for _, lease := range networkDHCPLeases {
		if lease.IPaddr != "" {
			return lease.IPaddr, nil
		}
	}
	return "", fmt.Errorf("Error for getting IP")

}
