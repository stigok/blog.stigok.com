---
layout: post
title: "Increase 5v rail amp output on Raspberry Pi 3/Zero"
date: 2017-10-24 21:41:02 +0200
categories: raspberrypi
redirect_from:
  - /post/increase-5v-rail-amp-output-on-raspberry-pi-3zero
---

Put `max_usb_current=1` in `/boot/config.txt`. This will lift the current limit from 600mAh to 1200mAh. 

> All that max_usb_current=1 does is to set GPIO38 high, which in turn turns on a FET, which connects a second 39K resistor in parallel to an existing one, on pin 5 of U13, the AP2553W6 USB power manager, lifting the current limit from 0.6A to double that (1.2A), see no possible scenario there why the PI resets because of that, except in case the gate of the FET Q4 is somehow shorted to GND. Which could be caused by a production fault. Inspect Q4, as look if there is solder shorting pins together. Also R6 (resistor mounted between gate of Q4 and GND) should be 100K not 0 Ohm. U13, Q4 and R6 should be near the USB ports.

## References
- https://www.raspberrypi.org/forums/viewtopic.php?p=594183#p594183