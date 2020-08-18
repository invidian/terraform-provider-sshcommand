# SSH Command Provider

The SSH Command provider allows to execute commands remotely via SSH and capture the output from them.

This provider uses built-in [crypto/ssh](https://godoc.org/golang.org/x/crypto/ssh) Golang library to act as a SSH Client. Currently the implementation is very limited, but can be easily extended.

## Example Usage

```hcl
terraform {
  required_providers {
    sshcommand = {
      source  = "invidian/sshcommand"
      version = "0.2.0"
    }
  }
}

resource "sshcommand_command" "ssh_host_fingerprints" {
  host               = "example"
  command            = "ssh-keygen -r $(hostname -f) | cut -d' ' -f4-6"
  private_key        = file(".ssh/id_rsa")
}

output "example" {
  value = "\n${sshcommand_command.ssh_host_fingerprints.result}"
}

```

## Argument Reference

This provider currently takes no arguments.
