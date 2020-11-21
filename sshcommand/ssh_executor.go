package sshcommand

import (
	"fmt"
	"log"
	"time"

	"golang.org/x/crypto/ssh"
)

// createSSHSession opens new SSH connection and configures the session.
//
// It is up to the caller to close the session.
func createSSHSession(sshConfig *ssh.ClientConfig, address string) (*ssh.Session, error) {
	connection, err := ssh.Dial("tcp", address, sshConfig)
	if err != nil {
		return nil, fmt.Errorf("opening SSH connection: %w", err)
	}

	session, err := connection.NewSession()
	if err != nil {
		return nil, fmt.Errorf("creating SSH session: %w", err)
	}

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,        // Disable echoing.
		ssh.TTY_OP_ISPEED: TTYSpeed, // Input speed = 14.4kbaud.
		ssh.TTY_OP_OSPEED: TTYSpeed, // Output speed = 14.4kbaud.
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return nil, fmt.Errorf("requesting pseudo terminal: %w", err)
	}

	return session, nil
}

// This function opens TCP connection, SSH connection, executes given command and returns it's output.
func executeSSH(sshConfig *ssh.ClientConfig, address string, command string) ([]byte, bool, error) {
	session, err := createSSHSession(sshConfig, address)
	if err != nil {
		return []byte{}, false, fmt.Errorf("creating SSH session: %w", err)
	}

	defer func() {
		if err := session.Close(); err != nil {
			log.Printf("%s: closing SSH session: %v", address, err)
		}
	}()

	output, err := session.Output(command)
	if err != nil {
		return []byte{}, true, fmt.Errorf("executing command: %w", err)
	}

	return output, false, nil
}

type sshExecutor struct {
	host                string
	command             string
	ignoreExecuteErrors bool
	retry               bool
	authMethods         []ssh.AuthMethod
	timeout             time.Duration
	user                string
	retryTimeout        time.Duration
	retryInterval       time.Duration
	port                int
}

func (e sshExecutor) execute() ([]byte, error) {
	sshConfig := &ssh.ClientConfig{
		Auth:            e.authMethods,
		Timeout:         e.timeout,
		User:            e.user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // nolint:gosec
	}

	address := fmt.Sprintf("%s:%d", e.host, e.port)

	output, executionError, err := executeSSH(sshConfig, address, e.command)

	// If no error occurred or execution error occurred and we are told to ignore it, just return the result.
	if err == nil || (err != nil && executionError && e.ignoreExecuteErrors) {
		return output, nil
	}

	// If error occursed and we are told to not retry, return the error.
	if err != nil && !e.retry {
		return nil, fmt.Errorf("execution error: %w", err)
	}

	start := time.Now()

	// Try again until we timeout.
	for time.Since(start) < e.retryTimeout {
		output, executionError, err = executeSSH(sshConfig, address, e.command)
		// If command executed successfully, we can finish.
		if err == nil {
			break
		}
		// Wait specified interval between attempts.
		time.Sleep(e.retryInterval)
	}

	// If command returned error, check if we can tolerate it.
	if err != nil && !(executionError && e.ignoreExecuteErrors) {
		return nil, err
	}

	return output, nil
}
