---
layout: post
title: "Enable SSH server sshd on raspbian"
date: 2017-10-12 17:33:44 +0200
categories: raspbian raspberrypi ssh
redirect_from:
  - /post/enable-ssh-server-sshd-on-raspbian
---

In raspbian, theres a systemd service `sshswitch.service` that checks whether a certain file exists or not.

    [Unit]
    Description=Turn on SSH if /boot/ssh is present
    ConditionPathExistsGlob=/boot/ssh{,.txt}
    After=regenerate_ssh_host_keys.service
    
    [Service]
    Type=oneshot
    ExecStart=/bin/sh -c "update-rc.d ssh enable && invoke-rc.d ssh start && rm -f /boot/ssh ; rm -f /boot/ssh.txt"
    
    [Install]
    WantedBy=multi-user.target

This means that to enable the `sshd.service` on a fresh install of raspbian you can create an empty file at `/boot/ssh`.

You can do this without power the Pi, by, for example:

    # mkdir -p /mnt/boot
    # mount /dev/mmcblk0p1 /mnt/boot
    # touch /mnt/boot/ssh
    $ sync
    # umount /mnt/boot

And the `sshd` should be running on the next boot.