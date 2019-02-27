---
layout: post
title: "Multi-monitor setup with i3 and xrandr on arch linux"
date: 2017-08-14 15:24:58 +0200
categories: i3wm multimonitor xrandr archlinux i3bar
redirect_from:
  - /post/multimonitor-setup-with-i3-and-xrandr-on-arch-linux
---

I am using i3 window manager and connecting a second screen. i3 is multi-monitor-ready and should work out of the box when a second monitor has been configured with `xrandr`. It is expected that the current setup with a single monitor is working, and that when executing `xrandr`, all connected screens are listed as connected.

Install the xrandr CLI

    # pacman -S xrandr-xorg

List connected screens

    $ xrandr
    Screen 0: minimum 8 x 8, current 3840 x 1080, maximum 32767 x 32767
    DP1 disconnected primary (normal left inverted right x axis y axis)
    DP2 disconnected (normal left inverted right x axis y axis)
    HDMI1 disconnected (normal left inverted right x axis y axis)
    HDMI2 connected 1920x1080+0+0 (normal left inverted right x axis y axis) 620mm x 340mm
       1920x1080     60.00*+  50.00    59.94  
       1920x1080i    60.00    50.00    59.94  
       1680x1050     59.88  
       1600x900      60.00  
       1280x1024     60.02  
       1280x800      59.91  
       1280x720      60.00    50.00    59.94  
       1024x768      60.00  
       832x624       74.55  
       800x600       75.00    60.32  
       720x576       50.00  
       720x576i      50.00  
       720x480       60.00    59.94  
       720x480i      60.00    59.94  
       640x480       75.00    60.00    59.94  
       720x400       70.08  
    VGA1 connected 1920x1080+1920+0 (normal left inverted right x axis y axis) 480mm x 270mm
       1920x1080     60.00*+
       1280x1024     75.02    60.02  
       1152x864      75.00  
       1024x768      75.03    60.00  
       800x600       75.00    60.32  
       640x480       75.00    59.94  
       720x400       70.08  
    VIRTUAL1 disconnected (normal left inverted right x axis y axis)
    DVI-I-1-1 disconnected (normal left inverted right x axis y axis)
    DVI-I-1-2 disconnected (normal left inverted right x axis y axis)

And since my second screen is placed to the right of my main one, I simply specify that as an argument to `xrandr`

    $ xrandr --output VGA1 --right-of HDMI2 --auto

Using `--auto` for making the monitor use its first preferred mode (or, something close to 96dpi if they have no preferred mode).

**A second screen with i3wm magically appears**

The mouse should now be able to travel through to the other screen. Move i3 windows to the second screen with the default i3 keyboard shortcuts: ($mod+Shift+Left / Right). Play around and enjoy the simplicity.

Set preferred primary output

    $ xrandr --output HDMI2 --primary

Make i3 place the tray area of i3bar on the primary output

    bar {
      status_command exec i3status
      tray_output primary
    }

## References
- https://i3wm.org/docs/userguide.html#multi_monitor
- `man xrandr`