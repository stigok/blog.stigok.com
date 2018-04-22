---
layout: post
title:  "Self healing reverse SSH setup with systemd"
date:   2018-04-22 19:48:07 +0200
categories: ssh systemd
---

I use a reverse SSH setup for my remotely deployed boxes in order to remotely manage them.
What is important to me with the setup is that

1. A persistent connection to the target host is established
2. The connection is securely encrypted
3. The client connects to a known target port for easy discovery
4. If port forwarding fails, the connection closes
5. It let's me connect to the clients via SSH

To achieve that, let's take a look at my `ssh` command arguments with excerpts from the manpage

```
-g  Allows remote hosts to connect to local forwarded ports
-N  Do not execute a remote command
-T  Disable pseudo-terminal allocation
-o  Used to give options in the format used in the configuration file (man ssh_config)
  ServerAliveInterval   Interval in seconds to ping the server while connection has been inactive
  ExitOnForwardFailure  Whether to terminate the connection if it cannot set up all requested port forwards
-R  Forward given remote TCP port (22221) to the local port (22)
-v  Verbose mode. More v's increase verbosity.
```

I call it self heailing as the `ExitOnForwardFailure` option is enabled, and will make sure that
the TCP port forward is successfully established. This option, combined with running the `ssh`
process with systemd, will establish the connection on boot and take care of service restarts.

Let's look at how it's set up in systemd. I have a service described in `/etc/systemd/system/ssh-reverse.ssh`

```systemd.service
[Unit]
Description=Reverse SSH connection
After=network.target

[Service]
Type=simple
ExecStart=/usr/bin/ssh -vvv -g -N -T -o "ServerAliveInterval 10" -o "ExitOnForwardFailure yes" -R 22221:localhost:22 198.51.100.77
Restart=always
RestartSec=5s

[Install]
WantedBy=default.target
```

This makes sure that the service is restarted whenever the process exits, for whatever reason, and that it restarts after 5 seconds.

Before starting the service, you should connect to this host manually. `ssh`  will ask if you trust the target host.
If this question is raised while running in systemd, the process will fail to connect, and the service will forever restart.
It is important to run this command as `root`, as that is the user who will be running the service process.

```
# ssh 198.51.100.77
The authenticity of host '198.51.100.77 (198.51.100.77)' can't be established.
ECDSA key fingerprint is SHA256:7IkLx6KW08tCqrDPlmHHBjRFTUHiH4lQFabcdefghijkl.
Are you sure you want to continue connecting (yes/no)?
```

After answering yes here, you should be automatically logged in to the remote host, and the port forward should be established.
If you are not immediately logged in, e.g. asked for password, you should set up authentication using an SSH key.

If you already have a key without a passphrase for the `root` user, you can transfer it to the remote system using

```
# ssh-copy-id $target_host
```

If you don't already have a key without a passphrase, to avoid being asked for login, you should create a key.

```
# ssh-keygen
```

