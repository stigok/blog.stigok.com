---
layout: post
title: "Clone corrupted SD-card and ignore read errors"
date: 2017-12-27 19:54:54 +0100
categories: recovery raspberrypi
redirect_from:
  - /post/clone-corrupted-sdcard-and-ignore-read-errors
---

I had an SD-card in one of my Raspberry Pi's die on me today. It has been in operation for over a year, and for some months, it has had systemd-journal persistent storage enabled with 5 minute write intervals. My immediate thoughts is that the journal is to blame, but who knows. Here's how I cloned the SD-card anyway.

> I would maybe use a different approach for a bricked device with spinning parts, like a hard-disk. Since this was an SD-card with non-critical data on it, I have not considered the possibility of further damaging the device while attempting to clone it. Take caution!

## Clone block device

I had success mounting the boot partition of the SD card, therefore I will only attempt to backup the root partition of the block device (`/dev/sda2`).

I'll us `dd` for this. Since I expect errors to occur, I don't want it to exit on read failures

    # dd if=/dev/sda2 of=corrupted-pi-root-partition.img conv=sync,noerror status=progress

The `sync` option pads invalid reads with zeroes (0) while maintaining the offsets. The `noerror` prevents `dd` from exiting on errors.

While cloning, I'm running `dmesg -w` to follow detailed bad sector information while the card is copied

    # dmesg -w

Here is some of the output I encountered while cloning a my corrupted card:

    ### output from dd
    dd: error reading '/dev/sda2': Input/output error
    8762028+269 records in
    8762297+0 records out
    4486296064 bytes (4.5 GB, 4.2 GiB) copied, 3965.33 s, 1.1 MB/s
    4486296576 bytes (4.5 GB, 4.2 GiB) copied, 3965 s, 1.1 MB/s


    ### output from dmesg
    [37716.939106] sd 0:0:0:0: timing out command, waited 180s
    [37716.939120] sd 0:0:0:0: [sda] tag#0 UNKNOWN(0x2003) Result: hostbyte=0x00 driverbyte=0x08
    [37716.939123] sd 0:0:0:0: [sda] tag#0 Sense Key : 0x4 [current] 
    [37716.939125] sd 0:0:0:0: [sda] tag#0 ASC=0x4b ASCQ=0x0 
    [37716.939128] sd 0:0:0:0: [sda] tag#0 CDB: opcode=0x28 28 00 00 48 f0 48 00 00 08 00
    [37716.939131] print_req_error: I/O error, dev sda, sector 4780104
    [37716.939136] Buffer I/O error on dev sda2, logical block 571657, async page read

From the output above, it's possible to determine which sectors are bad, and possibly use that information for advanced recovery.

A problem I did not solve was to decrease the bad sector read timeout from 180 seconds. Fortunately, I only had about ten bad sectors, so it didn't take *too* long, but for a device in really bad shape, this might be way too much. If you know how to decrease it, please post a comment.