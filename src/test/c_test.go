package gotest

import (
	"testing"
	"os"
	"os/exec"
	"bytes"
	"strings"
)

// HELPER
// copy test source file `*.c` to tmp dir
func copyCSourceFile(name string, t *testing.T) (string, string) {
	t.Logf("Copying file %s ...", name)

	absPath, _ := os.Getwd()
	baseDir, projectDir := absPath+"/tmp", absPath+"/../.."
	os.MkdirAll(baseDir, os.ModePerm)

	cpCmd := exec.Command("cp", projectDir+"/src/test/resources/c/"+name, baseDir+"/Main.c")
	cpErr := cpCmd.Run()

	if cpErr != nil {
		os.RemoveAll(baseDir + "/")
		t.Error(cpErr.Error())
		t.FailNow()
	}

	t.Log("Done")
	return baseDir, projectDir
}

// HELPER
// compile C source file
func compileC(name, baseDir, projectDir string, t *testing.T) (string) {
	t.Logf("Compiling file %s ...", name)

	var compilerStderr bytes.Buffer
	compilerCmd := exec.Command(projectDir+"/bin/c_compiler", "-basedir="+baseDir)
	compilerCmd.Stderr = &compilerStderr
	compilerErr := compilerCmd.Run()

	if compilerErr != nil {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerErr.Error())
		t.FailNow()
	}

	t.Log("Done")
	return compilerStderr.String()
}

// HELPER
// run C binary in our container
func runC(baseDir, projectDir, memory, timeout string, t *testing.T) (string) {
	t.Log("Running binary /Main ...")

	var containerStdout, containerStderr bytes.Buffer
	containerArgs := []string{"-basedir=" + baseDir, "-input=10:10:23PM", "-expected=22:10:23", "-memory=" + memory, "-timeout=" + timeout}
	containerCmd := exec.Command(projectDir+"/bin/c_container", containerArgs...)
	containerCmd.Stdout = &containerStdout
	containerCmd.Stderr = &containerStderr

	if err := containerCmd.Run(); err != nil {
		os.RemoveAll(baseDir + "/")
		t.Error(err.Error())
		t.FailNow()
	}

	t.Log(containerStderr.String())
	return containerStdout.String()
}

func Test_C_AC(t *testing.T) {
	name := "ac.c"
	baseDir, projectDir := copyCSourceFile(name, t)
	compilerStderr := compileC(name, baseDir, projectDir, t)

	if len(compilerStderr) > 0 {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr)
		t.FailNow()
	}

	containerErr := runC(baseDir, projectDir, "64", "1000", t)
	if !strings.Contains(containerErr, "\"status\":0") {
		os.RemoveAll(baseDir + "/")
		t.Error(containerErr + " => status != 0")
		t.FailNow()
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Compiler_Bomb_0(t *testing.T) {
	name := "compiler_bomb_0.c"
	baseDir, projectDir := copyCSourceFile(name, t)
	compilerStderr := compileC(name, baseDir, projectDir, t)

	if !strings.Contains(compilerStderr, "signal: killed") {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr + " => Compile error does not contain string `signal: killed`")
		t.FailNow()
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Compiler_Bomb_1(t *testing.T) {
	name := "compiler_bomb_1.c"
	baseDir, projectDir := copyCSourceFile(name, t)
	compilerStderr := compileC(name, baseDir, projectDir, t)

	if !strings.Contains(compilerStderr, "signal: killed") {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr + " => Compile error does not contain string `signal: killed`")
		t.FailNow()
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Compiler_Bomb_2(t *testing.T) {
	name := "compiler_bomb_2.c"
	baseDir, projectDir := copyCSourceFile(name, t)
	compilerStderr := compileC(name, baseDir, projectDir, t)

	if !strings.Contains(compilerStderr, "signal: killed") {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr + " => Compile error does not contain string `signal: killed`")
		t.FailNow()
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Fork_Bomb(t *testing.T) {
	name := "fork_bomb.c"
	baseDir, projectDir := copyCSourceFile(name, t)
	compilerStderr := compileC(name, baseDir, projectDir, t)

	if len(compilerStderr) > 0 {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr)
		t.FailNow()
	}

	containerErr := runC(baseDir, projectDir, "64", "1000", t)

	if !strings.Contains(containerErr, "Runtime Error") {
		os.RemoveAll(baseDir + "/")
		t.Error(containerErr)
		t.FailNow()
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Get_Host_By_Name(t *testing.T) {
	name := "get_host_by_name.c"
	baseDir, projectDir := copyCSourceFile(name, t)
	compilerStderr := compileC(name, baseDir, projectDir, t)

	if len(compilerStderr) > 0 {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr)
		t.FailNow()
	}

	containerErr := runC(baseDir, projectDir, "64", "1000", t)

	// Main.c:(.text+0x28): warning: Using 'gethostbyname' in statically linked applications
	// requires at runtime the shared libraries from the glibc version used for linking
	if !strings.Contains(containerErr, "\"status\":2") {
		os.RemoveAll(baseDir + "/")
		t.Error(containerErr)
		t.FailNow()
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Include_Leaks(t *testing.T) {
	name := "include_leaks.c"
	baseDir, projectDir := copyCSourceFile(name, t)
	compilerStderr := compileC(name, baseDir, projectDir, t)

	if !strings.Contains(compilerStderr, "/etc/shadow") {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr + " => Compile error does not contain string `/etc/shadow`")
		t.FailNow()
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Infinite_Loop(t *testing.T) {
	name := "infinite_loop.c"
	baseDir, projectDir := copyCSourceFile(name, t)
	compilerStderr := compileC(name, baseDir, projectDir, t)

	if len(compilerStderr) > 0 {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr)
		t.FailNow()
	}

	containerErr := runC(baseDir, projectDir, "64", "1000", t)

	if !strings.Contains(containerErr, "Runtime Error") {
		os.RemoveAll(baseDir + "/")
		t.Error(containerErr)
		t.FailNow()
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Memory_Allocation(t *testing.T) {
	name := "memory_allocation.c"
	baseDir, projectDir := copyCSourceFile(name, t)
	compilerStderr := compileC(name, baseDir, projectDir, t)

	if len(compilerStderr) > 0 {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr)
		t.FailNow()
	}

	containerErr := runC(baseDir, projectDir, "8", "5000", t)

	// `Killed` is sent to tty by kernel (and record will also be kept in /var/log/message)
	// both stdout and stderr are empty which will lead to status WA
	// OR
	// just running out of time
	if !strings.ContainsAny(containerErr, "\"status\":5 & Runtime Error") {
		os.RemoveAll(baseDir + "/")
		t.Error(containerErr)
		t.FailNow()
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Plain_Text(t *testing.T) {
	name := "plain_text.c"
	baseDir, projectDir := copyCSourceFile(name, t)

	compilerStderr := compileC(name, baseDir, projectDir, t)

	if !strings.Contains(compilerStderr, "error") {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr + " => Compile error does not contain string `error`")
		t.FailNow()
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Run_Command_Line_0(t *testing.T) {
	name := "run_command_line_0.c"
	baseDir, projectDir := copyCSourceFile(name, t)
	compilerStderr := compileC(name, baseDir, projectDir, t)

	if len(compilerStderr) > 0 {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr)
		t.FailNow()
	}

	containerErr := runC(baseDir, projectDir, "16", "1000", t)

	if !strings.Contains(containerErr, "\"status\":5") {
		os.RemoveAll(baseDir + "/")
		t.Error(containerErr)
		t.FailNow()
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Run_Command_Line_1(t *testing.T) {
	name := "run_command_line_1.c"
	baseDir, projectDir := copyCSourceFile(name, t)
	compilerStderr := compileC(name, baseDir, projectDir, t)

	if len(compilerStderr) > 0 {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr)
		t.FailNow()
	}

	containerErr := runC(baseDir, projectDir, "16", "1000", t)

	if !strings.Contains(containerErr, "\"status\":5") {
		os.RemoveAll(baseDir + "/")
		t.Error(containerErr)
		t.FailNow()
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_Syscall_0(t *testing.T) {
	name := "syscall_0.c"
	baseDir, projectDir := copyCSourceFile(name, t)
	compilerStderr := compileC(name, baseDir, projectDir, t)

	if len(compilerStderr) > 0 {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr)
		t.FailNow()
	}

	containerErr := runC(baseDir, projectDir, "16", "1000", t)

	if !strings.Contains(containerErr, "\"status\":5") {
		os.RemoveAll(baseDir + "/")
		t.Error(containerErr)
		t.FailNow()
	}

	os.RemoveAll(baseDir + "/")
}

func Test_C_TCP_Client(t *testing.T) {
	name := "tcp_client.c"
	baseDir, projectDir := copyCSourceFile(name, t)
	compilerStderr := compileC(name, baseDir, projectDir, t)

	if len(compilerStderr) > 0 {
		os.RemoveAll(baseDir + "/")
		t.Error(compilerStderr)
		t.FailNow()
	}

	containerErr := runC(baseDir, projectDir, "16", "5000", t)

	if !strings.Contains(containerErr, "Runtime Error") {
		os.RemoveAll(baseDir + "/")
		t.Error(containerErr)
		t.FailNow()
	}

	os.RemoveAll(baseDir + "/")
}
