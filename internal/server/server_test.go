package server_test

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/ambardhesi/runnable/internal/server"
)

func startServer(port int) *server.Server {
	config := server.Config{
		Port:           port,
		LogDir:         "logs",
		CertFilePath:   "../../certs/svr-cert.pem",
		KeyFilePath:    "../../certs/svr-key.pem",
		CaCertFilePath: "../../certs/ca-cert.pem",
		TestMode:       true,
	}

	s, err := server.NewServer(config)
	if err != nil {
		log.Printf("Failed to create server %v\n", err)
		os.Exit(1)
	}

	go func() {
		s.Start()
	}()
	return s
}

func TestUnauthorized(t *testing.T) {
	s := startServer(8080)
	defer s.Stop()

	// Alice starts a job
	args := []string{
		"--ca", "../../certs/ca-cert.pem",
		"--cert", "../../certs/alice-cert.pem",
		"--key", "../../certs/alice-key.pem",
		"start", "echo", "hello", "world",
	}

	cmd := exec.Command("../../runnable-client", args...)

	output, _ := cmd.CombinedOutput()
	strOutput := string(output)

	// hacky way of extracting the Job ID from the response, but oh well, running out of time
	i := strings.Index(strOutput, "jobID")
	jobID := strOutput[i+8 : len(strOutput)-3]

	// Bob tries to get it
	args = []string{
		"--ca", "../../certs/ca-cert.pem",
		"--cert", "../../certs/bob-cert.pem",
		"--key", "../../certs/bob-key.pem",
		"get", jobID,
	}

	cmd = exec.Command("../../runnable-client", args...)

	output, _ = cmd.CombinedOutput()
	strOutput = string(output)

	// We expect an error
	if !strings.Contains(strOutput, "User is unauthorized") {
		t.Errorf("expected user is unauthorized error")
	}
}

func TestHappyPath(t *testing.T) {
	s := startServer(8081)
	defer s.Stop()

	// Alice starts a job
	args := []string{
		"--ca", "../../certs/ca-cert.pem",
		"--cert", "../../certs/alice-cert.pem",
		"--key", "../../certs/alice-key.pem",
		"start", "echo", "hello", "world",
	}

	cmd := exec.Command("../../runnable-client", args...)

	output, _ := cmd.CombinedOutput()
	strOutput := string(output)

	// hacky way of extracting the Job ID from the response, but oh well, running out of time
	i := strings.Index(strOutput, "jobID")
	jobID := strOutput[i+8 : len(strOutput)-3]

	// Alice tries to get it
	args = []string{
		"--ca", "../../certs/ca-cert.pem",
		"--cert", "../../certs/alice-cert.pem",
		"--key", "../../certs/alice-key.pem",
		"get", jobID,
	}

	cmd = exec.Command("../../runnable-client", args...)

	output, _ = cmd.CombinedOutput()
	strOutput = string(output)

	// We expect no errors and job should be completed
	if !strings.Contains(strOutput, "\"state\":\"Completed\"") {
		t.Errorf("expected job to be completed")
	}
}
