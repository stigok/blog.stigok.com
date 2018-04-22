---
layout: post
title: "Mounting filesystem in ram with tmpfs"
date: 2017-07-12 16:30:30 +0200
categories: tmpfs mount linux raspberrypi
redirect_from:
  - /post/mounting-filesystem-in-ram-with-tmpfs
---

I'm recording some video on my Raspberry Pi and I don't want to write to the SD-card. I want to do some processing on the video before sending it out, though, so I need to save it somewhere. Saving it in RAM is ideal both for speed and its transientness.

First create the folder to hold the new file system

    $ mkdir /home/alarm/tmp

Then add the mount record to `/etc/fstab` for automatically mounting on boot

    # echo 'tmpfs /home/alarm/tmp tmpfs nodev,nosuid,uid=1000,gid=1000,size=50m,noexec,mode=1700 0 0' >> /etc/fstab

Mount the fs

    # mount -a

Test

    $ echo 'Hello, world!' >> ~/tmp/foobar

Now, reboot, and all files within the mount path will be gone. **Forever**

## References
- `man mount`