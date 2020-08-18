# SSH Command Terraform Provider [![Build Status](https://travis-ci.com/invidian/terraform-provider-sshcommand.svg?branch=master)](https://travis-ci.com/invidian/terraform-provider-sshcommand)

This provider allow to execute commands remotely via SSH and capture the output from them.

This provider uses built-in [crypto/ssh](https://godoc.org/golang.org/x/crypto/ssh) Golang library to act as a SSH Client. Currently the implementation is very limited, but can be easily extended.

## Table of contents
* [Requirements](#requirements)
* [Building](#building)
* [Installing the provider](#installing-the-provider)
* [Resources](#resources)
* [Example usage](#example-usage)
* [Authors](#authors)

## Requirements

-	[Terraform](https://www.terraform.io/downloads.html) 0.11.x
-	[Go](https://golang.org/doc/install) 1.11.x

## Building

For testing builds, simply run `docker build .`, which will download all dependencies, run build, test and linter.

For local builds, simply follow the steps from [Dockerfile](https://github.com/invidian/terraform-provider-sshcommand/blob/master/Dockerfile).

## Installing the provider

After building the provider, install it using the Terraform instructions for [installing a third party provider](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins).

## Resources

### sshcommand_command

This resource executes given command on remote system and stores it's output in Terraform.

#### Parameters
  - `host` - Hostname to connect.
  - `private_key` - SSH private key used for authentication (SSH Agent support is not implemented).
  - `command` - Command to execute.
  - `user` - User used for SSH log in. Defaults value is `root`.
  - `port` - Port to open SSH connection. Defaults is `22`.
  - `connection_timeout` - Timeout for opening TCP connection. This should be decreased when using `retry`. Defaults is `5m`.
  - `retry` - If this is set to true, plugin will retry to connect/execute command until `retry_timeout` is reached. Defaults to 'false'.
  - `retry_timeout` - Time after which retry logic should time out. Defaults to `5m`.
  - `retry_interval` - Specifies how long to wait between each attempt. Defaults to `5s`.
  - `ignore_execute_errors` - If true, resource will be created even if executed command returns non 0 exit code. Defaults to `false`.

#### Attributes
  - `result` - Output of executed command.

## Example usage
```hcl
provider "sshcommand" {
  version = "~> 0.1.0"
}

output "example" {
  value = "\n${sshcommand_command.ssh_host_fingerprints.result}"
}

resource "sshcommand_command" "ssh_host_fingerprints" {
  host               = "example"
  command            = "ssh-keygen -r $(hostname -f) | cut -d' ' -f4-6"
  private_key        = "${file(".ssh/id_rsa")}"
}

# Reboot server after OS installation
resource "sshcommand_command" "reboot" {
  host                  = "${var.node_ip}"
  command               = "reboot"
  private_key           = "${var.ssh_private_key}"
  ignore_execute_errors = true
  depends_on            = [ "null_resource.os_install" ]
}

# Make sure you SSH into correct system
resource "sshcommand_command" "wait_for_os" {
  host           = "${var.node_ip}"
  command        = "grep ID=flatcar /etc/os-release"
  private_key    = "${var.ssh_private_key}"
  # If grep fails or SSH connection gets refused, resource will be trying again.
  retry          = true
  retry_interval = "1s"
}
```

## Authors
* **Mateusz Gozdek** - *Initial work* - [invidian](https://github.com/invidian)
