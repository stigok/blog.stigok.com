---
layout: post
title:  "Updating and applying default GRUB settings in CentOS 7"
date:   2019-02-21 16:20:26 +0100
categories: centos, grub
---

**NOTE**: This is (most probably) distro specific to CentOS 7

I wanted to force `systemd-fsck` to run `fsck` on the next reboot.

Make desired changes in */etc/default/grub* or */etc/grub.d/40_custom*.
I appended `fsck.mode=force fsck.repair=preen` to the `GRUB_CMDLINE_LINUX` variable.

Depending on your system, apply changes using
- **BIOS**: `grub2-mkconfig -o /boot/grub2/grub.cfg`
- **EFI**:  `grub2-mkconfig -o /boot/efi/EFI/centos/grub.cfg`

Output on my BIOS system:
```terminal
$ sudo grub2-mkconfig -o /boot/grub2/grub.cfg
Generating grub configuration file ...
Found linux image: /boot/vmlinuz-3.10.0-957.5.1.el7.x86_64
Found initrd image: /boot/initramfs-3.10.0-957.5.1.el7.x86_64.img
Found linux image: /boot/vmlinuz-3.10.0-957.1.3.el7.x86_64
Found initrd image: /boot/initramfs-3.10.0-957.1.3.el7.x86_64.img
Found linux image: /boot/vmlinuz-3.10.0-862.14.4.el7.x86_64
Found initrd image: /boot/initramfs-3.10.0-862.14.4.el7.x86_64.img
Found linux image: /boot/vmlinuz-3.10.0-862.11.6.el7.x86_64
Found initrd image: /boot/initramfs-3.10.0-862.11.6.el7.x86_64.img
Found linux image: /boot/vmlinuz-3.10.0-862.el7.x86_64
Found initrd image: /boot/initramfs-3.10.0-862.el7.x86_64.img
Found linux image: /boot/vmlinuz-0-rescue-d46fbfa672984127815985c826ed7514
Found initrd image: /boot/initramfs-0-rescue-d46fbfa672984127815985c826ed7514.img
done
```

After a, hopefully, successful reboot, I'd want to remove these settings again.


## References
- `man systemd-fsck`
- https://wiki.centos.org/HowTos/Grub2
- https://unix.stackexchange.com/a/152249/28043

