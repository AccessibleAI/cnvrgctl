package main

import (
	"crypto/rand"
	"encoding/base64"
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

func createRandom(n int) string {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Println(err)
		//os.Exit(1)
	}
	return string(b)
}

func createUser() {
	encrypt := base64.StdEncoding.EncodeToString([]byte(createRandom(9)))
	argPass := []string{"-c", fmt.Sprintf("echo %s:%s | chpasswd", "cnvrg", encrypt)}
	passCmd := exec.Command("/bin/sh", argPass...)
	if out, err := passCmd.Output(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Output: %s\n", out)
	}

}
