---
layout: post
title: "Persistently disable wlan powersaving in NetworkManager"
date: 2017-03-15 14:52:51 +0100
categories: NetworkManager linux networking
redirect_from:
  - /post/persistently-disable-wlan-powersaving-in-networkmanager
---

After each time I had rebooted my machine I would have to manually disable power management on my wireless network card. I was using `iwconfig` for this, but then it occured to me that it could probably be done automatically in `NetworkManager` itself.

---

When I'm not connected to a network

    $ iwconfig wlp2s0
    wlp2s0    IEEE 802.11  ESSID:off/any  
              Mode:Managed  Access Point: Not-Associated   Tx-Power=20 dBm   
              Retry short limit:7   RTS thr=2347 B   Fragment thr:off
              Power Management:on

Manually disable power management

    # iwconfig wlp2s0 power off
    $ iwconfig wlp2s0
    wlp2s0    IEEE 802.11  ESSID:off/any  
              Mode:Managed  Access Point: Not-Associated   Tx-Power=20 dBm   
              Retry short limit:7   RTS thr=2347 B   Fragment thr:off
              Power Management:off

But instead of doing this manually; create a new drop-in configuration file in the NetworkManager `conf.d/` folder:

    #
    # /etc/NetworkManager/conf.d/10-powersaving.conf
    #
    [connection-wifi-wlp2s0]
    wifi.powersave=off

Reboot to see if it works as it should. Otherwise get some clues with `journalctl --boot` to read the log of the current boot.

Expected output after reboot

    $ iwconfig wlp2s0 | grep "Power Management"

# References

- `man NetworkManager`
- `man NetworkManager.conf`
- `man nm-settings`