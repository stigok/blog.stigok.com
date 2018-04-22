---
layout: post
title: "Format disk for Windows, Mac and Linux on Linux"
date: 2017-03-23 13:20:22 +0100
categories: linux partitioning storage
redirect_from:
  - /post/format-disk-for-windows-mac-and-linux-on-linux
---

> **NOTE:** If you mix up the disk names in the below commands the results may be fatal for your data. If you are not certain you are writing the correct drive paths, don't run the commands.

Commands starting with `$` are executed as normal user and `#` run as root (or sudo).

Identify the device you want to format

    $ lsblk
    NAME   MAJ:MIN RM   SIZE RO TYPE MOUNTPOINT
    sda      8:0    0 232.9G  0 disk 
    ├─sda1   8:1    0   512M  0 part /boot
    ├─sda2   8:2    0    20G  0 part /
    ├─sda3   8:3    0     8G  0 part [SWAP]
    └─sda4   8:4    0 204.4G  0 part /home
    sdb      8:16   0   1.8T  0 disk 
    ├─sdb1   8:17   0 316.7G  0 part 
    └─sdb2   8:18   0   1.5T  0 part

In my case, my OS installation is on *sda*, and the disk I want to format is *sdb*.
The digit after the disk name is the partition number.

I want to delete all partitions on the disk and create a new one. I will use `fdisk` for this.

    # fdisk /dev/sdb

The program will output some usage information, which might be smart to read. Note that nothing will be written to disk before you explicitly tell it to do so with the `w` command.

1. Press `o` to "create a new empty DOS partition table"
2. Press `n` to "add a new partition"
3. Press Enter for default partition type
4. Press Enter for default partition number
5. Press Enter for default first sector
6. Press Enter for default last sector
  - A new partition has been created
7. Press `t` to change type of a partition
  - Partition 1 should now have been automatically selected
8. Press `7` to select "HPFS/NTFS/exFAT" as the partition type
9. Press `w` to write table to disk and exit

.

    Command (m for help): w
    The partition table has been altered.
    Calling ioctl() to re-read partition table.
    Syncing disks.

After the partition has been created, create a file system on the partition.
To get the required utility for this, you'll need the `mkntfs` binary. On Arch it is comes with `extra/ntfs-3g`.
I'm running `mkntfs` with the `--debug` option for my curiosity's sake.

    # mkntfs --quick --debug /dev/sdb1
    Cluster size has been automatically set to 4096 bytes.
    Creating NTFS volume structures.
    Creating root directory (mft record 5)
    Creating $MFT (mft record 0)
    Creating $MFTMirr (mft record 1)
    Creating $LogFile (mft record 2)
    Creating $AttrDef (mft record 4)
    Creating $Bitmap (mft record 6)
    Creating $Boot (mft record 7)
    Creating backup boot sector.
    Creating $Volume (mft record 3)
    Creating $BadClus (mft record 8)
    Creating $Secure (mft record 9)
    Creating $UpCase (mft record 0xa)
    Creating $Extend (mft record 11)
    Creating system file (mft record 0xc)
    Creating system file (mft record 0xd)
    Creating system file (mft record 0xe)
    Creating system file (mft record 0xf)
    Creating $Quota (mft record 24)
    Creating $ObjId (mft record 25)
    Creating $Reparse (mft record 26)
    Syncing root directory index record.
    Syncing $Bitmap.
    Syncing $MFT.
    Updating $MFTMirr.
    Syncing device.
    mkntfs completed successfully. Have a nice day.

And now I have a disk that works with Windows, Mac and Linux (with NTFS support).