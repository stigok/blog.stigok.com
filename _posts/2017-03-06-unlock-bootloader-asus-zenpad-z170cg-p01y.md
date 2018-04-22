---
layout: post
title: "Unlock bootloader Asus Zenpad Z170CG (P01Y)"
date: 2017-03-06 09:33:37 +0100
categories: android adb bootloader zenpad
redirect_from:
  - /post/unlock-bootloader-asus-zenpad-z170cg-p01y
---

> This guide has been tested on Linux. It may also work nicely on OSX or Windows, but I haven't tested it out. If you experience any problems, leave a message in the comments.

Before continuing, you should have installed the `adb` and `fastboot` binaries, connected the device to your computer and enabled _ADB debugging_ on the device.

- Lines starting with `#` should be run as root or with `sudo`
- Lines starting with `$` should be run as a normal user
- Lines following the `#` and `$` lines are output from the command that was executed

Verify that the device is accessible

    # adb devices
    List of devices attached
    HELLOWORLDSECRETSERIAL	device

Reboot into fastboot

    # adb reboot fastboot

Unlock bootloader

    # fastboot oem asus-go
    ...
    (bootloader)  Unlocking the bootloader means the following:
    (bootloader)  All user data will be deleted
    (bootloader)  Any securely stored data will be inaccessible
    (bootloader)  Warranty will be void
    (bootloader)  After unlocking you have to execute
    (bootloader)  > fastboot format userdata
    (bootloader)  > fastboot format cache
    (bootloader)  or carry out a factory reset from recovery
    (bootloader)  To confirm the unlock, please execute the command
    (bootloader)  > fastboot oem unlock confirm
    OKAY [  0.057s]
    finished. total time: 0.057s

Confirm unlocking

    # fastboot oem asus-go confirm
    ...
    (bootloader)  Unlocking and rebooting into unlocked state
    OKAY [  0.070s]
    finished. total time: 0.070s

Get a var dump for later

    # fastboot getvar all
    (bootloader)  version-baseband: 24832
    (bootloader)  version-bootloader: 1638.500_M1S1
    (bootloader)  product: SF_3GR
    (bootloader)  secure: YES
    (bootloader)  unlocked: YES
    (bootloader)  off-mode-charge: 1
    (bootloader)  ========== parition type ==========
    (bootloader)  system parition type: ext4
    (bootloader)  userdata parition type: ext4
    (bootloader)  cache parition type: ext4
    (bootloader)  hypervisor parition type: raw
    (bootloader)  boot parition type: raw
    (bootloader)  recovery parition type: raw
    (bootloader)  splash parition type: raw
    (bootloader)  mvconfig parition type: raw
    (bootloader)  secvm parition type: raw
    (bootloader)  vrl parition type: raw
    (bootloader)  psi parition type: raw
    (bootloader)  slb parition type: raw
    (bootloader)  ucode_patch parition type: raw
    (bootloader)  APD parition type: ext4
    (bootloader)  ADF parition type: ext4
    (bootloader)  factory parition type: ext4
    (bootloader)  nvm_static_calib parition type: raw
    (bootloader)  nvm_static_fix parition type: raw
    (bootloader)  nvm_dynamic parition type: raw
    (bootloader)  linux_nvm_fs parition type: ext4
    (bootloader)  ===================================
    (bootloader)  ========== parition size ==========
    (bootloader)  system parition size: 0x0000000090000000
    (bootloader)  userdata parition size: 0x00000002f0380000
    (bootloader)  cache parition size: 0x0000000017200000
    (bootloader)  hypervisor parition size: 0x0000000000100000
    (bootloader)  boot parition size: 0x0000000001100000
    (bootloader)  recovery parition size: 0x0000000001100000
    (bootloader)  splash parition size: 0x0000000002080000
    (bootloader)  mvconfig parition size: 0x0000000000080000
    (bootloader)  secvm parition size: 0x0000000000c00000
    (bootloader)  vrl parition size: 0x0000000000040000
    (bootloader)  psi parition size: 0x0000000000020000
    (bootloader)  slb parition size: 0x0000000000100000
    (bootloader)  ucode_patch parition size: 0x0000000000003000
    (bootloader)  APD parition size: 0x0000000009600000
    (bootloader)  ADF parition size: 0x0000000000100000
    (bootloader)  factory parition size: 0x0000000000900000
    (bootloader)  nvm_static_calib parition size: 0x0000000000800000
    (bootloader)  nvm_static_fix parition size: 0x0000000000100000
    (bootloader)  nvm_dynamic parition size: 0x0000000000100000
    (bootloader)  linux_nvm_fs parition size: 0x0000000000500000
    (bootloader)  ===================================
    (bootloader)  max-download-size: 0x3cc00000
    all:
    finished. total time: 0.261s

Continue booting the device normally

    # fastboot continue

Next step for me is to build and install a custom bootloader. To be continued...

## References

  - https://forum.xda-developers.com/android/development/asus-zenpad-c-7-0-z170c-p01z-root-method-t3311752/post70595649#post70595649