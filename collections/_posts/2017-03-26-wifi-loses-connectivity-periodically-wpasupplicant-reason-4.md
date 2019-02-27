---
layout: post
title: "WiFi loses connectivity periodically wpa_supplicant reason 4"
date: 2017-03-26 12:28:30 +0200
categories: wifi linux networkmanager realtek
redirect_from:
  - /post/wifi-loses-connectivity-periodically-wpasupplicant-reason-4
---

> These are some issues I initially had with my Lenovo 100S with a Realtek 8821ae wifi adapter

I have been losing my wifi connection a lot lately with wpa_supplicant exiting with

    wpa_supplicant[474]: wlp2s0: CTRL-EVENT-DISCONNECTED bssid=xx:xx:xx:xx:xx reason=4 locally_generated=1

along with other random disconnects. It became more and more frequent by the day to the point where the connection would be dropped after just seconds of connectivity.


I found a list of [wpa_supplicant reason codes]() that claims reason 4 means *Disassociated due to inactivity*. The connection is certainly not idle when I lose my connectivity, so maybe some connectivity check is failing.

So I enabled debugging in NetworkManager by writing to `/etc/NetworkManager/conf.d/00-defaults.conf`

    [logging]
    level=DEBUG

and followed along whenever I lost connectivity `journalctl -u NetworkManager -f`.

## Realtek network card

After a lot of digging, I found a post explaining that Realtek adapters are know for losing connectivity, which seemed like an instant match for my troubles.
It also pointed in the direction of a GitHub repo containing sources of the latest updates to Realtek adapter drivers.

    $ git clone https://github.com/lwfinger/rtlwifi-new ~/src/rtlwifi-new
    $ cd ~/src/rtlwifi-new

But when ready to run `make` I did not have all the dependencies

    $ make
    make -C /lib/modules/4.10.5-1-ARCH/build M=/home/noop/src/rtlwifi_new modules
    make[1]: Entering directory '/usr/lib/modules/4.10.5-1-ARCH/build'
    make[1]: *** No rule to make target 'modules'.  Stop.
    make[1]: Leaving directory '/usr/lib/modules/4.10.5-1-ARCH/build'
    make: *** [Makefile:58: all] Error 2

In my case, on Arch Linux 4.10.5, I needed the `linux-headers` package

    # pacman -S linux-headers

Now it looks better

    $: make -C /lib/modules/4.10.5-1-ARCH/build M=/home/noop/src/rtlwifi_new modules
    make[1]: Entering directory '/usr/lib/modules/4.10.5-1-ARCH/build'
      Building modules, stage 2.
      MODPOST 15 modules
    make[1]: Leaving directory '/usr/lib/modules/4.10.5-1-ARCH/build'
    Making backups
    Install rtlwifi SUCCESS

Then install

    # make install

Enable module

    # modprobe rtl8821ae

And restart the machine

## Connectivity still lags

The driver certainly fixed the connection drops, but now a different issue arised. At specific intervals I'm temporarily losing connectivity.
The interfaces are still up and running, but I am getting random ping spikes of > 1000 ms in a time frame of about 10 seconds before returning to normal.
During the spikes, the debug logs are giving me some valuable hints:

    Mar 27 10:48:59 NetworkManager[419]: <debug> [1490604539.7178] device[0xfa8c50] (wlp2s0): wifi-scan: scanning requested
    Mar 27 10:48:59 NetworkManager[419]: <debug> [1490604539.7180] device[0xfa8c50] (wlp2s0): wifi-scan: no SSIDs to probe scan
    Mar 27 10:48:59 NetworkManager[419]: <debug> [1490604539.7193] device[0xfa8c50] (wlp2s0): add_pending_action (1): 'wifi-scan'
    Mar 27 10:48:59 NetworkManager[419]: <debug> [1490604539.7195] device[0xfa8c50] (wlp2s0): wifi-scan: scheduled in 120 seconds (interval now 120 seconds)
    Mar 27 10:48:59 NetworkManager[419]: <debug> [1490604539.7592] device[0xfa8c50] (wlp2s0): wifi-scan: scanning-state: scanning
    Mar 27 10:49:07 NetworkManager[419]: <debug> [1490604547.8645] ndisc-lndp[0x10fd0e0,"wlp2s0"]: processing libndp events
    Mar 27 10:49:08 NetworkManager[419]: <debug> [1490604548.7870] ndisc-lndp[0x10fd0e0,"wlp2s0"]: processing libndp events
    Mar 27 10:49:10 NetworkManager[419]: <debug> [1490604550.6191] ndisc-lndp[0x10fd0e0,"wlp2s0"]: processing libndp events
    Mar 27 10:49:13 NetworkManager[419]: <debug> [1490604553.1484] device[0xfa8c50] (wlp2s0): wifi-scan: scan-done callback: successful

While debugging the previous issue, i stumbled upon an article talking about [lost connections when wifi scanning was taking place](https://blogs.gnome.org/dcbw/2016/05/16/networkmanager-and-wifi-scans/),
which seems to be closely related to the issues I'm having now.
So I'll bind the BSSID on the networks I'm having issues on, and hoping for the best

## Other things I tried
### Start from scratch

After changing laptop a while back, I used unedited connection profiles from `/etc/NetworkManager/system-connections`.
This may have been causing issues with invalid connection UIDS and device hardware addresses.
I just deleted all my connections and I'm re-creating them from scratch.
It did not make a difference at first, but may be an important factor later on.

### Verify that NetworkManager connectivity check works

Then I run `journalctl -u NetworkManager -f` before restarting NetworkManager with `systemctl restart NetworkManager` and start watching the logs. The first thing I see is what configuration files are being read:

    <info>  [1490522215.5734] Read config: /etc/NetworkManager/NetworkManager.conf (lib: 20-connectivity.conf) (etc: 00-defaults.conf, 10-powersaving.conf

I know the settings of defaults and powersaving, but the `lib: 20-connectivity.conf` seems interesting. `/usr/lib/NetworkManager/conf.d/20-connectivity.conf` contains the following:

    [connectivity]
    uri=http://www.archlinux.org/check_network_status.txt

But later I find the logs says its fine:

    NetworkManager[30021]: <debug> [1490522524.6651] connectivity: check: send periodic request to 'http://www.archlinux.org/check_network_status.txt'
    NetworkManager[30021]: <debug> [1490522524.7821] connectivity: check for uri 'http://www.archlinux.org/check_network_status.txt' successful.

... which I can confirm by reading the CONNECTIVITY SECTION of `man NetworkManager.conf`

    $ curl http://www.archlinux.org/check_network_status.txt
    NetworkManager is online

For those interested, I have [set up my own connectivity check endpoint](https://blog.stigok.com/post/roll-your-own-networkmanager-connectivity-check-endpoint-with-nginx).


## References
- [wpa_supplicant reason codes]()

[wpa_supplicant reason codes]: http://www.aboutcher.co.uk/2012/07/linux-wifi-deauthenticated-reason-codes/