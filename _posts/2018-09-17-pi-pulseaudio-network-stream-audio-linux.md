---
layout: post
title:  "Stream audio over network with PulseAudio running on Raspberry Pi"
date:   2018-09-17 19:33:26 +0200
categories: pulseaudio audio streaming
---

## Introduction

I am tired of plugging in the mini-jack to my laptop in order to listen to
some music on my living room speakers. I want to be able to stream over my
local 5GHz wifi network.

## Server installation

The server can be any kind of box with pulseaudio installed. My server is now
a Raspberry Pi connected with the mini-jack to my speakers.

Install some pre-reqs

```
# apt update
# apt install pulseaudio pulseaudio-module-zeroconf avahi-daemon xvfb dbus-launch dbus-x11
```

First I have to create a virtual X11 frame-buffer, because pulseaudio wants
an X11 `DISPLAY` to use. I am creating a systemd service for this.

```
# cat /etc/systemd/system/xvfb@.service
[Unit]
Description=virtual frame buffer X server for display %I
After=network.target

[Service]
ExecStart=/usr/bin/Xvfb %I -screen 0 1280x1024x24

[Install]
WantedBy=multi-user.target
```

I start this service right away and make sure it starts on boot as well. Since
the service name ends with `@`, it can be invoked with an argument that will
be interpolated through the rest of the service file where `%I` appears.

The below is to create a `DISPLAY` at `:77`

```
# systemctl daemon-reload
# systemctl start xvfb@:77.service
# systemctl enable xvfb@:77.service
```

Then defining the pulseaudio service file to make it also start on boot.
In this file, we are giving it a `DISPLAY` environment variable specifying
the display that got created with `xvfb`

```
# cat /etc/systemd/system/pulseaudio.service
[Unit]
Description=PulseAudio Sound System
Before=sound.target

[Service]
Environment=DISPLAY=:77
BusName=org.pulseaudio.Server
ExecStart=/usr/bin/pulseaudio --system
Restart=always

[Install]
WantedBy=default.target
```

Since I had some problems running pulseaudio on the Pi outside of `--system`
mode, I am just running with it, although the logs tells me it is not advised.

Appending some lines to the system config file of pulseaudio to enable network
audio. Allowing link-local IPv6 addresses is vital if you're on an IPv6-enabled
network and Avahi is configured for it.

```
# tee -a /etc/pulse/system.pa
load-module module-native-protocol-tcp auth-ip-acl=127.0.0.1;10.1.1.0/24;fe80::/64
load-module module-zeroconf-publish
```

Enabling and starting this service too

```
# systemctl daemon-reload
# systemctl start pulseaudio.service
# systemctl enable pulseaudio.service
```

Check the logs for errors

```
# journalctl -u pulseaudio
```

## Client setup

**The below might be Arch Linux specific!**

Set `default-server` to the IP address of the remote pulseaudio
server inside `/etc/pulse/client.conf`, then restart the user
service.

```
$ systemctl --user restart pulseaudio.service
```

### Note

The above solution is the only one I've found so far.
However, this removes the possibility of playing sound locally, and
everything is sent to the remote server. This *may* be what you want,
but I'd like to send only select application audio output to the remote
while the rest is local. For this I would use `pavucontrol`.

This post will be updates as soon as I figure out how to have them both
available as outputs at the same time.


## References
- https://raspberrypi.stackexchange.com/questions/11735/using-pi-to-stream-all-audio-output-from-my-pc-to-my-stereo
- https://wiki.archlinux.org/index.php/PulseAudio#Networked_audio
- https://raspberrypi.stackexchange.com/questions/639/how-to-get-pulseaudio-running/44767#44767
- https://superuser.com/questions/319040/proper-way-to-start-xvfb-on-startup-on-centos/912648#912648
