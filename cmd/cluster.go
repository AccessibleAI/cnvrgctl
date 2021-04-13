package cmd

import (
	"bufio"
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/cnvrgctl/pkg/cnvrg"
	emoji "github.com/kyokomi/emoji/v2"
	"github.com/markbates/pkger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"net"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"
)

var (
	cnvrgUser     = "cnvrg"
	home          = "/home/cnvrg"
	sshPrivateKey = "/home/cnvrg/.ssh/id_rsa"
	sshPublicKey  = "/home/cnvrg/.ssh/id_rsa.pub"
	encPass       = "paMfuNMgwFAX2"
	rkeDir        = "/home/cnvrg/rke-cluster"
	sshdConfig    = "/etc/ssh/sshd_config"
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
		createUser()
		generateKeys()
		saveTools()
		generateRkeClusterManifest()
		allowTcpForwarding()
		fixPermissions()
		rkeUp()
		checkClusterReady()
	},
}

func allowTcpForwarding() {
	shouldEnableTcpForwarding := true
	sshdConfigData := []string{}
	file, err := os.Open(sshdConfig)
	if err != nil {
		logrus.Errorf("err: %v, error generating private key", err)
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		sshdConfigData = append(sshdConfigData, scanner.Text())
	}
	for i := 0; i < len(sshdConfigData); i++ {
		if strings.HasPrefix(sshdConfigData[i], "#") {
			continue
		}
		if strings.Contains(sshdConfigData[i], "AllowTcpForwarding") {
			if strings.TrimSpace(strings.ReplaceAll(sshdConfigData[i], "AllowTcpForwarding", "")) == "no" {
				sshdConfigData[i] = "AllowTcpForwarding yes"
				shouldEnableTcpForwarding = false
			}
			if strings.TrimSpace(strings.ReplaceAll(sshdConfigData[i], "AllowTcpForwarding", "")) == "yes" {
				shouldEnableTcpForwarding = false
			}
		}
	}
	if shouldEnableTcpForwarding {
		sshdConfigData = append(sshdConfigData, "AllowTcpForwarding yes")
	}
	file, err = os.OpenFile(sshdConfig, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModeAppend)
	if err != nil {
		logrus.Errorf("err: %v, error opening sshd config", err)
		panic(err)
	}
	if err := file.Truncate(0); err != nil {
		logrus.Errorf("err: %v, error truncating sshd config", err)
		panic(err)
	}
	if _, err := file.WriteString(strings.Join(sshdConfigData, "\n")); err != nil {
		logrus.Errorf("err: %v, error writing sshd config", err)
		panic(err)
	}
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
	if _, err := os.Stat(sshPrivateKey); os.IsNotExist(err) {
		return false
	}
	return true
}

func createUser() {
	if !isUserExists("cnvrg") {
		argUser := []string{"-m", "-d", home, "-s", "/bin/bash", "-p", encPass, "--groups", "docker,sudo", cnvrgUser}
		userCmd := exec.Command("useradd", argUser...)
		if out, err := userCmd.CombinedOutput(); err != nil {
			logrus.Errorf("err: %v, there was an error by adding user cnvrg", err)
			panic(err)
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
	src := sshPublicKey
	dst := home + "/.ssh/authorized_keys"
	in, err := os.Open(src)
	if err != nil {
		logrus.Fatal(err)
		panic(err)
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		logrus.Fatal(err)
		panic(err)
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		logrus.Fatal(err)
		panic(err)
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
		logrus.Errorf("err: %v, faild to create %v", err, sshKeysDir)
		panic(err)
	}

	privateKey, err := generatePrivateKey(bitSize)
	if err != nil {
		logrus.Errorf("err: %v, error generating private key", err)
		panic(err)
	}

	publicKeyBytes, err := generatePublicKey(&privateKey.PublicKey)
	if err != nil {
		logrus.Errorf("err: %v, error generating public key", err)
		panic(err)
	}

	privateKeyBytes := encodePrivateKeyToPEM(privateKey)

	err = writeKeyToFile(privateKeyBytes, sshPrivateKey)
	if err != nil {
		logrus.Error(err.Error())
		panic(err)
	}

	err = writeKeyToFile([]byte(publicKeyBytes), sshPublicKey)
	if err != nil {
		logrus.Error(err.Error())
		panic(err)
	}

	createAuthorizedKeysFile()
}

func getUserUidGid() (uid, gid int) {
	u, err := user.Lookup(cnvrgUser)
	if err != nil {
		logrus.Fatal(err)
		panic(err)
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
		panic(err)
	}

}

func saveTools() {
	binTools := []string{"k9s", "kubectl", "rke"}
	for _, toolName := range binTools {
		logrus.Infof("dumping %s", toolName)
		dst := "/usr/local/bin/" + toolName
		saveAsset(dst, toolName)
	}

	manifests := []string{"cnvrg-crds.yaml", "cnvrg-operator.yaml"}
	for _, manifest := range manifests {
		logrus.Infof("dumping %s", manifest)
		dst := rkeDir + "/" + manifest
		saveAsset(dst, manifest)
	}
}

func saveAsset(dst string, toolName string) {
	logrus.Infof("dumping %s", toolName)
	f, err := pkger.Open("/pkg/assets/" + toolName)
	if err != nil {
		logrus.Fatal(err)
		panic(err)
	}
	destination, err := os.Create(dst)
	if err != nil {
		logrus.Fatal(err)
		panic(err)
	}
	defer destination.Close()
	_, err = io.Copy(destination, f)
	if err != nil {
		logrus.Fatal(err)
		panic(err)
	}
	uid, gid := getUserUidGid()
	if err = os.Chown(dst, uid, gid); err != nil {
		logrus.Fatal(err)
		panic(err)
	}
	if err := os.Chmod(dst, 0755); err != nil {
		logrus.Fatal(err)
		panic(err)
	}
}

func getMainIp() string {
	var ipv4Addr net.IP

	nic, err := net.InterfaceByName(getMainNic())
	if err != nil {
		logrus.Errorf("%s can't get interface", err)
		panic(err)
	}
	addrs, err := nic.Addrs()
	if err != nil { // get addresses
		logrus.Errorf("%s can't get interface addesses", err)
		panic(err)
	}
	for _, addr := range addrs { // get ipv4 address
		ipv4Addr = addr.(*net.IPNet).IP.To4()
		if ipv4Addr != nil {
			break
		}
	}
	if ipv4Addr == nil {
		logrus.Errorf("interface does not have any IP addesses")
		panic(err)
	}
	logrus.Infof("node main IP: %s", ipv4Addr.String())
	return ipv4Addr.String()

}

func getMainNic() string {
	// /proc/net/route
	procRouteFile := "/proc/net/route"
	b, err := ioutil.ReadFile(procRouteFile)
	if err != nil {
		logrus.Error(err)
		panic(err)
	}
	routeData := strings.Split(string(b), "\n")
	if len(routeData) < 2 {
		logrus.Errorf("%s doesn't contains enougth information, %v", procRouteFile, routeData)
		panic(err)
	}
	nic := strings.Split(routeData[1], "\t")
	if len(nic) < 1 || nic[0] == "" {
		logrus.Errorf("%s doesn't contains enougth information, %v", procRouteFile, nic)
		panic(err)
	}
	logrus.Infof("detected node main nic: %s", nic[0])
	return nic[0]
}

func generateRkeClusterManifest() {
	var tpl bytes.Buffer
	templateData := map[string]interface{}{
		"Data": map[string]interface{}{
			"Server":        getMainIp(),
			"User":          cnvrgUser,
			"SshPrivateKey": sshPrivateKey,
		},
	}

	clusterManifestTpl := "/pkg/assets/cluster.tpl"
	f, err := pkger.Open(clusterManifestTpl)
	if err != nil {
		logrus.Errorf("error reading cluster.tpl %v", err)
		panic(err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		logrus.Errorf("%v, error reading file: %v", err, clusterManifestTpl)
		panic(err)
	}
	clusterTmpl, err := template.New(strings.ReplaceAll(clusterManifestTpl, "/", "-")).Parse(string(b))
	if err != nil {
		logrus.Errorf("%v, template: %v", err, clusterManifestTpl)
		panic(err)
	}
	if err = clusterTmpl.Execute(&tpl, templateData); err != nil {
		logrus.Errorf("err: %v rendering template error", err)
		panic(err)
	}

	// dir for rke
	if err := os.MkdirAll(rkeDir, os.ModePerm); err != nil {
		logrus.Errorf("err: %v, faild to create %v", err, rkeDir)
		panic(err)
	}

	if err := ioutil.WriteFile(rkeDir+"/cluster.yml", tpl.Bytes(), 0644); err != nil {
		logrus.Errorf("err: %v, faild to cluster.yml %v", err, rkeDir)
		panic(err)
	}
}

func rkeUp() {
	args := []string{"-c", fmt.Sprintf(`su - cnvrg -c "cd %s && rke up"`, rkeDir)}
	cmd := exec.Command("/bin/bash", args...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		logrus.Errorf("%v error creating StdoutPipe for cmd", os.Stderr)
		panic(err)
	}
	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			logrus.Infof("%s", scanner.Text())
		}
	}()
	if err := cmd.Start(); err != nil {
		logrus.Error(err)
		panic(err)
	}
	if err := cmd.Wait(); err != nil {
		logrus.Error(err)
		panic(err)
	}

	// copy kubeconfig file to default location
	args = []string{"-c", fmt.Sprintf(`su - cnvrg -c "cp %s/kube_config_cluster.yml %s/.kube/config"`, rkeDir, home)}
	cmd = exec.Command("/bin/bash", args...)
	if _, err := cmd.CombinedOutput(); err != nil {
		logrus.Error(err)
		panic(err)
	}
}

func checkClusterReady() {
	viper.Set("kubeconfig", fmt.Sprintf("%s/.kube/config", home))
	for i := 1; i <= 20; i++ {
		ready, err := cnvrg.CheckNodesReadyStatus()
		if err != nil {
			logrus.Errorf("err: %v, can't list K8s nodes", err)
			panic(err)
		}
		if ready {
			logrus.Infof("K8s is ready, time to deploy cnvrg %s! (run: cnvrgctl cnvrg up -h %s)", emoji.Sprint(":rocket:!"), emoji.Sprint(":nerd:"))
			break
		}
		logrus.Infof("checking k8s ready status, attempt left: %d ...", 20-i)
		time.Sleep(10 * time.Second)
	}
}
