---
layout: post
title: "Update BIOS on a Thinkpad with Linux"
date: 2018-05-13 04:00:42 +0200
categories: thinkpad
---

## Motivation

Wanted to update my BIOS in a hope to fix some ACPI warnings on boot

    des. 06 01:56:14 archlinux kernel: ACPI: Core revision 20170728
    des. 06 01:56:14 archlinux kernel: ACPI Error: [\_SB_.PCI0.XHC_.RHUB.HS11] Namespace lookup failure, AE_NOT_FOUND (20170728/dswload-210)
    des. 06 01:56:14 archlinux kernel: ACPI Exception: AE_NOT_FOUND, During name lookup/catalog (20170728/psobject-252)
    des. 06 01:56:14 archlinux kernel: ACPI Exception: AE_NOT_FOUND, (SSDT:ProjSsdt) while loading table (20170728/tbxfload-228)
    des. 06 01:56:14 archlinux kernel: ACPI Error: 1 table load failures, 10 successful (20170728/tbxfload-246)

First problem is that Lenovo only releases official update tools for Windows.
After some research I eventually found a solution that could be performed
using only a Linux system -- the same system I wanted to upgrade.

## Guide

I have a Lenovo ThinkPad T470p (20J6), and support downloads are located at the
[Lenovo support pages](https://pcsupport.lenovo.com/no/en/products/LAPTOPS-AND-NETBOOKS/THINKPAD-T-SERIES-LAPTOPS/THINKPAD-T470P/20J6/downloads/DS120708).

If you have different model just use the search bar. If you're not sure what
model you have, you can use `dmidecode` to find out. On my system the relevant
output looks like the below

    # dmidecode | grep 'System Information' -A 8
    System Information
    Manufacturer: LENOVO
    Product Name: 20J677LOVE
    Version: ThinkPad T470p
    Serial Number: PF0T****
    UUID: E379BE4C-2D34-11B2-A85C-************
    Wake-up Type: Power Switch
    SKU Number: LENOVO_MT_20J6_BU_Think_FM_ThinkPad T470p
    Family: ThinkPad T470p

Download the "BIOS Update (Bootable CD) [...]" ISO image. You don't need the
utility. In the details section of the download, **read the README** so that
you understand the risks of updating the BIOS. The file I downloaded was
called *r0fuj15wd.iso*.

Download a tool called [`geteltoro`][geteltoro] to extract the boot image from
the ISO file. I went through a great deal of trouble to figure out how to make
a proper bootable USB pen drive from that ISO. This tool was the only one that
actually yielded proper results. You will need Perl in order to run this
program.

    $ geteltoro.pl -o ~/lenovo-bios.img ~/downloads/r0fuj15wd.iso

Now you should have a proper disk image ready to be written to a USB drive.
Write to disk and follow up with a sync of pending IO writes.

    # dd if=~/lenovo-bios.img of=/dev/sdc bs=1M status=progress
    # sync

Reboot and see if you're able to boot into the USB. If not, you may have to
change some settings in BIOS before you'll be able to contine.

  - Disable Secure Boot
  - Enable Legacy Boot option (not just force UEFI)
  - Set the boot order
    - Optionally you can press F12 when the full screen logo appears when
      booting to trigger the Boot device selection menu.

Follow the on-screen instructions to update the BIOS.

There's also an option in there to change the model name suffix of the system.
I don't know if it will interfere with future updates, so *might* not be a
great idea.

## Aftermath

After the update it looked like this instead

    mai 13 03:20:30 prodigy kernel: ACPI: Core revision 20180105
    mai 13 03:20:30 prodigy kernel: ACPI BIOS Error (bug): Failure looking up [\_SB.PCI0.XHC.RHUB.HS11], AE_NOT_FOUND (20180105/dswload-211)
    mai 13 03:20:30 prodigy kernel: ACPI Error: AE_NOT_FOUND, During name lookup/catalog (20180105/psobject-252)
    mai 13 03:20:30 prodigy kernel: ACPI Error: AE_NOT_FOUND, (SSDT:ProjSsdt) while loading table (20180105/tbxfload-228)
    mai 13 03:20:30 prodigy kernel: ACPI Error: 1 table load failures, 10 successful (20180105/tbxfload-246)

Didn't solve my ACPI errors, but my screen is brighter on clearer than before.
Or I'm imagining things. Not sure. Either way I like to have the latest BIOS.

## References
- https://workaround.org/article/updating-the-bios-on-lenovo-laptops-from-linux-using-a-usb-flash-stick/

[geteltorito]: https://userpages.uni-koblenz.de/~krienke/ftp/noarch/geteltorito/
