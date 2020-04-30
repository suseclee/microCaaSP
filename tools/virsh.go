package tools

import (
	"log"
)

func ActivateNetwork(networkFilePath string, networkName string) {

	virshCmd := []string{"virsh", "net-define", networkFilePath}
	if err := Shell(virshCmd); err != nil {
		log.Fatal(err)
	}

	virshCmd = []string{"virsh", "net-start", networkName}
	if err := Shell(virshCmd); err != nil {
		log.Fatal(err)
	}

	virshCmd = []string{"virsh", "net-autostart", networkName}
	if err := Shell(virshCmd); err != nil {
		log.Fatal(err)
	}
}

func InstallDomain(imagePath string) {
	virshCmd := []string{"virt-install", "--connect", "qemu:///system",
		"--virt-type", "kvm", "--name", "microCaaSP", "--ram", "4056", "--vcpus=4",
		"--os-type", "linux", "--os-variant", "sle15", "--disk", "path=" + imagePath + ",format=qcow2",
		"--import", "--network", "network=microCaaSP-network,mac=52:54:00:9e:1d:ed", "--noautoconsole"}

	if err := Shell(virshCmd); err != nil {
		log.Fatal(err)
	}
}
func TerminateDomain(domain string) {
	virshCmd := []string{"virsh", "shutdown", "--domain", domain}
	if err := Shell(virshCmd); err != nil {
		log.Fatal(err)
	}

	virshCmd = []string{"virsh", "destroy", "--domain", domain}
	if err := Shell(virshCmd); err != nil {
		log.Fatal(err)
	}

	virshCmd = []string{"virsh", "undefine", "--domain", domain}
	if err := Shell(virshCmd); err != nil {
		log.Fatal(err)
	}
}

func TerminateNetwork(network string) {
	virshCmd := []string{"virsh", "net-destroy", network}
	if err := Shell(virshCmd); err != nil {
		log.Fatal(err)
	}

	virshCmd = []string{"virsh", "net-undefine", network}
	if err := Shell(virshCmd); err != nil {
		log.Fatal(err)
	}
}

func TerminatePool(pool string) {
	virshCmd := []string{"virsh", "pool-destroy", pool}
	if err := Shell(virshCmd); err != nil {
		log.Fatal(err)
	}

	virshCmd = []string{"virsh", "pool-undefine", pool}
	if err := Shell(virshCmd); err != nil {
		log.Fatal(err)
	}
}
