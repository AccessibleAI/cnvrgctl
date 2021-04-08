package cmd

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/markbates/pkger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	cnvrgUser = "cnvrg"
	home      = "/home/cnvrg"
	encPass   = "paMfuNMgwFAX2"
)

var ClusterUpParams = []Param{
	{Name: "single-node", Value: true, Usage: "create single node K8s cnvrg cluster"},
}

var ClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "deploy single node cnvrg K8s cluster",
}

var ClusterUpCmd = &cobra.Command{
	Use:   "up",
	Short: "bring up cnvrg single nodes k8s cluster",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Infof("deploying k8s cluster")
		//createUser()
		//generateKeys()
		//fixPermissions()
		saveRke()

	},
}

func isUserExists(user string) bool {

	file, err := os.Open("/etc/passwd")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	r := bufio.NewScanner(file)
	for r.Scan() {
		lines := r.Text()
		parts := strings.Split(lines, ":")
		if parts[0] == user {
			return true
		}
	}
	return false
}

func isKeysExists() bool {
	if _, err := os.Stat(home + "/.ssh/rke_id_rsa"); os.IsNotExist(err) {
		return false
	}
	return true
}

func createUser() {
	if !isUserExists("cnvrg") {
		argUser := []string{"-m", "-d", home, "-s", "/bin/bash", "-p", encPass, "--groups", "docker", cnvrgUser}
		userCmd := exec.Command("useradd", argUser...)
		if out, err := userCmd.CombinedOutput(); err != nil {
			logrus.Errorf("err: %v, there was an error by adding user cnvrg", err)
		} else {
			logrus.Info(string(out))
		}
	} else {
		logrus.Warn("skip user creation, cnvrg user already exists")
	}
}

func generatePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	// Private Key generation
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}

	// Validate Private Key
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}

	logrus.Info("private key generated")
	return privateKey, nil
}

func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	// Get ASN.1 DER format
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)

	// pem.Block
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}

	// Private key in PEM format
	privatePEM := pem.EncodeToMemory(&privBlock)

	return privatePEM
}

func generatePublicKey(privatekey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(privatekey)
	if err != nil {
		return nil, err
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)

	logrus.Info("public key generated")
	return pubKeyBytes, nil
}

func writeKeyToFile(keyBytes []byte, saveFileTo string) error {
	err := ioutil.WriteFile(saveFileTo, keyBytes, 0600)
	if err != nil {
		return err
	}
	logrus.Infof("key saved to: %s", saveFileTo)
	u, err := user.Lookup(cnvrgUser)
	if err != nil {
		logrus.Error(err)
		return err
	}
	uid, _ := strconv.Atoi(u.Uid)
	gid, _ := strconv.Atoi(u.Gid)
	err = os.Chown(saveFileTo, uid, gid)
	return nil
}

func createAuthorizedKeysFile() {
	src := home + "/.ssh/rke_id_rsa.pub"
	dst := home + "/.ssh/authorized_keys"
	in, err := os.Open(src)
	if err != nil {
		logrus.Fatal(err)
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		logrus.Fatal(err)
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		logrus.Fatal(err)
	}
	defer out.Close()
}

func generateKeys() {

	if isKeysExists() {
		logrus.Info("keys exists, no need to generate")
		return
	}
	bitSize := 2048
	sshKeysDir := home + "/.ssh"

	if err := os.MkdirAll(sshKeysDir, os.ModePerm); err != nil {
		logrus.Fatalf("err: %v, faild to create %v", err, sshKeysDir)
	}

	privateKey, err := generatePrivateKey(bitSize)
	if err != nil {
		logrus.Fatalf("err: %v, error generating private key", err)
	}

	publicKeyBytes, err := generatePublicKey(&privateKey.PublicKey)
	if err != nil {
		logrus.Fatalf("err: %v, error generating public key", err)
	}

	privateKeyBytes := encodePrivateKeyToPEM(privateKey)

	err = writeKeyToFile(privateKeyBytes, sshKeysDir+"/rke_id_rsa")
	if err != nil {
		log.Fatal(err.Error())
	}

	err = writeKeyToFile([]byte(publicKeyBytes), sshKeysDir+"/rke_id_rsa.pub")
	if err != nil {
		log.Fatal(err.Error())
	}

	createAuthorizedKeysFile()
}

func getUserUidGid() (uid, gid int) {
	u, err := user.Lookup(cnvrgUser)
	if err != nil {
		logrus.Fatal(err)
	}
	uid, _ = strconv.Atoi(u.Uid)
	gid, _ = strconv.Atoi(u.Gid)
	return uid, gid
}

func fixPermissions() {
	logrus.Info("fixing permissions")
	uid, gid := getUserUidGid()
	err := filepath.Walk(home, func(name string, info os.FileInfo, err error) error {
		if err == nil {
			err = os.Chown(name, uid, gid)
		}
		return err
	})
	if err != nil {
		logrus.Fatal(err)
	}

}

func saveRke() {
	dst := home + "/rke"
	f, err := pkger.Open("/pkg/assets/rke_linux-amd64")
	if err != nil {
		logrus.Fatal(err)
	}
	destination, err := os.Create(dst)
	if err != nil {
		logrus.Fatal(err)
	}
	defer destination.Close()
	_, err = io.Copy(destination, f)
	if err != nil {
		logrus.Fatal(err)
	}
	uid, gid := getUserUidGid()
	if err = os.Chown(dst, uid, gid); err != nil {
		logrus.Fatal(err)
	}
	if err := os.Chmod(dst, 0755); err != nil {
		logrus.Fatal(err)
	}

}
