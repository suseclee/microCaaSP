package tools

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/suseclee/microCaaSP/configs/constants"

	"github.com/briandowns/spinner"
	"github.com/mattn/go-isatty"
	"k8s.io/klog"
)

func Shell(cmds []string, debug bool) (string, error) {
	// https://golang.org/src/os/exec/example_test.go
	cmd := exec.Command(cmds[0], cmds[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	if debug {
		fmt.Printf("%s", string(out))
	}
	return string(out), nil
}

func ShellSpin(cmds []string) error {

	spin := spinner.New(spinner.CharSets[36], 1000*time.Millisecond)
	isTerminal := isatty.IsTerminal(os.Stdout.Fd())

	cmd := exec.Command(cmds[0], cmds[1:]...)
	if err := cmd.Start(); err != nil {
		return err
	}
	if isTerminal {
		spin.Start()
	}
	if err := cmd.Wait(); err != nil {
		if isTerminal {
			spin.Stop()
		}
		log.Println(strings.Join(cmds, " "))
		return err
	}
	if isTerminal {
		spin.Stop()
	}
	return nil
}

func readerStreamer(reader io.Reader, outputChan chan<- string, description string) {
	result := bytes.Buffer{}
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		result.Write([]byte(scanner.Text()))
		if description == "stdout" {
			klog.V(2).Infof("%s", scanner.Text())
		} else if description == "stderr" {
			klog.Errorf("%s", scanner.Text())
		}
	}
	outputChan <- result.String()
}

func CreateImageStorage(dir string) {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		os.Mkdir(dir, 0755)
	}
}

func RemoveImageStorage(dir string) {
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		os.RemoveAll(dir)
	}
}

func Download(backupDir string, tempDir string, url string, version string, fileName string) error {
	tempFilePath := path.Join(tempDir, fileName)
	backupFilePath := path.Join(backupDir, fileName)
	downloadFileURL := url + path.Join(version, fileName)

	silent := true
	if _, err := os.Stat(tempFilePath); os.IsNotExist(err) {
		if _, err := os.Stat(backupFilePath); os.IsNotExist(err) {
			if filepath.Ext(fileName) == ".qcow2" {
				fmt.Printf("Downloading %s (~4GB : < 10 min)...\n", fileName)
			}
			CreateImageStorage(backupDir)
			getDownload := []string{"wget", downloadFileURL, "-P", backupDir}
			if err := ShellSpin(getDownload); err != nil {
				log.Println("Check your microCaaSP version. Avialable versions are ", strings.Join(constants.GetAvilableVersions(), ","))
				RemoveImageStorage(backupDir)
			}
			CreateImageStorage(tempDir)
			cp := []string{"cp", backupFilePath, tempDir + "/"}
			if _, e := Shell(cp, silent); e != nil {
				RemoveImageStorage(tempDir)
				log.Fatal(e)
			}
		} else {
			cp := []string{"cp", backupFilePath, tempDir + "/"}
			if _, e := Shell(cp, silent); e != nil {
				log.Fatal(e)
			}
		}
	}
	return nil
}
