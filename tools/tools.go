package tools

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	//	"golang.org/x/crypto/ssh"
)

// https://golang.org/src/os/exec/example_test.go
func Shellx(cmds []string) {
	log.Fatal(strings.Join(cmds, " "))
	cmd := exec.Command(cmds[0], cmds[1:]...)
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	log.Printf("Waiting for command to finish...")
	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}

}

func Shell(cmds []string) error {
	fmt.Printf("%s\n", strings.Join(cmds, " "))
	cmd := exec.Command(cmds[0], cmds[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	fmt.Printf("%s\n", string(out))
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

/*
//https://github.com/SUSE/skuba/blob/master/internal/pkg/skuba/deployments/ssh/ssh.go
type Target struct {
	target       *deployments.Target
	user         string
	targetName   string
	sudo         bool
	port         int
	verboseLevel string
	client       *ssh.Client
}

func (t *Target) Ssh(silent bool, stdin string, command string, args ...string) (stdout string, stderr string, error error) {
	if t.client == nil {
		if err := t.initClient(); err != nil {
			return "", "", errors.Wrap(err, "failed to initialize client")
		}
	}
	session, err := t.client.NewSession()
	if err != nil {
		return "", "", err
	}
	if len(stdin) > 0 {
		session.Stdin = bytes.NewBufferString(stdin)
	}
	stdoutReader, err := session.StdoutPipe()
	if err != nil {
		return "", "", err
	}
	stderrReader, err := session.StderrPipe()
	if err != nil {
		return "", "", err
	}
	finalCommand := strings.Join(append([]string{command}, args...), " ")
	if t.sudo {
		finalCommand = fmt.Sprintf("sudo sh -c '%s'", finalCommand)
	}
	if !silent {
		klog.V(2).Infof("running command: %q", finalCommand)
	}
	if err := session.Start(finalCommand); err != nil {
		return "", "", err
	}
	stdoutChan := make(chan string)
	stderrChan := make(chan string)
	go readerStreamer(stdoutReader, stdoutChan, "stdout", silent)
	go readerStreamer(stderrReader, stderrChan, "stderr", silent)
	if err := session.Wait(); err != nil {
		return "", "", err
	}
	stdout = <-stdoutChan
	stderr = <-stderrChan
	return stdout, stderr, nil
}
*/
