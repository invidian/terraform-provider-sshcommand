package sshcommand

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
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

//nolint:funlen
func resourceCommandSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"host": {
			Type:     schema.TypeString,
			Required: true,
		},
		"private_key": {
			Type:         schema.TypeString,
			Optional:     true,
			Sensitive:    true,
			ValidateFunc: validatePrivateKeyFunc(),
		},
		"password": {
			Type:      schema.TypeString,
			Optional:  true,
			Sensitive: true,
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

func resourceCommandToSSHExecutor(d *schema.ResourceData) (*sshExecutor, error) {
	connectionTimeout, _ := time.ParseDuration(d.Get("connection_timeout").(string))
	retryTimeout, _ := time.ParseDuration(d.Get("retry_timeout").(string))
	retryInterval, _ := time.ParseDuration(d.Get("retry_interval").(string))

	authMethods := []ssh.AuthMethod{}

	if pk := d.Get("private_key").(string); pk != "" {
		signer, _ := ssh.ParsePrivateKey([]byte(d.Get("private_key").(string)))

		authMethods = append(authMethods, ssh.PublicKeys(signer))
	}

	if p := d.Get("password").(string); p != "" {
		authMethods = append(authMethods, ssh.Password(p))
	}

	if len(authMethods) == 0 {
		return nil, fmt.Errorf("no auth methods specified")
	}

	return &sshExecutor{
		host:                d.Get("host").(string),
		command:             d.Get("command").(string),
		ignoreExecuteErrors: d.Get("ignore_execute_errors").(bool),
		retry:               d.Get("retry").(bool),
		authMethods:         authMethods,
		timeout:             connectionTimeout,
		user:                d.Get("user").(string),
		retryTimeout:        retryTimeout,
		retryInterval:       retryInterval,
		port:                d.Get("port").(int),
	}, nil
}

func resourceCommandCreate(d *schema.ResourceData, meta interface{}) error {
	e, err := resourceCommandToSSHExecutor(d)
	if err != nil {
		return fmt.Errorf("initializing: %w", err)
	}

	// Execute the command.
	output, err := e.execute()
	if err != nil {
		return fmt.Errorf("execution: %w", err)
	}

	// Save result on success.
	if err := d.Set("result", string(output)); err != nil {
		return fmt.Errorf("setting %q field: %w", "result", err)
	}

	d.SetId(sha256sum(fmt.Sprintf("%s-%s", e.host, e.command)))

	return nil
}

func validatePrivateKeyFunc() schema.SchemaValidateFunc {
	return func(v interface{}, k string) (we []string, errors []error) {
		if _, err := ssh.ParsePrivateKey([]byte(v.(string))); err != nil {
			errors = append(errors, fmt.Errorf("parsing private key: %w", err))
		}

		return
	}
}

func validateTimeoutFunc() schema.SchemaValidateFunc {
	return func(v interface{}, k string) (we []string, errors []error) {
		if _, err := time.ParseDuration(v.(string)); err != nil {
			errors = append(errors, fmt.Errorf("parsing duration: %w", err))
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
