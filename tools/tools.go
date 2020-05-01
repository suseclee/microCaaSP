package tools

import (
	"fmt"
	"os"
	"os/exec"
	"path"
)

func Shell(cmds []string) error {
	// https://golang.org/src/os/exec/example_test.go
	cmd := exec.Command(cmds[0], cmds[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Printf("%s", string(out))
	return nil
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
	//log.Fatalf("%s : %s ", tempFilePath, backupFilePath)
	if _, err := os.Stat(tempFilePath); os.IsNotExist(err) {
		if _, err := os.Stat(backupFilePath); os.IsNotExist(err) {
			getDownload := []string{"wget", downloadFileURL, "-P", backupDir}
			Shell(getDownload)
			cp := []string{"cp", backupFilePath, tempDir + "/"}
			Shell(cp)
		} else {
			cp := []string{"cp", backupFilePath, tempDir + "/"}
			Shell(cp)
		}
	}
	return nil
}
