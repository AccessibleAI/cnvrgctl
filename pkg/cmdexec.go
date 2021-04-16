package pkg

import (
	"bufio"
	"fmt"
	"github.com/pkg/sftp"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"strings"
	"time"
)

type SshBashCommand struct {
	SshHost string
	SshPort int
	SshUser string
	SshPass string
	SshKey  string
	Command string
	Output  []string
	Hidden  bool

}

func NewCmd(command string) *SshBashCommand {
	c := SshBashCommand{
		SshHost: viper.GetString("host"),
		SshPort: viper.GetInt("port"),
		SshUser: viper.GetString("ssh-user"),
		SshPass: viper.GetString("ssh-pass"),
		SshKey:  viper.GetString("ssh-key"),
		Command: fmt.Sprintf("/bin/bash -lc '%s'", command),
		Hidden:  false,
	}
	return &c
}

func (c *SshBashCommand) Exec() {
	if viper.GetBool("dry-run") {
		logrus.Info("dry-run enabled, skipping ExecSshCommand")
		return
	}
	config := c.sshClientConfig()
	sshHost := c.hostAddress()
	conn, err := ssh.Dial("tcp", sshHost, config)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	if c.Hidden == false {
		logrus.Infof("[%s] executing %s ", c.SshHost, c.Command)
	}
	c.Output = runCommand(c.Command, conn)
}

func (c *SshBashCommand) Copy(scriptStr, dstPath string) error {
	if viper.GetBool("dry-run") {
		logrus.Info("dry-run enabled, skipping SSHCopyFile")
		return nil
	}

	config := c.sshClientConfig()
	sshHost := c.hostAddress()
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

func (c *SshBashCommand) getAuthMethod() []ssh.AuthMethod {
	if c.SshPass != "" {
		return []ssh.AuthMethod{
			ssh.Password(c.SshPass),
		}
	}
	if c.SshKey != "" {
		return []ssh.AuthMethod{
			publicKey(c.SshKey),
		}
	}
	logrus.Error("both --ssh-pass and --ssh-key flags wasn't set, can't continue")
	panic("--ssh-pass and --ssh-key flags wasn't set, can't continue")

}

func (c *SshBashCommand) sshClientConfig() *ssh.ClientConfig {
	timeout, _ := time.ParseDuration("10s")
	return &ssh.ClientConfig{
		Config:          ssh.Config{},
		User:            viper.GetString("ssh-user"),
		Auth:            c.getAuthMethod(),
		Timeout:         timeout,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
}

func (c *SshBashCommand) hostAddress() string {
	connAddress := fmt.Sprintf("%s:%d", viper.GetString("host"), viper.GetInt("ssh-port"))
	logrus.Infof("connecting to: %s", connAddress)
	return connAddress
}

func runCommand(cmd string, conn *ssh.Client) []string {
	var output []string
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
			stdout := stdoutScanner.Text()
			output = append(output, stdout)
			logrus.Infof("%s", stdout)
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

	err = sess.Run(cmd)
	if err != nil {
		panic(err)
	}
	return output
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
