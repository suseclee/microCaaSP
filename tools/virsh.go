package tools

import (
	"log"
	"os"

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

	virshCmd = []string{"virsh", "net-autostart", networkName}
	if err := Shell(virshCmd, debug); err != nil {
		log.Fatal(err)
	}
}

func InstallDomain(imagePath string) {
	debug := constants.DEBUGMODE
	virshCmd := []string{"virt-install", "--connect", "qemu:///system",
		"--virt-type", "kvm", "--name", "microCaaSP", "--ram", "4056", "--vcpus=4",
		"--os-type", "linux", "--os-variant", "sle15", "--disk", "path=" + imagePath + ",format=qcow2",
		"--import", "--network", "network=microCaaSP-network,mac=52:54:00:9e:1d:ed", "--noautoconsole"}

	if err := Shell(virshCmd, debug); err != nil {
		log.Fatal(err)
	}
}
func TerminateDomain(domain string) {
	debug := constants.DEBUGMODE
	virshCmd := []string{"virsh", "shutdown", "--domain", domain}
	if err := Shell(virshCmd, debug); err != nil {
		log.Fatal(err)
	}

	virshCmd = []string{"virsh", "destroy", "--domain", domain}
	if err := Shell(virshCmd, debug); err != nil {
		log.Fatal(err)
	}

	virshCmd = []string{"virsh", "undefine", "--domain", domain}
	if err := Shell(virshCmd, debug); err != nil {
		log.Fatal(err)
	}
}

func TerminateNetwork(network string) {
	debug := constants.DEBUGMODE
	virshCmd := []string{"virsh", "net-destroy", network}
	if err := Shell(virshCmd, debug); err != nil {
		log.Fatal(err)
	}

	virshCmd = []string{"virsh", "net-undefine", network}
	if err := Shell(virshCmd, debug); err != nil {
		log.Fatal(err)
	}
}

func TerminatePool(pool string) {
	debug := constants.DEBUGMODE
	virshCmd := []string{"virsh", "pool-destroy", pool}
	if err := Shell(virshCmd, debug); err != nil {
		log.Fatal(err)
	}

	virshCmd = []string{"virsh", "pool-undefine", pool}
	if err := Shell(virshCmd, debug); err != nil {
		log.Fatal(err)
	}
}

func DeletePreviousNetwork(network string) {
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		log.Fatal("Error: connect to qemu:///system")
		os.Exit(1)
	}
	if net, err := conn.LookupNetworkByName(constants.VIRSHNETWORK); err == nil {
		net.Destroy()
		net.Undefine()
	}
}

func CheckDomain(domain string) {
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		log.Fatal("Error: connect to qemu:///system")
		os.Exit(1)
	}
	if net, err := conn.LookupNetworkByName(constants.VIRSHNETWORK); err == nil {
		net.Destroy()
		net.Undefine()
	}
}
