---
layout: post
title: "Create bootable USB drive with ThinkPad Drive Erase Utility on Linux"
date: 2017-12-04 13:04:39 +0100
categories: thinkpad ssd encryption
redirect_from:
  - /post/create-bootable-usb-drive-with-thinkpad-drive-erase-utility-on-linux
---

I am creating a bootable USB pen drive containing [Drive Erase Utility for the Resetting the Cryptographic Key and the Erasing the Solid State Drive - ThinkPad](https://pcsupport.lenovo.com/no/en/products/LAPTOPS-AND-NETBOOKS/THINKPAD-T-SERIES-LAPTOPS/THINKPAD-T470P/downloads/DS019026) for my Lenovo T470P.

This is one of multiple steps taken to prepare for full disk encryption.

Download the .zip file containing the utility. This file is a little under 80KB.

## Prepare USB disk

Run `fdisk` to format and prepare the USB disk. **Make sure** the device path to the USB disk is correct. `/dev/sdb` is the correct one for me, but may not be the same device on your computer.

1. `sudo fdisk /dev/sdb` to start the utility. Next steps are keypresses.
2. <kbd>o</kbd> to create a new DOS disklabel
3. <kbd>n</kbd> for new partition
4. <kbd>Enter</kbd> for default (primary partition)
5. <kbd>Enter</kbd> for default partition number
6. <kbd>Enter</kbd> for default first sector
7. <kbd>Enter</kbd> for default last sector
8. **Conditional:** If prompted, enter <kbd>Y</kbd> to overwrite existing file system signature
9. <kbd>t</kbd> to change partition type
10. <kbd>c</kbd> for *W95 FAT32 (LBA)*
11. <kbd>a</kbd> to toggle (enable) bootable flag
12. <kbd>w</kbd> to write changes to disk and exit fdisk

Create a filesystem on the drive. Now, append a `1` to the original device path to select the first and single partition created in the previous step.

    # mkfs.vfat -F 32 /dev/sdb1

Mount the partition

    # mount /dev/sdb1 /mnt

Unzip the contents of the utility zip previously downloaded

    # cd /mnt
    # unzip ~/Downloads/83fd04ww.zip

Create UEFI folder and move the `BootX64.efi` into that folder

    # mkdir -p /mnt/EFI/BOOT
    # mv /mnt/BootX64.efi /mnt/EFI/BOOT

Sync pending disk operations, change out of the directory, and unmount the flash drive

    # sync
    # cd /
    # unmount /mnt

## Run the utility

Insert the USB drive into the ThinkPad and hammer F12 (alternatively F1, then F12) after powering on to be able to select "Boot from other device". Select the USB Flash Drive, and the utility should start.

Disk utility steps
- 1 for delete
- Yes to confirm
- Yes to REALLY confirm
- Write down the Request Key written on screen
- Press ENTER to Restart and let the boot process do its thing without interfering
- Enter the request key previously written down and ENTER to continue
- Enter to confirm destroying of all data
- Wait for operation to complete and press any key to restart