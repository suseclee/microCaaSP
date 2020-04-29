package microCaaSP

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/libvirt/libvirt-go"
	"github.com/suseclee/microCaaSP/tools"
)

type MicroCaaSP struct {
	tempDir   string
	backupDir string
	ip        string
	username  string
	password  string
	url       string
	files     []string
}

func (m *MicroCaaSP) Init() {
	m.tempDir = path.Join("/tmp", "microCaaSP")
	m.backupDir = path.Join(os.Getenv("HOME"), ".microCaaSP")
	m.files = []string{"microCaaSP.xml", "microCaaSP.qcow2"}
	m.ip = "192.168.2.0"
	m.username = "sles"
	m.password = "linux"
	m.url = "http://10.84.128.39/repo/SUSE/Images/microCaaSP/"
}
func (m *MicroCaaSP) Deploy() {
	m.Init()

	for _, fileName := range m.files {
		tools.Download(m.backupDir, m.tempDir, m.url, fileName)
	}

	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		log.Fatal("Error: connect to qemu:///system")
		os.Exit(1)
	}

	if net, err := conn.LookupNetworkByName("microCaaSP-network"); err == nil {
		fmt.Println("virsh cmds")
		net.Destroy()
		net.Undefine()
	}

	virshCmd := []string{"virsh", "net-define", path.Join(m.tempDir, "microCaaSP.xml")}
	if err := tools.Shell(virshCmd); err != nil {
		log.Fatal(err)
	}

	virshCmd = []string{"virsh", "net-start", "microCaaSP-network"}
	if err := tools.Shell(virshCmd); err != nil {
		log.Fatal(err)
	}

	virshCmd = []string{"virsh", "net-autostart", "microCaaSP-network"}
	if err := tools.Shell(virshCmd); err != nil {
		log.Fatal(err)
	}

	virshCmd = []string{"virt-install", "--connect", "qemu:///system",
		"--virt-type", "kvm", "--name", "microCaaSP", "--ram", "4056", "--vcpus=4",
		"--os-type", "linux", "--os-variant", "sle15", "--disk", "path=" + path.Join(m.tempDir, "microCaaSP.qcow2") + ",format=qcow2",
		"--import", "--network", "network=microCaaSP-network,mac=52:54:00:9e:1d:ed", "--noautoconsole"}

	if err := tools.Shell(virshCmd); err != nil {
		log.Fatal(err)
	}
}

func (m *MicroCaaSP) Login() {

}

func (m *MicroCaaSP) Destroy() {
	m.Init()

	virshCmd := []string{"virsh", "shutdown", "--domain", "microCaaSP"}
	if err := tools.Shell(virshCmd); err != nil {
		log.Fatal(err)
	}

	virshCmd = []string{"virsh", "destroy", "--domain", "microCaaSP"}
	if err := tools.Shell(virshCmd); err != nil {
		log.Fatal(err)
	}

	virshCmd = []string{"virsh", "undefine", "--domain", "microCaaSP"}
	if err := tools.Shell(virshCmd); err != nil {
		log.Fatal(err)
	}

	virshCmd = []string{"virsh", "net-destroy", "microCaaSP-network"}
	if err := tools.Shell(virshCmd); err != nil {
		log.Fatal(err)
	}

	virshCmd = []string{"virsh", "net-undefine", "microCaaSP-network"}
	if err := tools.Shell(virshCmd); err != nil {
		log.Fatal(err)
	}

	virshCmd = []string{"virsh", "pool-destroy", "microCaaSP"}
	if err := tools.Shell(virshCmd); err != nil {
		log.Fatal(err)
	}

	virshCmd = []string{"virsh", "pool-undefine", "microCaaSP"}
	if err := tools.Shell(virshCmd); err != nil {
		log.Fatal(err)
	}

	cleanTempDir := []string{"rm", "-r", m.tempDir}
	tools.Shell(cleanTempDir)

}
