package main

import (
	"bufio"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"os/exec"
	"strconv"
	"strings"
)

var (
	cnvrgUser = "cnvrg"
	home      = "/home/cnvrg"
	encPass   = "paMfuNMgwFAX2"
)

var clusterUpParams = []param{
	{name: "single-node", value: true, usage: "create single node K8s cnvrg cluster"},
}

var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "deploy single node cnvrg K8s cluster",
}

var clusterUpCmd = &cobra.Command{
	Use:   "up",
	Short: "bring up cnvrg single nodes k8s cluster",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Infof("deploying k8s cluster")
		//createUser()
		generateKeys()

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
	if _, err := os.Stat(home + "/.ssh/id_rsa"); os.IsNotExist(err) {
		return false
	}
	return true
}

func createUser() {
	if !isUserExists("cnvrg") {
		argUser := []string{"-m", "-d", home, "-s", home, "-p", encPass, "--groups", "docker", cnvrgUser}
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

	logrus.Info("Private Key generated")
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

	logrus.Info("Public key generated")
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

func generateKeys() {

	if isKeysExists() {
		return
	}
	bitSize := 2048
	sshKeysDir := home + "./ssh"

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
}
