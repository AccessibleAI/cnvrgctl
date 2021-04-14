package pkg

import (
	"bufio"
	"fmt"
	"github.com/pkg/sftp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

func ExecBashCommand(command string, args []string) (string, error) {
	userCmd := exec.Command(command, args...)
	if out, err := userCmd.CombinedOutput(); err != nil {
		return string(out), err
	} else {
		return string(out), nil
	}
}

func ExecBashScript(script string) {
	args := append([]string{"-lc"}, script)
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
}

func ExecSshBashScript(command string) {
	if viper.GetBool("dry-run") {
		logrus.Info("dry-run enabled, skipping ExecSshCommand")
		return
	}
	config := sshClientConfig()
	sshHost := hostAddress()
	conn, err := ssh.Dial("tcp", sshHost, config)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	command = fmt.Sprintf("/bin/bash -lc '%s'", command)
	logrus.Infof("[%s] executing %s ", sshHost, command)
	runCommand(command, conn)
}

func SSHCopyFile(scriptStr, dstPath string) error {
	if viper.GetBool("dry-run") {
		logrus.Info("dry-run enabled, skipping SSHCopyFile")
		return nil
	}

	config := sshClientConfig()
	sshHost := hostAddress()
	logrus.Infof("[%s] copying %s ", sshHost, dstPath)
	conn, err := ssh.Dial("tcp", sshHost, config)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	// open an SFTP session over an existing ssh connection.
	sftp, err := sftp.NewClient(conn)
	if err != nil {
		return err
	}
	defer sftp.Close()
	// Create the destination file
	dstFile, err := sftp.Create(dstPath)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	// write to file
	if _, err := dstFile.ReadFrom(strings.NewReader(scriptStr)); err != nil {
		return err
	}
	return nil
}

func getAuthMethod() []ssh.AuthMethod {
	if viper.GetString("ssh-pass") != "" {
		return []ssh.AuthMethod{
			ssh.Password(viper.GetString("ssh-pass")),
		}
	}
	if viper.GetString("ssh-key") != "" {
		return []ssh.AuthMethod{
			publicKey(viper.GetString("ssh-key")),
		}
	}
	logrus.Error("both --ssh-pass and --ssh-key flags wasn't set, can't continue")
	panic("--ssh-pass and --ssh-key flags wasn't set, can't continue")

}

func runCommand(cmd string, conn *ssh.Client) {
	sess, err := conn.NewSession()
	if err != nil {
		panic(err)
	}
	defer sess.Close()

	sessStdOut, err := sess.StdoutPipe()
	if err != nil {
		panic(err)
	}
	stdoutScanner := bufio.NewScanner(sessStdOut)
	go func() {
		for stdoutScanner.Scan() {
			logrus.Infof("%s", stdoutScanner.Text())
		}
	}()

	sessStderr, err := sess.StderrPipe()
	if err != nil {
		panic(err)
	}
	stderrScanner := bufio.NewScanner(sessStderr)
	go func() {
		for stderrScanner.Scan() {
			logrus.Errorf("%s", stderrScanner.Text())
		}
	}()
	//
	//go func() {
	//	io.Copy(os.Stderr, sessStderr)
	//}()

	err = sess.Run(cmd)
	if err != nil {
		panic(err)
	}
}

func publicKey(path string) ssh.AuthMethod {
	key, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		panic(err)
	}
	return ssh.PublicKeys(signer)
}

func sshClientConfig() *ssh.ClientConfig {
	timeout, _ := time.ParseDuration("10s")
	return &ssh.ClientConfig{
		Config:          ssh.Config{},
		User:            viper.GetString("ssh-user"),
		Auth:            getAuthMethod(),
		Timeout:         timeout,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

}

func hostAddress() string {
	connAddress := fmt.Sprintf("%s:%d", viper.GetString("host"), viper.GetInt("ssh-port"))
	logrus.Infof("connecting to: %s", connAddress)
	return connAddress
}
