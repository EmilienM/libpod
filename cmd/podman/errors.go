// +build !remoteclient

package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/sirupsen/logrus"
)

func outputError(err error) {
	if MainGlobalOpts.LogLevel == "debug" {
		logrus.Errorf(err.Error())
	} else {
		ee, ok := err.(*exec.ExitError)
		if ok {
			if status, ok := ee.Sys().(syscall.WaitStatus); ok {
				exitCode = status.ExitStatus()
			}
		}
		fmt.Fprintln(os.Stderr, "Error:", err.Error())
	}
}
