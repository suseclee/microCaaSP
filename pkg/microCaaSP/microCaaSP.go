package microCaaSP

import (
	"log"
	"os"
	"path"

	"github.com/libvirt/libvirt-go"
	"github.com/suseclee/microCaaSP/configs/constants"
	"github.com/suseclee/microCaaSP/tools"
)

type MicroCaaSP struct {
	tempDir   string
	backupDir string
	url       string
	files     []string
}

func (m *MicroCaaSP) Init() {
	m.tempDir = constants.GetTempDir()
	m.backupDir = constants.GetBackupDir()
	m.files = constants.GetDownloadFiles()
	m.url = constants.URL
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

	if net, err := conn.LookupNetworkByName(constants.VIRSHNETWORK); err == nil {
		net.Destroy()
		net.Undefine()
	}

	tools.ActivateNetwork(path.Join(m.tempDir, m.files[0]), constants.VIRSHNETWORK)
	tools.InstallDomain(path.Join(m.tempDir, m.files[1]))
}

func (m *MicroCaaSP) Login() {
	tools.Terminal()
}

func (m *MicroCaaSP) Destroy() {
	tools.TerminateDomain(constants.VIRSHDOMAIN)
	tools.TerminateNetwork(constants.VIRSHNETWORK)
	tools.TerminatePool(constants.VIRSHPOOL)

	cleanTempDir := []string{"rm", "-r", constants.GetTempDir()}
	tools.Shell(cleanTempDir)
}
