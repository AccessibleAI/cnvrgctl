package cmd

import (
	"bytes"
	"fmt"
	"github.com/cnvrgctl/pkg"
	"github.com/kyokomi/emoji/v2"
	"github.com/manifoldco/promptui"
	"github.com/markbates/pkger"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io/ioutil"
	"strings"
	"text/template"
)

var ClusterParams = []Param{
	{Name: "single-node", Value: true, Usage: "create single node K8s cnvrg cluster"},
	{Name: "cnvrg-user", Value: "cnvrg", Usage: "user for managing cnvrg stack"},
	{Name: "install-docker", Value: true, Usage: "set to false to disable docker installation"},
	{Name: "download-tools", Value: true, Usage: "set to false to disable tools download (rke,kubectl,k9s)"},
	{Name: "host", Value: "", Usage: "destination host for K8s deployment"},
	{Name: "ssh-port", Value: 22, Usage: "ssh port"},
	{Name: "ssh-user", Value: "", Usage: "user for ssh connection"},
	{Name: "ssh-pass", Value: "", Usage: "password for ssh connection"},
	{Name: "ssh-key", Value: "", Usage: "private key for ssh connection"},
}

var ClusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "deploy single node cnvrg K8s cluster",
}

var ClusterRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "remove K8s cnvrg cluster",
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Warn("K8s cluster and cnvrg user will be removed, all data will be lost!")
		prompt := promptui.Prompt{Label: "Delete cnvrg K8s cluster", IsConfirm: true}
		result, err := prompt.Run()
		if err != nil {
			return
		}
		if result == "y" {
			cleanup()
		}
	},
}
var ClusterUpCmd = &cobra.Command{
	Use:   "up",
	Short: "bring up cnvrg single nodes k8s cluster",
	Run: func(cmd *cobra.Command, args []string) {
		cnvrgUser := viper.GetString("cnvrg-user")
		logrus.Infof("deploying k8s cluster")

		//prepare env
		prepareSetup()

		if viper.GetBool("download-tools") {
			logrus.Info("downloading tools")
			pkg.NewCmd("sudo cluster-setup.sh downloadTools").Exec()
		}

		//installing docker
		if viper.GetBool("install-docker") {
			logrus.Infof("installing docker")
			pkg.NewCmd(`sudo cluster-setup.sh installDocker`).Exec()
		}

		// create user for cnvrg
		logrus.Infof("creating %s user", cnvrgUser)
		pkg.NewCmd(`sudo cluster-setup.sh createUser`).Exec()

		// create workdirs
		logrus.Infof("creating workdirs")
		pkg.NewCmd(`sudo cluster-setup.sh workdirs`).Exec()

		// add user to sudo,docker groups
		logrus.Infof("adding %s user to groups", cnvrgUser)
		pkg.NewCmd(`sudo cluster-setup.sh addUserToGroups`).Exec()

		// generate ssh keys
		logrus.Infof("generating ssh keys")
		pkg.NewCmd(`sudo su ` + cnvrgUser + ` -c "cluster-setup.sh generateSSHKeys"`).Exec()

		prepareRkeSetup()
		rkeUp()

		//dumpDeploymentAssets()
		//generateRkeClusterManifest()
		//fixPermissions()
		//rkeUp()
		//checkClusterReady()
	},
}

func prepareSetup() {
	tmpDir := fmt.Sprintf("/tmp/%s", viper.GetString("ssh-user"))
	tmpClusterSetup := fmt.Sprintf("%s/cluster-setup.sh", tmpDir)
	pkg.NewCmd(fmt.Sprintf("mkdir -p %s", tmpDir)).Exec()

	if err := pkg.NewCmd("").Copy(generateClusterSetupScript(), tmpClusterSetup); err != nil {
		logrus.Error(err)
		panic(err)
	}

	// prepare deployment scripts
	logrus.Infof("copying deployment scripts")
	pkg.NewCmd(`chmod 0755 ` + tmpClusterSetup).Exec()
	cmd := pkg.NewCmd(fmt.Sprintf(`PASSWD=%s `+tmpClusterSetup+` patchSshUser`, viper.GetString("ssh-pass")))
	cmd.Hidden = true
	cmd.Exec()
	pkg.NewCmd(`sudo mv ` + tmpClusterSetup + " /usr/local/bin/cluster-setup.sh").Exec()

}

func generateRkeClusterManifest() string {
	var tpl bytes.Buffer
	templateData := map[string]interface{}{
		"Data": map[string]interface{}{
			"Server":        getMainIp(),
			"User":          viper.GetString("cnvrg-user"),
			"SshPrivateKey": fmt.Sprintf("/home/%s/.ssh/id_rsa", viper.GetString("cnvrg-user")),
			"ExternalIp":    viper.GetString("host"),
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

	return tpl.String()
}

func prepareRkeSetup() {
	cnvrgUser := viper.GetString("cnvrg-user")
	sshUser := viper.GetString("ssh-user")
	rkeClusterTmpFile := "/tmp/" + sshUser + "/cluster.yml"
	rkeClusterFinalFile := "/home/" + cnvrgUser + "/rke-cluster/cluster.yml"
	logrus.Infof("copying rke cluster.yml")
	pkg.NewCmd(fmt.Sprintf("sudo rm -fr %s", rkeClusterTmpFile)).Exec()
	if err := pkg.NewCmd("").Copy(generateRkeClusterManifest(), rkeClusterTmpFile); err != nil {
		logrus.Error(err)
		panic(err)
	}
	cmd := fmt.Sprintf(`sudo chown %s:%s %s && sudo mv %s %s`, cnvrgUser, cnvrgUser, rkeClusterTmpFile, rkeClusterTmpFile, rkeClusterFinalFile)
	pkg.NewCmd(cmd).Exec()
}

func rkeUp() {
	cnvrgUser := viper.GetString("cnvrg-user")
	logrus.Infof("rke up")
	rkeUpCmd := fmt.Sprintf(`sudo su %s -c "cd /home/%s/rke-cluster && rke -d up --ignore-docker-version && cp kube_config_cluster.yml ~/.kube/config"`, cnvrgUser, cnvrgUser)
	pkg.NewCmd(rkeUpCmd).Exec()
	logrus.Infof("K8s is ready, time to deploy cnvrg %s! (docs: https://github.com/accessibleAI/cnvrgio-operator)", emoji.Sprint(":rocket:"))
}

func generateClusterSetupScript() string {
	templateData := map[string]interface{}{
		"Data": map[string]interface{}{
			"CnvrgUser": viper.GetString("cnvrg-user"),
			"SshUser":   viper.GetString("ssh-user"),
		},
	}
	buffer, err := renderTemplate("/pkg/assets/cluster-setup.sh", templateData)
	if err != nil {
		logrus.Errorf("error generating cluster setup script err: %v", err)
		panic(err)
	}

	if viper.GetBool("dry-run") {
		logrus.Infof("\n%s", buffer.String())
		return ""
	}
	//if err := ioutil.WriteFile("/usr/local/bin/cluster-setup.sh", buffer.Bytes(), 0755); err != nil {
	//	logrus.Errorf("err: %v, faild to save cluster-setup.sh %v", err, rkeDir)
	//	panic(err)
	//}
	return buffer.String()

}

func getMainIp() string {
	// get node main ip
	logrus.Infof("getting node main ip")
	cmd := pkg.NewCmd("sudo cluster-setup.sh getMainIp")
	cmd.Exec()
	if len(cmd.Output) < 1 {
		logrus.Error("can't detect main IP")
		panic("can't detect main IP")
	}
	return strings.TrimSpace(cmd.Output[0])
}

func renderTemplate(templateFile string, templateData map[string]interface{}) (*bytes.Buffer, error) {
	var tpl bytes.Buffer
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
}

func cleanup() {
	logrus.Info("removing cnvrg K8s cluster and cnvrg user")
	cnvrgUser := viper.GetString("cnvrg-user")
	cmd := fmt.Sprintf(`if [ $(cat /etc/passwd | grep %s | wc -l) -eq 1 ]; then sudo su %s -c "cluster-setup.sh removeRke"; fi`, cnvrgUser, cnvrgUser)
	pkg.NewCmd(cmd).Exec()
	pkg.NewCmd("sudo cluster-setup.sh delUser").Exec()
	logrus.Info("cleanup successfully finished!")
}
