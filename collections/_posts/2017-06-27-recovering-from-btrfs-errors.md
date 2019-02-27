---
layout: post
title: "Recovering from btrfs errors"
date: 2017-06-27 14:04:53 +0200
categories: btrfs recovery luks
redirect_from:
  - /post/recovering-from-btrfs-errors
---

After vacuuming under my desk, my SSD went into read-only mode while the system
was on. I have a [LVM on Luks and btrfs](https://fogelholk.io/installing-arch-with-lvm-on-luks-and-btrfs/)
setup, and after rebooting, it went straight into btrfs emergency rescue mode.

![btrfs lvm luks errors](https://public.stigok.com/img/1498564677730815124.jpg)

## Attempt backup

Before continuing, I mirrored the drive to a backup drive with a 2048 block size.
The reason for that block size is that it is what I determined to be the sector
size of the SSD. Whether that is the correct assumption or not I cannot guarantee.

    # dd if=/dev/sda of=/dev/sdb bs=2048

I created a btrfs filesystem on a usb pen drive, as the btrfs recovery shell
didn't support vfat nor ext4. Inserted the pen drive in the recovering computer
and mounted it. I had two other disks connected, so my USB drive became **sdc**.

    # mkdir /usbmnt
    # mount /dev/sdc1 /usbmnt

Then attempt to recover my git repositories that I, of course, hadn't pushed
back to origin.

    # btrfs restore -c --path-regex '^/(|home/(|/sshow(|/repos/(|/.*))))$' \
        /dev/mapper/lvmvg-rootvol /usbmnt

And my dotfiles (this copies over all directories starting with a dot too, but
not their contents, so it's okay for me)

    # btrfs restore -c --path-regex '^/(|home/(|/sshow(|/\.\w*)))$' \
        /dev/mapper/lvmvg-rootvol /usbmnt
        
And I want my `.config` folder

    # btrfs restore -c --path-regex '^/(|home/(|/sshow(|/.config(|/.*))))$' \
        /dev/mapper/lvmvg-rootvol /usbmnt

## Attempt repair

> **IMPORTANT:** At this point, I have the backup that I want, and I'm continuing
> to figure btrfs out by blindly attempting recovery options with no fear of
> permanently damaging my data or drive, even though that might actually happen.

Now that I have the backup, I'll try to run the `--repair` tool.

    # btrfs check --repair /dev/mapper/lvmvg-rootvol

It outputted the same errors as `check` without `--repair`, so I stepped up my
faith game and attempted `init-csum-tree` and `init-extent-tree` as well. I
haven't even read it up, but according to the errors from a plain `btrfs check`,
this is exactly what the disk was having problems with.

    # btrfs check --repair --init-csum-tree --init-extent-tree /dev/mapper/lvmvg-rootvol

![btrfs recovery complete](https://public.stigok.com/img/1498564318738706074.jpg)
    
Now ended with success and reporting 0 errors. Time for a reboot.

It went pretty good. Now booting properly again.

## Aftermath

System seems okay. After boot I had a corrupted journal file, which was rotated,
but other than that, it seems fine. Some questions to ask myself though:

- Should I continue to use this disk, even though it *might* have hardware issues?
- Should I continue to use btrfs + luks + lvm, even though btrfs is a bit experimental?

Yes to both. Let's live on in certain uncertainty &lt;3

## References
- https://btrfs.wiki.kernel.org/index.php/Restore
- http://linux-btrfs.vger.kernel.narkive.com/MiGhpa5X/help-with-corrupt-filesystem-btrfs-free-extent-5236-io-failure