---
layout: post
title: "Troubleshooting bluescreen of death ntoskrnl.exe Windows 7"
date: 2017-04-05 15:03:44 +0200
categories: windows bsod
redirect_from:
  - /post/troubleshooting-bluescreen-of-death-ntoskrnlexe-windows-7
---

Had issues with my Windows 7 box got lots of bluescreens, but it happened mostly while I wasn't around. I came back and saw my computer either being rebooted or hanging on the blue screen.

*TLDR;* scroll to bottom

## Reading the memory dump

Make sure you have enabled memory dumps on system failures. Open *System* in *All Control Panel Items* (shortcut Win+Break), then open *Advanced system settings* in the sidebar.

![Advanced system settings](https://s.42.fm/img/advanced-system.png)

Select *Advanced*, then *Startup and Recovery -> Settings...*. Enable *Write debugging information" and select *Kernel memory dump*.

![Write debugging information](https://s.42.fm/img/startup-recovery.PNG)

On the next BSOD, the current system memory will be dumped into a .dmp file for debugging purposes. This will help in determining the cause of the error.

I use [BlueScreenView](http://www.nirsoft.net/utils/blue_screen_view.html) to view the dumps afterward

![BlueScreenView](https://s.42.fm/img/bluescreenview.PNG)

All of my incidents seemed to occur within `ntoskrnl.exe`. A lot of articles were suggesting issues related to storage and NTFS, but it didn't really seem like the cause for mine.

## Troubleshooting

- Memory faulty?
  - Ran 13 passes of memtest86 without errors
- SATA faulty?
  - Switched SATA cables
  - Chose a different port on the motherboard
- Switch BIOS settings for on-board devices
  - Disabled COM Super I/O
  - Disconnected front USB hub (as I've had issues with connectivity on those ports)
  - Intel system drivers update
- Prevent OS from sleeping
  - I was seeing a pattern of when was getting blue screens; either when going to, or waking up from sleep.
  - **Disabling sleep solved the problems**

So **preventing my OS from going to sleep solved my problems**, and I don't mind it not sleeping. It's still turning off the monitors without issues, so I'm happy.