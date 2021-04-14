package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/cnvrgctl/pkg"
	"github.com/cnvrgctl/pkg/cnvrg"
	emoji "github.com/kyokomi/emoji/v2"
	"github.com/markbates/pkger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
	home          = "/home/cnvrg"
	sshPrivateKey = "/home/cnvrg/.ssh/id_rsa"
	rkeDir        = "/home/cnvrg/rke-cluster"
)

var ClusterUpParams = []Param{
	{Name: "single-node", Value: true, Usage: "create single node K8s cnvrg cluster"},
	{Name: "cnvrg-user", Value: "cnvrg", Usage: "user for managing cnvrg stack"},
	{Name: "install-docker", Value: true, Usage: "set to false to disable docker installation"},
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
		generateClusterSetupScript()
		if viper.GetBool("install-docker") {
			pkg.ExecBashScript("cluster-setup.sh installDocker")
		}
		pkg.ExecBashScript("cluster-setup.sh createUser")
		//dumpDeploymentAssets()
		//generateRkeClusterManifest()
		//fixPermissions()
		//rkeUp()
		//checkClusterReady()
	},
}

func generateClusterSetupScript() {
	templateData := map[string]interface{}{
		"Data": map[string]interface{}{
			"CnvrgUser": viper.GetString("cnvrg-user"),
		},
	}
	buffer, err := renderTemplate("/pkg/assets/cluster-setup.sh", templateData)
	if err != nil {
		logrus.Errorf("error generating cluster setup script err: %v", err)
		panic(err)
	}

	if viper.GetBool("dry-run") {
		logrus.Infof("\n%s", buffer.String())
		return
	}
	if err := ioutil.WriteFile("/usr/local/bin/cluster-setup.sh", buffer.Bytes(), 0755); err != nil {
		logrus.Errorf("err: %v, faild to save cluster-setup.sh %v", err, rkeDir)
		panic(err)
	}

}

func getUserUidGid() (uid, gid int) {
	u, err := user.Lookup(viper.GetString("cnvrg-user"))
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

func dumpDeploymentAssets() {
	manifests := []string{"cnvrg-crds.yaml", "cnvrg-operator.yaml"}
	for _, manifest := range manifests {
		logrus.Infof("dumping %s", manifest)
		dst := rkeDir + "/" + manifest
		logrus.Infof("dumping %s", manifest)
		f, err := pkger.Open("/pkg/assets/" + manifest)
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

func renderTemplate(templateFile string, templateData map[string]interface{}) (*bytes.Buffer, error) {
	var tpl bytes.Buffer
	//templateData := map[string]interface{}{
	//	"Data": map[string]interface{}{
	//		"Server":        getMainIp(),
	//		"User":          cnvrgUser,
	//		"SshPrivateKey": sshPrivateKey,
	//	},
	//}

	//clusterManifestTpl := "/pkg/assets/cluster.tpl"
	f, err := pkger.Open(templateFile)
	if err != nil {
		logrus.Errorf("error reading cluster.tpl %v", err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		logrus.Errorf("%v, error reading file: %v", err, templateFile)
		return nil, err
	}

	clusterTmpl, err := template.New(strings.ReplaceAll(templateFile, "/", "-")).Parse(string(b))
	if err != nil {
		logrus.Errorf("%v, template: %v", err, templateFile)
		return nil, err
	}

	if err = clusterTmpl.Execute(&tpl, templateData); err != nil {
		logrus.Errorf("err: %v rendering template error", err)
		return nil, err
	}

	return &tpl, nil
	//// dir for rke
	//if err := os.MkdirAll(rkeDir, os.ModePerm); err != nil {
	//	logrus.Errorf("err: %v, faild to create %v", err, rkeDir)
	//	panic(err)
	//}
	//
	//if err := ioutil.WriteFile(rkeDir+"/cluster.yml", tpl.Bytes(), 0644); err != nil {
	//	logrus.Errorf("err: %v, faild to cluster.yml %v", err, rkeDir)
	//	panic(err)
	//}
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

	// copy kubeconfig file for cnvrg user
	args = []string{"-lc", fmt.Sprintf(`su - cnvrg -c "cp %s/kube_config_cluster.yml %s/.kube/config"`, rkeDir, home)}
	cmd = exec.Command("/bin/bash", args...)
	if _, err := cmd.CombinedOutput(); err != nil {
		logrus.Error(err)
		panic(err)
	}

	//// copy kubeconfig file for root user (or current user)
	//args = []string{"-lc", fmt.Sprintf(`mkdir -p ~/.kube/config && cp %s/kube_config_cluster.yml ~/.kube/config"`, rkeDir)}
	//cmd = exec.Command("/bin/bash", args...)
	//if _, err := cmd.CombinedOutput(); err != nil {
	//	logrus.Error(err)
	//	panic(err)
	//}
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
			logrus.Infof("K8s is ready, time to deploy cnvrg %s! (run: cnvrgctl cnvrg up -h)", emoji.Sprint(":rocket:"))
			break
		}
		logrus.Infof("checking k8s ready status, attempt left: %d ...", 20-i)
		time.Sleep(10 * time.Second)
	}
}

func userHomeDir() string {
	return fmt.Sprintf("/home/%s", viper.GetString("cnvrg-user"))
}
