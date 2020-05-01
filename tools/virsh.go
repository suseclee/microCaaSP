package tools

import (
	"log"

	"github.com/libvirt/libvirt-go"
	"github.com/suseclee/microCaaSP/configs/constants"
)

func ActivateNetwork(networkFilePath string, networkName string) {
	debug := constants.DEBUGMODE
	virshCmd := []string{"virsh", "net-define", networkFilePath}
	if err := Shell(virshCmd, debug); err != nil {
		log.Fatal(err)
	}
	virshCmd = []string{"virsh", "net-start", networkName}
	if err := Shell(virshCmd, debug); err != nil {
		log.Fatal(err)
	}

	conn := GetConnection()
	if net, errp := conn.LookupNetworkByName(constants.VIRSHNETWORK); errp == nil {
		if err := net.SetAutostart(true); err != nil {
			log.Fatal(err)
		}
	}
}

func InstallDomain(imagePath string) {
	debug := constants.DEBUGMODE
	virshCmd := []string{"virt-install", "--connect", "qemu:///system",
		"--virt-type", "kvm", "--name", "microCaaSP", "--ram", "4096", "--vcpus=4",
		"--os-type", "linux", "--os-variant", "sle15", "--disk", "path=" + imagePath + ",format=qcow2",
		"--import", "--network", "network=microCaaSP-network,mac=52:54:00:9e:1d:ed", "--noautoconsole"}

	if err := Shell(virshCmd, debug); err != nil {
		log.Fatal(err)
	}
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

func MicroCaaSPNetworkExist() bool {
	conn := GetConnection()
	if _, err := conn.LookupNetworkByName(constants.VIRSHDOMAIN); err == nil {
		return true
	}
	return false
}

func MicroCaaSPStoragePoolExist() bool {
	conn := GetConnection()
	if _, err := conn.LookupStoragePoolByName(constants.VIRSHPOOL); err == nil {
		return true
	}
	return false
}
