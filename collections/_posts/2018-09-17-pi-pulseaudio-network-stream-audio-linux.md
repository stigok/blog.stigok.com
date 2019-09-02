---
layout: post
title:  "Stream audio over network with PulseAudio running on Raspberry Pi"
date:   2018-09-17 19:33:26 +0200
categories: pulseaudio audio streaming raspberrypi
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
# apt install pulseaudio pulseaudio-module-zeroconf avahi-daemon
```

Define the pulseaudio service file to make it start on boot and log to the
journal.

```
# tee /etc/systemd/system/pulseaudio.service
[Unit]
Description=PulseAudio system-wide server

[Service]
Type=forking
PIDFile=/var/run/pulse/pid
ExecStart=/usr/bin/pulseaudio --daemonize --system --realtime --log-target=journal
ExecStop=/usr/bin/pulseaudio -k
Restart=on-failure
LimitRTPRIO=1000
LimitNICE=-20

[Install]
WantedBy=multi-user.target
```

Since I had some problems running pulseaudio on the Pi in user mode, I'm just
going with `--system` mode, although the logs tells me it is not advised.
This Pi is only used headlessly by me anyway.

Appending some lines to the pulseaudio system configuration file to enable network
audio. Allow all IPv4 private internet addresses, and allowing link-local IPv6 addresses
is vital if you're on an IPv6-enabled network and Avahi is configured for it.

```
# tee -a /etc/pulse/system.pa
load-module module-native-protocol-tcp auth-ip-acl=127.0.0.0/8;10.0.0.0/8;172.16.0.0/12;192.168.0.0/16;fe80::/10
load-module module-zeroconf-publish
```

Start it right away, but also enable it to make it start on boot automatically.

```
# systemctl daemon-reload
# systemctl start pulseaudio.service
# systemctl enable pulseaudio.service
```

Check the logs for errors

```
# journalctl -u pulseaudio
```

If you need more verbosity in the logs, you can set `log-level = debug` in the
pulseaudio *daemon.conf*.

Start Avahi to enable automatic service discovery

```
# systemctl start avahi-daemon.service
# systemctl enable avahi-daemon.service
```

## Client setup

**The below might be Arch Linux specific!**

Start Avahi to enable automatic remote speaker discovery

```
# systemctl start avahi-daemon.service
# systemctl enable avahi-daemon.service
```

Restart the user service

```
$ systemctl --user restart pulseaudio.service
```

Some times I was experiencing choppy sound. I did the following to remedy this

```
# tee /etc/libao.conf
default_driver=pulse
quiet
buffer_time=50
dev=combined
server=localhost
```

And whenever I'm getting choppy sound again, as I usually do for the initial
song I play on SoundCloud, I restart the pulseaudio user service again and
all is fine.

Now I'm able to use `pavucontrol` to select audio outputs for specific
applications on my system.

## References
- https://raspberrypi.stackexchange.com/questions/11735/using-pi-to-stream-all-audio-output-from-my-pc-to-my-stereo
- https://wiki.archlinux.org/index.php/PulseAudio#Networked_audio
- https://raspberrypi.stackexchange.com/questions/639/how-to-get-pulseaudio-running/44767#44767
- https://superuser.com/questions/319040/proper-way-to-start-xvfb-on-startup-on-centos/912648#912648
- https://wiki.archlinux.org/index.php/PulseAudio/Examples
- https://partofthething.com/thoughts/multi-room-audio-over-wi-fi-with-pulseaudio-and-raspberry-pis/
