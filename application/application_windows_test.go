// this is a copy of functions and logic from os/signal/signal_windows_test.go

//go:build windows
// +build windows

package application

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
	"testing"
	"time"
)

func sendCtrlBreak(t *testing.T, pid int) {
	d, e := syscall.LoadDLL("kernel32.dll")
	if e != nil {
		t.Fatalf("LoadDLL: %v\n", e)
	}
	p, e := d.FindProc("GenerateConsoleCtrlEvent")
	if e != nil {
		t.Fatalf("FindProc: %v\n", e)
	}
	r, _, e := p.Call(syscall.CTRL_BREAK_EVENT, uintptr(pid))
	if r == 0 {
		t.Fatalf("GenerateConsoleCtrlEvent: %v\n", e)
	}
}

func TestSignalStop(t *testing.T) {
	// create source file
	const source = `
package main

import (
    "log"
	"merqio/libs/application"
)

func main() {
	app := application.New()

	var deferred bool
	app.Defer(func() { deferred = true })

	var stop bool
	app.OnStop(func() { stop = true })

	app.Run()

	if deferred == false {
		log.Fatal("Failed to invoke deferred functions.")
	}

	if stop == false {
		log.Fatal("Failed to invoke stop functions.")
	}
}
`
	tmp, err := ioutil.TempDir("", "TestCtrlBreak")
	if err != nil {
		t.Fatal("TempDir failed: ", err)
	}
	defer os.RemoveAll(tmp)

	// write ctrlbreak.go
	name := filepath.Join(tmp, "ctlbreak")
	src := name + ".go"
	f, err := os.Create(src)
	if err != nil {
		t.Fatalf("Failed to create %v: %v", src, err)
	}
	defer f.Close()
	f.Write([]byte(source))

	// compile it
	exe := name + ".exe"
	defer os.Remove(exe)
	o, err := exec.Command("go", "build", "-o", exe, src).CombinedOutput()
	if err != nil {
		t.Fatalf("Failed to compile: %v\n%v", err, string(o))
	}

	// run it
	cmd := exec.Command(exe)
	var b bytes.Buffer
	cmd.Stdout = &b
	cmd.Stderr = &b
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: syscall.CREATE_NEW_PROCESS_GROUP,
	}
	err = cmd.Start()
	if err != nil {
		t.Fatalf("Start failed: %v", err)
	}
	go func() {
		time.Sleep(1 * time.Second)
		sendCtrlBreak(t, cmd.Process.Pid)
	}()
	err = cmd.Wait()
	if err != nil {
		t.Fatalf("Program exited with error: %v\n%v", err, string(b.Bytes()))
	}
}
