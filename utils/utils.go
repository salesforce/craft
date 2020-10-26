// Copyright (c) 2020, salesforce.com, inc.
// All rights reserved.
// SPDX-License-Identifier: BSD-3-Clause
// For full license text, see LICENSE.txt file in the repo root or https://opensource.org/licenses/BSD-3-Clause

package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

func FileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func CmdExec(cmdStr, dir string) {
	log.Debugf("$(%s) %s", dir, cmdStr)
	cmdList := strings.Split(cmdStr, " ")

	out := exec.Command(cmdList[0], cmdList[1:]...)
	out.Dir = dir
	// stdoutStderr, err := out.CombinedOutput()
	stdoutStderr, err := out.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	out.Stderr = out.Stdout
	done := make(chan struct{})
	scanner := bufio.NewScanner(stdoutStderr)
	go func() {
		for scanner.Scan() {
			output := scanner.Text()
			if strings.Contains(output, "msg") {
				slice := strings.SplitAfter(output, "msg=")
				output = slice[len(slice)-1]
				log.Infof(strings.Trim(output, "\""))
			} else {
				log.Infof(output)
			}
		}
		done <- struct{}{}
	}()
	err = out.Start()
	if err != nil {
		log.Fatal(err)
	}
	<-done
	err = out.Wait()
	if err != nil {
		log.Fatal(err)
	}
}

func MinCmdExec(cmdStr, dir string) {
	log.Debugf("$(%s): %s", dir, cmdStr)
	cmdList := strings.Split(cmdStr, " ")

	out := exec.Command(cmdList[0], cmdList[1:]...)
	out.Dir = dir
	out.Env = os.Environ()
	stdoutStderr, err := out.CombinedOutput()
	log.Debug("%s", stdoutStderr)
	if err != nil {
		MinCmdExec(cmdStr, dir)
		//log.Fatal(err)
	}
}

func EnvCmdExec(cmdStr, dir string, env []string) {
	log.Debugf("$(%s): %s", dir, cmdStr)
	log.Debugf("env %s", env)
	cmdList := strings.Split(cmdStr, " ")

	out := exec.Command(cmdList[0], cmdList[1:]...)
	out.Dir = dir
	out.Env = os.Environ()
	for _, e := range env {
		out.Env = append(out.Env, e)
	}
	stdoutStderr, err := out.CombinedOutput()
	log.Infof("%s", stdoutStderr)
	if err != nil {
		log.Fatal(err)
	}

}
func Exists(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Fatal(err)
	}
}

func CheckGoPath() {
	if os.ExpandEnv("$GOPATH") == "" {
		fmt.Println("$GOPATH is not set.")
		os.Exit(0)
	}
}
