---
layout: post
title: "Update EFI boot entries with UEFI shell"
date: 2018-06-04 14:44:35 +0200
categories: uefi
---

## Motivation

The motherboard on my Lenovo T470P was replaced and in that process I
(naturally) lost my UEFI boot entries. I wanted them back.

## Add a new entry with UEFI Shell

I used the Arch Linux installation media on a USB stick to get access to a UEFI
shell. Download the Arch `.iso` from [ArchLinux.org downloads](https://www.archlinux.org/download/)
and drop it on a USB stick your favorite tool. I use `dd if=arch.iso of=/dev/sdX bs=4M status=progress`.

Boot it up and the initial menu should give you some different options

- Arch Linux archiso, dumping you into a systemd booted live shell
- UEFI shell v1 or v2
- EFI Default Loader
- Reboot into firware interface

From that list, enter UEFI Shell v2, and you should be dumped into the shell.
The initial output shows you the different file systems and block devices
it sees.

On my computer, the USB device I'm booting from is `FS0` and my actual EFI
partition on my physical drive is `FS1`.
You can look through the file systems by writing e.g. `FS1:`, then use `ls` and
`cd` to browse your way around.

I determined that my EFI image was at `FS1:\EFI\grub\grubx64.efi`, so I created
a new boot entry for that path. The program takes care of mapping `FS1` to a
GUID for you.

    FS1:\> bcfg boot add 0 fs1:\EFI\grub\grubx64.efi
    Target = 0000
    bcfg: Add Boot0000 as 0

List all boot entries (`-v` verbose, `-b` paged)

    bcfg boot dump -v -b

Looks good to me, so I'll reboot

    FS1:\> reset

Hope it worked ;)

## References
- https://wiki.archlinux.org/index.php/Unified_Extensible_Firmware_Interface#Important_UEFI_Shell_commands

