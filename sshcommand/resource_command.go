package sshcommand

import (
	"crypto/sha256"
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"golang.org/x/crypto/ssh"
)

const (
	// DefaultSSHPort represents default port used for SSH connections.
	DefaultSSHPort = 22

	// DefaultTimeout represents default timeout for long-standing operations
	// like connecting or retrying the execution.
	DefaultTimeout = "5m"

	// TTYSpeed defines virtual terminal default input and output speed.
	TTYSpeed = 14400
)

func resourceCommand() *schema.Resource {
	return &schema.Resource{
		Create: resourceCommandCreate,
		// Those 2 functions below does nothing, but must be implemented.
		Read:   resourceCommandRead,
		Delete: resourceCommandDelete,
		// Reuse create for updating.
		Update: resourceCommandCreate,
		Schema: resourceCommandSchema(),
	}
}

func resourceCommandSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"host": {
			Type:     schema.TypeString,
			Required: true,
		},
		"private_key": {
			Type:         schema.TypeString,
			Required:     true,
			Sensitive:    true,
			ValidateFunc: validatePrivateKeyFunc(),
		},
		"command": {
			Type:     schema.TypeString,
			Required: true,
		},
		"user": {
			Type:     schema.TypeString,
			Optional: true,
			Default:  "root",
		},
		"port": {
			Type:     schema.TypeInt,
			Optional: true,
			Default:  DefaultSSHPort,
		},
		"connection_timeout": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      DefaultTimeout,
			ValidateFunc: validateTimeoutFunc(),
		},
		"ignore_execute_errors": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"retry": {
			Type:     schema.TypeBool,
			Optional: true,
			Default:  false,
		},
		"retry_timeout": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      DefaultTimeout,
			ValidateFunc: validateTimeoutFunc(),
		},
		"retry_interval": {
			Type:         schema.TypeString,
			Optional:     true,
			Default:      "5s",
			ValidateFunc: validateTimeoutFunc(),
		},
		"result": {
			Type:     schema.TypeString,
			Computed: true,
		},
	}
}

// This function opens TCP connection, SSH connection, executes given command and returns it's output.
func executeSSH(sshConfig *ssh.ClientConfig, address string, command string) ([]byte, bool, error) {
	connection, err := ssh.Dial("tcp", address, sshConfig)
	if err != nil {
		return []byte{}, false, fmt.Errorf("opening SSH connection: %s", err)
	}

	session, err := connection.NewSession()
	if err != nil {
		return []byte{}, false, fmt.Errorf("creating SSH session: %s", err)
	}

	defer func() {
		if err := session.Close(); err != nil {
			log.Printf("%s: closing SSH session: %v", address, err)
		}
	}()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,        // Disable echoing.
		ssh.TTY_OP_ISPEED: TTYSpeed, // Input speed = 14.4kbaud.
		ssh.TTY_OP_OSPEED: TTYSpeed, // Output speed = 14.4kbaud.
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return []byte{}, false, fmt.Errorf("requesting pseudo terminal: %s", err)
	}

	output, err := session.Output(command)
	if err != nil {
		return []byte{}, true, fmt.Errorf("executing command: %v", err)
	}

	return output, false, nil
}

func resourceCommandCreate(d *schema.ResourceData, meta interface{}) error {
	host := d.Get("host").(string)
	command := d.Get("command").(string)
	ignoreExecuteErrors := d.Get("ignore_execute_errors").(bool)
	retry := d.Get("retry").(bool)

	signer, err := ssh.ParsePrivateKey([]byte(d.Get("private_key").(string)))
	if err != nil {
		return fmt.Errorf("Unable to parse private key: %v", err)
	}

	connectionTimeout, err := time.ParseDuration(d.Get("connection_timeout").(string))
	if err != nil {
		return fmt.Errorf("Unable to parse connection timeout: %v", err)
	}

	retryTimeout, err := time.ParseDuration(d.Get("retry_timeout").(string))
	if err != nil {
		return fmt.Errorf("Unable to parse connection timeout: %v", err)
	}

	retryInterval, err := time.ParseDuration(d.Get("retry_interval").(string))
	if err != nil {
		return fmt.Errorf("Unable to parse retry interval: %v", err)
	}

	sshConfig := &ssh.ClientConfig{
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		Timeout:         connectionTimeout,
		User:            d.Get("user").(string),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // nolint:gosec
	}

	address := fmt.Sprintf("%s:%d", host, d.Get("port").(int))

	var output []byte

	var executionError bool

	// If retry is enabled, try to run command until we timeout
	if retry {
		start := time.Now()
		// Try until we timeout
		for time.Since(start) < retryTimeout {
			output, executionError, err = executeSSH(sshConfig, address, command)
			// If command executed successfully, we can finish
			if err == nil {
				break
			}
			// Wait specified interval between attempts
			time.Sleep(retryInterval)
		}

		// If command returned error, check if we can tolerate it
		if err != nil && !(executionError && ignoreExecuteErrors) {
			return err
		}
	} else {
		output, executionError, err = executeSSH(sshConfig, address, command)
		if err != nil && !(executionError && ignoreExecuteErrors) {
			return err
		}
	}

	if err := d.Set("result", string(output)); err != nil {
		return fmt.Errorf("setting %q field: %w", "result", err)
	}

	d.SetId(sha256sum(fmt.Sprintf("%s-%s", host, command)))

	return nil
}

func validatePrivateKeyFunc() schema.SchemaValidateFunc {
	return func(v interface{}, k string) (we []string, errors []error) {
		if _, err := ssh.ParsePrivateKey([]byte(v.(string))); err != nil {
			errors = append(errors, fmt.Errorf("parsing private key: %v", err))
		}

		return
	}
}

func validateTimeoutFunc() schema.SchemaValidateFunc {
	return func(v interface{}, k string) (we []string, errors []error) {
		if _, err := time.ParseDuration(v.(string)); err != nil {
			errors = append(errors, fmt.Errorf("parsing duration: %v", err))
		}

		return
	}
}

func resourceCommandRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceCommandDelete(d *schema.ResourceData, m interface{}) error {
	d.SetId("")

	return nil
}

func sha256sum(data string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(data)))
}
