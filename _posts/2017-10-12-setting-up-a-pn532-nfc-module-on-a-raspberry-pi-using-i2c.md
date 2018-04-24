---
layout: post
title: "Setting up a PN532 NFC module on a Raspberry Pi using I2C"
date: 2017-10-12 17:56:43 +0200
categories: raspberrypi nfc rfid i2c
redirect_from:
  - /post/setting-up-a-pn532-nfc-module-on-a-raspberry-pi-using-i2c
---

![PN532](https://public.stigok.com/img/1507824522111557533.jpg)

## Parts

- 1 x Raspberry Pi 3 Model B
- 1 x [PN532 NFC/RFID module](https://www.m.nu/rfid-nfc/pn532-nfc-rfid-module-v3)

## Preparing the Pi

- Download a fresh version of [raspbian][]
- Enable I2C interface on the Pi `raspi config`
- Enable SSH server `echo 'This enables SSH on boot' | sudo tee /boot/ssh`
- Install I2C utility binaries `sudo apt install i2c-tools`
- Configure the NFC module to use I2C by physically flipping a SMB header on the PCB

### Wire up the PN532

**Turn off the power to the Pi while you are wiring**

- [Raspberry Pi GPIO I2C map][GPIO map]

This NFC module has multiple interfaces to connect with. For I2C I'm physically connecting the GPIO pins to the 4-pin interface on the NFC module PCB.

    # NFC module pin -> Pi GPIO physical pin #
    GND -> 6
    VCC -> 4
    SDA -> 3
    SCL -> 5

Probe for I2C devices:

    root@raspberrypi:~# i2cdetect -y 1
         0  1  2  3  4  5  6  7  8  9  a  b  c  d  e  f
    00:          -- -- -- -- -- -- -- -- -- -- -- -- -- 
    10: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
    20: -- -- -- -- 24 -- -- -- -- -- -- -- -- -- -- -- 
    30: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
    40: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
    50: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
    60: -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- -- 
    70: -- -- -- -- -- -- -- --

In the above results, a device exists at address `0x24` which is the NFC module.

![Raspberry Pi 3 with PN532 connected over I2C](https://public.stigok.com/img/1507896705867113013.jpg)

### libnfc

Install NFC tools:

    sudo apt install libnfc5 libnfc-bin libnfc-examples

I had a problem detecting the I2C device with `nfc-list -v` and `nfc-scan-device -v` and it was due to libnfc not scanning for I2C devices out of the box.

Let libnfc know the device address of the reader in `/etc/nfc/libnfc.conf`:

    device.name = "PN532 over I2C"
    device.connstring = "pn532_i2c:/dev/i2c-1"

List connected NFC readers:

    pi@raspberrypi:~ $ nfc-scan-device -v
    nfc-scan-device uses libnfc 1.7.1
    1 NFC device(s) found:
    - pn532_i2c:/dev/i2c-1:
        pn532_i2c:/dev/i2c-1
    chip: PN532 v1.6
    initator mode modulations: ISO/IEC 14443A (106 kbps), FeliCa (424 kbps, 212 kbps), ISO/IEC 14443-4B (106 kbps), Innovision Jewel (106 kbps), D.E.P. (424 kbps, 212 kbps, 106 kbps)
    target mode modulations: ISO/IEC 14443A (106 kbps), FeliCa (424 kbps, 212 kbps), D.E.P. (424 kbps, 212 kbps, 106 kbps)

Read a card or tag by first starting `nfc-poll` then physically holding a NFC/RFID tag or card in front of the reader:

    pi@raspberrypi:~ $ nfc-poll 
    nfc-poll uses libnfc 1.7.1
    NFC reader: pn532_i2c:/dev/i2c-1 opened
    NFC device will poll during 30000 ms (20 pollings of 300 ms for 5 modulations)
    ISO/IEC 14443A (106 kbps) target:
        ATQA (SENS_RES): 00  04  
           UID (NFCID1): 65  f7  0a  ab  
          SAK (SEL_RES): 08  
    nfc_initiator_target_is_present: Target Released
    Waiting for card removing...done.

Now I am free to use the device using `libnfc`.

## Troubleshooting 

### i2cdetect can't find the NFC device

Make sure the wiring is correct. Maybe SDA and SCL have been switched. Check the [Pi wiring diagram][GPIO map].

### libnfc can't find the NFC device

If it is already appearing in `i2cdetect`, make sure you have properly set the I2C address in `libnfc.conf`

## References
- http://www.byteparadigm.com/applications/introduction-to-i2c-and-spi-protocols/
- https://pinout.xyz/pinout/i2c
- https://learn.sparkfun.com/tutorials/raspberry-pi-spi-and-i2c-tutorial
- https://www.element14.com/community/community/raspberry-pi/blog/2012/12/14/nfc-on-raspberrypi-with-pn532-py532lib-and-i2c
- https://www.scribd.com/document/374277518/PN532-Manual-V3-pdf

[raspbian]: https://www.raspberrypi.org/downloads/raspbian/
[GPIO map]: https://pinout.xyz/pinout/i2c