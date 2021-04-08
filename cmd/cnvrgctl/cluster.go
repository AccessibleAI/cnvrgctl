package main

import (
	"bufio"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
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
		createUser()
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

func createUser() {
	if !isUserExists("cnvrg") {
		argUser := []string{
			"-m",
			"-d",
			"/home/cnvrg",
			"-s",
			"/bin/bash",
			"-p",
			"paMfuNMgwFAX2",
			"--groups",
			"docker",
			"cnvrg"}
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
