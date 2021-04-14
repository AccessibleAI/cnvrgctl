package pkg

import (
	"bufio"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
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
