---
layout: post
title:  "Updating and applying default GRUB settings in CentOS 7"
date:   2019-02-21 16:20:26 +0100
categories: centos, grub
---

I wanted to forcefully make `systemd-fsck` perform `fsck` on boot.

Make desired changes in */etc/default/grub* or */etc/grub.d/40_custom*.
I appended `fsck.mode=force fsck.repair=preen` to the `GRUB_CMDLINE_LINUX` variable:

```terminal
$ cat /etc/default/grub
GRUB_TIMEOUT=5
GRUB_DISTRIBUTOR="$(sed 's, release .*$,,g' /etc/system-release)"
GRUB_DEFAULT=saved
GRUB_DISABLE_SUBMENU=true
GRUB_TERMINAL_OUTPUT="console"
GRUB_CMDLINE_LINUX="crashkernel=auto rd.lvm.lv=centos_4/root rd.lvm.lv=centos_4/swap rhgb quiet fsck.mode=force fsck.repair=preen"
GRUB_DISABLE_RECOVERY="true"
```

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
- https://hackerific.net/2016/03/09/fscking-centos-7/

