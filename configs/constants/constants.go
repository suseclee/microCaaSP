package constants

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
)

// microCaaSP constants
const (
	NODEIP       = "192.168.2.0"
	USERNAME     = "sles"
	URL          = "http://10.84.128.39/repo/SUSE/Images/microCaaSP/"
	PASSPHRASE   = "suse"
	PASSWORDFILE = "/configs/constants/.passwd"
	VIRSHNETWORK = "microCaaSP-network"
	VIRSHDOMAIN  = "microCaaSP"
	VIRSHPOOL    = "microCaaSP"
)

func GetTempDir() string {
	return path.Join("/tmp", "microCaaSP")
}

func GetBackupDir() string {
	return path.Join(os.Getenv("HOME"), ".microCaaSP")
}

func GetDownloadFiles() []string {
	//[0] must be networkFileName
	return []string{"microCaaSP.xml", "microCaaSP.qcow2"}
}

func GetPassword() string {
	return string(decryptFile(PASSWORDFILE, PASSPHRASE))
}

func createHash(key string) string {
	hasher := md5.New()
	hasher.Write([]byte(key))
	return hex.EncodeToString(hasher.Sum(nil))
}

func decrypt(data []byte, passphrase string) []byte {
	key := []byte(createHash(passphrase))
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return plaintext
}

func decryptFile(filename string, passphrase string) []byte {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadFile(path.Join(dir, filename))
	if err != nil {
		log.Fatal(err)

	}
	return decrypt(data, passphrase)
}
