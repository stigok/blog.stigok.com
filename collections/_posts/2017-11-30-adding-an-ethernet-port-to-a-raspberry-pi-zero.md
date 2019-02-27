---
layout: post
title: "Adding an ethernet port to a Raspberry Pi Zero"
date: 2017-11-30 21:44:16 +0100
categories: raspberrypi hardware spi ethernet
redirect_from:
  - /post/adding-an-ethernet-port-to-a-raspberry-pi-zero
---

This also applies to a Raspberry Pi A+, which doesn't come with any network interfaces at all.

This is an adaptation of [a post at raspi.tv](http://raspi.tv/2015/ethernet-on-pi-zero-how-to-put-an-ethernet-port-on-your-pi), describing the same procedure, but with different hardware and parts. All credits to the original author.

## Parts

I'm using a mini ENC28J60 ethernet chip with headers already soldered on for easy wiring. In this article, I'm connecting it to a Raspberry Pi Model A+ V1.1, but the procedure, wiring and pin numbering is identical to the Pi 2, Pi 3 and both of the Pi Zero's.

This ethernet chip has a 25MHz crystal. This number matters in the interface configuration later in the post. If you're using a different chip, make note of what kind of crystal it has.

Pictures pending...

## Wiring

My wiring is a bit different from the original post, as I have a different chip.

    Chip Pin  Pi Pin (BCM)
    NT        GP25
    SO        GP9
    SCK       GP11
    VCC       3V3
    SI        GP10
    CS        GP8
    GND       GND

Pictures pending...

## Configuration

Add the `dtoverlay` enabling the chip. This seem to imply `dtparam=spi=on`, so I'm not specifying it explicitly.

    # The ethernet chip has a 25MHz crystal
    dtoverlay=enc28j60,int_pin=25,speed=25000000

The crystal on my chip says `25 000` (picture pending...), which means 25 000 KHz, which equals 25 MHz, which again equals `25 000 000`.

## Test

Time to see if it works. Before powering up the Pi, **make sure** the wiring is correct. Really, double (triple), check the wiring to avoid damaging the Pi and the Card.

If you're running a headless setup, you can see if it's getting an IP address by [snooping DHCP on the local network](https://blog.stigok.com/post/watch-dump-dhcp-requests-on-local-network).

(pictures pending...)

### Connection speed

I'm getting **630kbps** with default clock, and **969kbps** with the 25MHz clock. This equals roughly 8Mbit and is acceptable for my use.

I tested this using `curl` on the local network. This test *might* not be a good one, as the bottleneck may have been SD card write speeds on the Pi itself.

## References
- raspi.tv/2015/ethernet-on-pi-zero-how-to-put-an-ethernet-port-on-your-pi