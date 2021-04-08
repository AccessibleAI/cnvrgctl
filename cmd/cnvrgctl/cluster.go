package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os/exec"
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

func createUser() {

	argUser := []string{"-m", "-d", "/home/cnvrg", "-s", "/bin/sh", "-p", "paMfuNMgwFAX", "cnvrg"}
	userCmd := exec.Command("useradd", argUser...)


	if out, err := userCmd.CombinedOutput(); err != nil {
		fmt.Println(err, "There was an error by adding user", "cnvrg")
	} else {
		fmt.Printf("Output: %s\n", out)
	}
}
