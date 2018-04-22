---
layout: post
title: "Flashing SD card with Arch Linux ARM for Raspberry Pi 2 and 3"
date: 2017-06-01 14:59:34 +0200
categories: arch linux raspberrypi
redirect_from:
  - /post/flashing-sd-card-with-arch-linux-arm-for-raspberry-pi-2-and-3
---

## Prereqs
- dosfstools (which supplies mkfs.vfat)
- fdisk

## Preparing SD card

Start fdisk to partition the SD card:

    fdisk /dev/sdX

- Type o. This will clear out any partitions on the drive.
- Type p to list partitions. There should be no partitions left.
- Type n, then p for primary, 1 for the first partition on the drive, press ENTER to accept the default first sector, then type +100M for the last sector.
- Type t, then c to set the first partition to type W95 FAT32 (LBA).
- Type n, then p for primary, 2 for the second partition on the drive, and then press
- ENTER twice to accept the default first and last sector.
- Write the partition table and exit by typing w.

Create and mount the FAT filesystem:

    mkfs.vfat /dev/sdX1
    mkdir boot
    mount /dev/sdX1 boot

Create and mount the ext4 filesystem:

    mkfs.ext4 /dev/sdX2
    mkdir root 
    mount /dev/sdX2 root

Download archive

    curl -LO http://os.archlinuxarm.org/os/ArchLinuxARM-rpi-3-latest.tar.gz

Extract the root filesystem (as root, not via sudo)

    bsdtar -xpf ArchLinuxARM-<rpi-version>-latest.tar.gz -C root
    sync

Move boot files to the first partition:

    mv root/boot/* boot

Unmount the two partitions:

    umount boot root

## Installing

Refer to the [Arch Linux Installation Guide](https://wiki.archlinux.org/index.php/Installation_guide) for detailed information.

Connect the Pi to the local network and do an `arp-scan` to find its IP address

    # arp-scan -l --interface=enp3s0

SSH to it and login with alarm:alarm

    $ ssh alarm@$PI_ADDRESS

Change password immediately

    alarm@alarm$ passwd

Use `su` to get root and change root password. Default root password is root

    alarm@alarm$ su
    root@alarm# passwd

## References
- <http://elinux.org/ArchLinux_Install_Guide>