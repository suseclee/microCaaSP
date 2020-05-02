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
	"time"

	"github.com/briandowns/spinner"
	"github.com/mattn/go-isatty"
	"k8s.io/klog"
)

func Shell(cmds []string, debug bool) error {
	// https://golang.org/src/os/exec/example_test.go
	cmd := exec.Command(cmds[0], cmds[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	if debug {
		fmt.Printf("%s", string(out))
	}
	return nil
}

func ShellSpin(cmds []string) {

	spin := spinner.New(spinner.CharSets[36], 1000*time.Millisecond)
	isTerminal := isatty.IsTerminal(os.Stdout.Fd())

	cmd := exec.Command(cmds[0], cmds[1:]...)
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	if isTerminal {
		spin.Start()
	}
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	if isTerminal {
		spin.Stop()
	}
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

func Download(backupDir string, tempDir string, url string, fileName string) error {
	CreateImageStorage(backupDir)
	CreateImageStorage(tempDir)
	tempFilePath := path.Join(tempDir, fileName)
	backupFilePath := path.Join(backupDir, fileName)
	downloadFileURL := url + fileName
	silent := true
	if _, err := os.Stat(tempFilePath); os.IsNotExist(err) {
		if _, err := os.Stat(backupFilePath); os.IsNotExist(err) {
			if filepath.Ext(fileName) == ".qcow2" {
				fmt.Printf("Downloading %s (~4GB : < 10 min)...\n", fileName)
			}
			getDownload := []string{"wget", downloadFileURL, "-P", backupDir}
			ShellSpin(getDownload)
			cp := []string{"cp", backupFilePath, tempDir + "/"}
			Shell(cp, silent)
		} else {
			cp := []string{"cp", backupFilePath, tempDir + "/"}
			Shell(cp, silent)
		}
	}
	return nil
}
