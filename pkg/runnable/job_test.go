package runnable_test

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"testing"
	"time"

	"github.com/ambardhesi/runnable/pkg/runnable"
	"github.com/ambardhesi/runnable/pkg/runnable/mock"
)

func TestNewJob(t *testing.T) {
	job, _ := runnable.NewJob("ownerID", "command name", "arg1", "arg2")
	b := &bytes.Buffer{}
	job.SetLogWriter(&mock.BufWriteCloser{b})

	if job.ID == "" {
		t.Errorf("expected non empty string for job ID")
	}

	if job.Status().ExitCode != -1 {
		t.Errorf("expected default exit code of -1")
	}

	if job.OwnerID != "ownerID" {
		t.Errorf("expected ownerID to be ownerID")
	}

	if job.Status().State != runnable.NotStarted {
		t.Errorf("expected initial state to be NotStarted")
	}
}

func TestStartEcho(t *testing.T) {
	cmd := "echo"
	args := []string{"hello world"}
	job, _ := runnable.NewJob("ownerID", cmd, args...)
	b := &bytes.Buffer{}
	job.SetLogWriter(&mock.BufWriteCloser{b})

	job.Cmd = fakeCmd(cmd, args...)
	job.Start()

	time.Sleep(500 * time.Millisecond)

	if ec := job.Status().ExitCode; ec != 0 {
		t.Errorf("expected exit code %v, got %v", 0, ec)
	}

	if state := job.Status().State; state != runnable.Completed {
		t.Errorf("expected state %v, got %v", runnable.Completed, state)
	}
}

func TestStartExit(t *testing.T) {
	cmd := "exit"
	// cmd is to exit with code 1
	args := []string{"1"}
	job, _ := runnable.NewJob("ownerID", cmd, args...)
	b := &bytes.Buffer{}
	job.SetLogWriter(&mock.BufWriteCloser{b})

	job.Cmd = fakeCmd(cmd, args...)
	job.Start()

	time.Sleep(500 * time.Millisecond)

	if ec := job.Status().ExitCode; ec != 1 {
		t.Errorf("expected exit code %v, got %v", 1, ec)
	}

	if state := job.Status().State; state != runnable.Completed {
		t.Errorf("expected state %v, got %v", runnable.Completed, state)
	}
}

func TestStopRunningJob(t *testing.T) {
	job, _ := runnable.NewJob("ownerID", "sleep", "1")
	b := &bytes.Buffer{}
	job.SetLogWriter(&mock.BufWriteCloser{b})

	job.Cmd = fakeCmd("sleep", "1")
	job.Start()

	job.Stop()
	// give job a chance to be stopped (or completed) before doing any assertions
	time.Sleep(500 * time.Millisecond)

	if ec := job.Status().ExitCode; ec != 0 {
		t.Errorf("expected exit code %v, got %v", 0, ec)
	}

	if state := job.Status().State; state != runnable.Stopped {
		t.Errorf("expected state %v, got %v", runnable.Stopped, state)
	}
}

func TestStopCompletedJob(t *testing.T) {
	job, _ := runnable.NewJob("ownerID", "exit", "0")
	b := &bytes.Buffer{}
	job.SetLogWriter(&mock.BufWriteCloser{b})

	job.Cmd = fakeCmd("exit", "0")
	job.Start()

	// give job a chance to complete before attempting to stop
	time.Sleep(200 * time.Millisecond)
	job.Stop()

	if ec := job.Status().ExitCode; ec != 0 {
		t.Errorf("expected exit code %v, got %v", 0, ec)
	}

	if state := job.Status().State; state != runnable.Completed {
		t.Errorf("expected state %v, got %v", runnable.Completed, state)
	}
}

func TestHelperProcess(t *testing.T) {
	if os.Getenv("GO_WANT_HELPER_PROCESS") != "1" {
		return
	}
	defer os.Exit(0)

	args := os.Args
	for len(args) > 0 {
		if args[0] == "--" {
			args = args[1:]
			break
		}
		args = args[1:]
	}
	if len(args) == 0 {
		os.Exit(2)
	}

	cmd, args := args[0], args[1:]
	switch cmd {
	case "echo":
		iargs := []interface{}{}
		for _, s := range args {
			iargs = append(iargs, s)
		}
	case "exit":
		n, _ := strconv.Atoi(args[0])
		os.Exit(n)
	case "sleep":
		s, _ := strconv.Atoi(args[0])
		time.Sleep(time.Duration(s) * time.Second)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command %q\n", cmd)
		os.Exit(2)
	}
}

// ref : https://npf.io/2015/06/testing-exec-command/
// This is also how the exec package does its tests.
// ref : https://github.com/golang/go/blob/master/src/os/exec/exec_test.go#L71
// Creates a fake cmd, which will call TestHelperProcess with the our cmd and its args.
// Inside TestHelperProcess, we can mock out the functionality.
func fakeCmd(command string, args ...string) *exec.Cmd {
	cs := []string{"-test.run=TestHelperProcess", "--", command}
	cs = append(cs, args...)
	cmd := exec.Command(os.Args[0], cs...)
	cmd.Env = []string{"GO_WANT_HELPER_PROCESS=1"}
	return cmd
}
