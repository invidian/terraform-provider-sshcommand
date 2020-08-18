# Command Data Source

This data executes given command on remote system on every Terraform run and provide it's output as an attribute.

## Example Usage

```hcl
data "sshcommand_command" "ssh_host_fingerprints" {
  host               = "example"
  command            = "ssh-keygen -r $(hostname -f) | cut -d' ' -f4-6"
  private_key        = file(".ssh/id_rsa")
}
```

## Argument Reference

* `host` - (Required) Hostname to connect.
* `private_key` - (Required) SSH private key used for authentication (SSH Agent support is not implemented).
* `command` - (Required) Command to execute.
* `user` - (Optional) User used for SSH log in. Defaults value is `root`.
* `port` - (Optional) Port to open SSH connection. Defaults is `22`.
* `connection_timeout` - (Optional) Timeout for opening TCP connection. This should be decreased when using `retry`. Defaults is `5m`.
* `retry` - (Optional) If this is set to true, plugin will retry to connect/execute command until `retry_timeout` is reached. Defaults to 'false'.
* `retry_timeout` - (Optional) Time after which retry logic should time out. Defaults to `5m`.
* `retry_interval` - (Optional) Specifies how long to wait between each attempt. Defaults to `5s`.
* `ignore_execute_errors` - (Optional) If true, resource will be created even if executed command returns non 0 exit code. Defaults to `false`.

## Attribute Reference

* `result` - Output of executed command.
