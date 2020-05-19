package microCaaSP

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/gookit/color"
	"github.com/suseclee/microCaaSP/configs/constants"
	"github.com/suseclee/microCaaSP/tools"
)

type MicroCaaSP struct {
	tempDir   string
	backupDir string
	url       string
	files     []string
	version   string
}

func (m *MicroCaaSP) Init(version string) {
	m.tempDir = constants.GetTempDir()
	m.backupDir = path.Join(constants.GetBackupDir(), version)
	m.files = constants.GetDownloadFiles()
	m.url = constants.URL
	m.version = version
}
func (m *MicroCaaSP) Deploy() {
	if tools.MicroCaaSPDomainExist() {
		log.Fatalf("%s is already deployed", constants.VIRSHDOMAIN)
	}

	fmt.Println("** This is a TP version and NOT intended for production usage. **")
	fmt.Println("** MicroCaaSP TP works only under SUSE VPN **")
	fmt.Println("** Deploying microCaaSP version ", m.version, " **")

	for _, fileName := range m.files {
		if err := tools.Download(m.backupDir, m.tempDir, m.url, m.version, fileName); err != nil {
			log.Fatalf("Error on downloading %s: (check your VPN)", fileName)
		}
	}

	tools.TerminateNetwork(constants.VIRSHNETWORK)
	tools.ActivateNetwork(path.Join(m.tempDir, m.files[0]), constants.VIRSHNETWORK)
	tools.InstallDomain(path.Join(m.tempDir, m.files[1]))

	fmt.Println("Getting microCaaSP IP ...")
	tools.WaitForMicroCaaSPNetworkReady()

	fmt.Println("microCaaSP is booting ...")
	err := tools.WaitForLogin()
	if err != nil {
		fmt.Printf("\n%s\n", err)
	} else {
		fmt.Printf("\nNow ready for %s\n", color.FgGreen.Render("microCaaSP login"))
	}
}

func (m *MicroCaaSP) Login() {
	if !tools.MicroCaaSPDomainExist() {
		log.Fatalf("Domain %s is not deployed. \nUse %s first", color.FgRed.Render(constants.VIRSHDOMAIN), color.FgGreen.Render("microCaaSP deploy"))
	}
	tools.Terminal()
}

func (m *MicroCaaSP) Destroy() {
	tools.TerminateDomain(constants.VIRSHDOMAIN)
	tools.TerminateNetwork(constants.VIRSHNETWORK)
	tools.TerminatePool(constants.VIRSHPOOL)
	if _, err := os.Stat(constants.GetTempDir()); !os.IsNotExist(err) {
		cleanTempDir := []string{"rm", "-r", constants.GetTempDir()}
		if _, err := tools.Shell(cleanTempDir, constants.DEBUGMODE); err != nil {
			log.Fatalf("Error on cleaning up %s", constants.GetTempDir())
		}
	}
	fmt.Println("microCaaSP is destroyed successfully")
}
