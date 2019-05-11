package main

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"golang.org/x/crypto/ssh"
)

func resourceCommand() *schema.Resource {
	return &schema.Resource{
		Create: resourceCommandCreate,
		// Those 2 functions below does nothing, but must be implemented.
		Read:   resourceCommandRead,
		Delete: resourceCommandDelete,
		// Reuse create for updating
		Update: resourceCommandCreate,

		Schema: map[string]*schema.Schema{
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
				Default:  22,
			},
			"connection_timeout": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "5m",
				ValidateFunc: validateTimeoutFunc(),
			},
			"ignore_execute_errors": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"result": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceCommandCreate(d *schema.ResourceData, meta interface{}) error {
	host := d.Get("host").(string)
	command := d.Get("command").(string)
	ignore_execute_errors := d.Get("ignore_execute_errors").(bool)
	signer, err := ssh.ParsePrivateKey([]byte(d.Get("private_key").(string)))
	if err != nil {
		return fmt.Errorf("Unable to parse private key: %v", err)
	}

	connection_timeout, err := time.ParseDuration(d.Get("connection_timeout").(string))
	if err != nil {
		return fmt.Errorf("Unable to parse connection timeout: %v", err)
	}

	sshConfig := &ssh.ClientConfig{
		Auth: []ssh.AuthMethod{
			ssh.PublicKeys(signer),
		},
		Timeout:         connection_timeout,
		User:            d.Get("user").(string),
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	connection, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, d.Get("port").(int)), sshConfig)
	if err != nil {
		return fmt.Errorf("Failed to open SSH connection: %s", err)
	}

	session, err := connection.NewSession()
	if err != nil {
		return fmt.Errorf("Failed to create session: %s", err)
	}
	defer session.Close()

	modes := ssh.TerminalModes{
		ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", 80, 40, modes); err != nil {
		return fmt.Errorf("request for pseudo terminal failed: %s", err)
	}

	output, err := session.Output(command)
	if err != nil && !ignore_execute_errors {
		return fmt.Errorf("Command execution failed: %v", err)
	}

	d.Set("result", string(output)) //nolint:errcheck

	d.SetId(sha256sum(fmt.Sprintf("%s-%s", host, command)))

	return nil
}

func validatePrivateKeyFunc() schema.SchemaValidateFunc {
	return func(v interface{}, k string) (we []string, errors []error) {
		_, err := ssh.ParsePrivateKey([]byte(v.(string)))
		if err != nil {
			errors = append(errors, fmt.Errorf("Unable to parse private key: %v", err))
		}
		return
	}
}

func validateTimeoutFunc() schema.SchemaValidateFunc {
	return func(v interface{}, k string) (we []string, errors []error) {
		_, err := time.ParseDuration(v.(string))
		if err != nil {
			errors = append(errors, fmt.Errorf("Unable to parse connection timeout: %v", err))
		}
		return
	}
}

func resourceCommandRead(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceCommandDelete(d *schema.ResourceData, m interface{}) error {
	return nil
}

func sha256sum(data string) string {
	return fmt.Sprintf("%x", sha256.Sum256([]byte(data)))
}
