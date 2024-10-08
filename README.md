# Nitro Logger

A simple utility that listens on a vsock port for lines of text and logs them to a file.
Supports file rotation and cleanup.

The typical use case is running on a host machine running an AWS Nitro instance which produces logs.
The software in the enclave should forward console logs to the specified vsock port, which would otherwise be
inaccessible due to the isolation of the enclave.

## Build

There are no binaries for this just yet, but compiling is trivial.
Requirements:

- Go 1.22 or later (although this is probably compatible with much earlier versions)
- Make

Just run ```make build``` and find a nitro-logger binary in the root of the project.

## Use

Example usage:

```shell
nohup ./nitro-logger -file /var/log/my-nitro-app.log -port 8090 &
```

To customize behavior, ```nitro-logger -h``` will show you the available options and their default values.
You can customize the max file size before rotation, the number of files to keep, their max age and whether they should
be gzipped.

Note: `vsock` is a Linux-specific feature, however to facilitate testing and development, the logger will fall back to
listening on a regular TCP socket on platforms in which vsock is unsupported.

Also note that it has no way of detecting other instances of itself running, so make sure you only run one.
To terminate it run:
```shell
ps aux | grep nitro-logger | grep -v grep | awk '{print $2}' | xargs -r sudo kill -9
```
Eventually, you might want to daemonize this process, so it boots with the OS.

## Future development

Potentially add support for forwarding to services such as DataDog our Cloudwatch.
