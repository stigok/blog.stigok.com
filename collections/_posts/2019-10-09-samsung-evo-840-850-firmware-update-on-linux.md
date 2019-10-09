---
layout: post
title:  "Update firmware of Samsung EVO 840 and 850 SSD drives on Linux"
date:   2019-10-09 23:06:25 +0200
categories: ssd, firmware
---

## Firmware upgrade procedure
### Samsung EVO 840

For the EVO 840 I downloaded the firmware for Mac, which was an ISO file.
I mounted the ISO to a local directory, copied out the *.img* file I found,
then wrote the *.img* file directly to a USB flash drive.
My USB flash drive specified as */dev/sdx*.

```
# mount Samsung_SSD_840_EVO_Series_EXT0BB6Q_Mac.iso /mnt
# dd if=/mnt/isolinux/btdsk.iso of=/dev/sdx
# sync
# umount /mnt
```

Download from <https://www.samsung.com/no/support/model/MZ-7TE500LW/#downloads>.

I rebooted the system and enabled *Legacy USB boot* in my UEFI (BIOS),
and followed the procedure. At one point, I was asked to power cycle the
disk before confirming I had done so. I pulled the power plug of the disk
and put it back in, before I confirmed with `Y`.

### Samsung EVO 850

For the EVO 850, the firmware update ISO can be directly written to a
USB flash drive.

```
# dd if=Samsung_SSD_850_EVO_EMT02B6Q_Mac.iso of=/dev/sdx
# sync
```

<https://www.samsung.com/no/support/model/MZ-75E250B/EU/#downloads>

## Tools

To check the firmware version of my disks, I use `smartctl` from the
*smartmontools* package in Arch.

```
# smartctl /dev/sdb --info
smartctl 7.0 2018-12-30 r4883 [x86_64-linux-5.3.4-arch1-1-ARCH] (local build)
Copyright (C) 2002-18, Bruce Allen, Christian Franke, www.smartmontools.org

=== START OF INFORMATION SECTION ===
Model Family:     Samsung based SSDs
Device Model:     Samsung SSD 840 EVO 500GB
Serial Number:    [redacted]
LU WWN Device Id: [redacted]
Firmware Version: EXT0BB6Q
User Capacity:    500,107,862,016 bytes [500 GB]
Sector Size:      512 bytes logical/physical
Rotation Rate:    Solid State Device
Device is:        In smartctl database [for details use: -P show]
ATA Version is:   ACS-2, ATA8-ACS T13/1699-D revision 4c
SATA Version is:  SATA 3.1, 6.0 Gb/s (current: 6.0 Gb/s)
Local Time is:    Wed Oct  9 23:28:48 2019 CEST
SMART support is: Available - device has SMART capability.
SMART support is: Enabled
```

The above output was taken *after* the upgrade of the EVO 840.

## References
- <https://unix.stackexchange.com/questions/333853/update-firmware-of-samsung-840-pro>
- <https://www.archlinux.org/packages/extra/x86_64/smartmontools/>
