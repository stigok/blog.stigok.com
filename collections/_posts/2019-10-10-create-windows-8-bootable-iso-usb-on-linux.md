---
layout: post
title:  "Create a Windows 8 bootable USB stick from an ISO on Linux"
date:   2019-10-10 12:28:39 +0200
categories: windows recovery
---

I downloaded the Windows 8 ISO from <https://www.microsoft.com/en-us/software-download/windows8ISO>.

Format and prepare the USB flash drive

```
# fdisk /dev/sdx
```

- `o` new DOS disklabel
- `n` new partition
- `[default]` partition type
- `[default]` partition number
- `[default]` first sector
- `+8GiB` last sector, since the ISO is around 4GiB
- *CONDITIONAL:* if notified about an existing signature, confirm changes with `Y`
- `t` to select partition type code
- `c` for FAT32 (LBA)
- `a` to make bootable
- `w` to write changes and exit

The output of the `p` command in `fdisk` should now read something similar to:

```
# fdisk -l /dev/sdx1
Disk /dev/sdx: 149.5 GiB, 160041885696 bytes, 312581808 sectors
Disk model: [redacted]
Units: sectors of 1 * 512 = 512 bytes
Sector size (logical/physical): 512 bytes / 512 bytes
I/O size (minimum/optimal): 512 bytes / 33553920 bytes
Disklabel type: dos
Disk identifier: 0x12345678

Device     Boot Start      End  Sectors Size Id Type
/dev/sdx1  *    65535 10551134 10485600   5G  c FAT32 (LBA)
```

Create NTFS file system. Use `-f` for fast format, to really speed things up.
(NTFS disk tools can be installed from the AUR as `ntfs-3g`)

```
# mkfs.ntfs -f /dev/sdx1
```

Mount both the ISO and the NTFS partition

```
# mkdir -p /mnt/disk /mnt/loop
# mount /dev/sdx1 /mnt/disk
# mount -o loop ~/downloads/Win8.1_Pro_N_EnglishInternational_x64.iso /mnt/loop
```

Copy the files from the loop device (ISO) to the USB flash drive NTFS partition.
I like to use `rsync` for its progress output, but you might as well use `cp -av`.

```
# cp -av /mnt/loop/* /mnt/disk/
```

Don't mind *Operation not permitted* messages when attempting to set file permissions.
This is expected with the FAT32 filesystem.

Synchronise cached writes and unmount the disk

```
# sync /dev/sdx
# unmount /mnt/disk
```

## References
- <https://thornelabs.net/posts/create-a-bootable-windows-7-or-10-usb-drive-in-linux.html>
- <https://itsfoss.com/bootable-windows-usb-linux/>
