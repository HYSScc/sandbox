package main

import (
	"bytes"
	"flag"
	"os"
	"os/exec"
	"syscall"
	"time"
	"github.com/getsentry/raven-go"
)

func main() {
	compiler := flag.String("compiler", "g++", "CPP compiler with abs path")
	basedir := flag.String("basedir", "/tmp", "basedir of tmp CPP code snippet")
	filename := flag.String("filename", "Main.cpp", "name of file to be compiled")
	timeout := flag.Int("timeout", 10, "timeout in seconds")
	flag.Parse()

	var stdout, stderr bytes.Buffer
	cmd := exec.Command(*compiler, *filename, "-save-temps", "-std=gnu++11", "-fmax-errors=10", "-o", "Main")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = *basedir

	time.AfterFunc(time.Duration(*timeout)*time.Millisecond, func() {
		syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
	})
	err := cmd.Run()

	if err != nil {
		os.Stderr.WriteString(err.Error() + "\n")
		os.Stderr.WriteString(stderr.String())
		// compiler bomb hangs for <-timeout>
		if len(stderr.String()) == 0 {
			raven.SetDSN("http://e79ebf76a31a43d18ef7bdfa7381537e:5b21a25106584b39ac22ebf0752412db@127.0.0.1:12000/3")
			raven.CaptureMessageAndWait(*basedir, map[string]string{"error": "CompileTimeExceeded"})
		}
		return
	}

	os.Stdout.WriteString("Compile OK")
}

