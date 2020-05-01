package microCaaSP

import (
	"log"
	"path"

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

	if tools.MicroCaaSPDomainExist() {
		log.Fatalf("%s is already deployed", constants.VIRSHDOMAIN)
	}

	for _, fileName := range m.files {
		if err := tools.Download(m.backupDir, m.tempDir, m.url, fileName); err != nil {
			log.Fatalf("Error on downloading %s: (check your VPN)", fileName)
		}
	}

	tools.TerminateNetwork(constants.VIRSHNETWORK)
	tools.ActivateNetwork(path.Join(m.tempDir, m.files[0]), constants.VIRSHNETWORK)
	tools.InstallDomain(path.Join(m.tempDir, m.files[1]))
}

func (m *MicroCaaSP) Login() {
	if !tools.MicroCaaSPDomainExist() {
		log.Fatalf("%s is not deployed. Deploy first and login", constants.VIRSHDOMAIN)
	}
	tools.Terminal()
}

func (m *MicroCaaSP) Destroy() {
	tools.TerminateDomain(constants.VIRSHDOMAIN)
	tools.TerminateNetwork(constants.VIRSHNETWORK)
	tools.TerminatePool(constants.VIRSHPOOL)

	cleanTempDir := []string{"rm", "-r", constants.GetTempDir()}
	tools.Shell(cleanTempDir, constants.DEBUGMODE)
}
