---
layout: post
title: "gpg refresh keys failed"
date: 2017-03-07 16:59:18 +0100
categories: gpg dirmngr linux
redirect_from:
  - /post/gpg-refresh-keys-failed
---

Trying to `--refresh-keys` my gpg key-chain, but gpg doesn't let me.

    $ gpg --refresh-keys
    gpg: refreshing 42 keys from hkp://keys.gnupg.net
    gpg: keyserver refresh failed: No keyserver available

## Things tried

- The [Arch Linux Wiki page on gpg](https://wiki.archlinux.org/index.php/GnuPG#gpg_hanged_for_all_keyservers_.28when_trying_to_receive_keys.29) mentioned that killing `dirmngr` before trying again might work out, but to no avail.
- Upgrading all packages `sudo pacman -Syu` did not help
- Rebooting did not help
- Specifying a keyserver manually `gpg --refresh-keys --keyserver keyserver.ubuntu.com:80` doesn't change a thing

## Maybe my firewall log can give me some hints.

    $ gpg --refresh-keys 2> /dev/null
    # journalctl --since "30 seconds ago"
    Mar 01 14:33:55 user kernel: [UFW AUDIT] IN= OUT=lo SRC=127.0.0.1 DST=127.0.0.1 LEN=60 TOS=0x00 PREC=0x00 TTL=64 ID=37191 DF PROTO=TCP SPT=50354 DPT=9050 WINDOW=43690 RES=0x00 SYN URGP=0
    Mar 01 14:33:55 user kernel: [UFW AUDIT] IN= OUT=lo SRC=127.0.0.1 DST=127.0.0.1 LEN=60 TOS=0x00 PREC=0x00 TTL=64 ID=37186 DF PROTO=TCP SPT=35828 DPT=9150 WINDOW=43690 RES=0x00 SYN URGP=0

Every time I run the command, something tries to connect to _127.0.0.1_ on port _9050_, then _9150_. Looking up those ports on Wikipedia tells me these two ports are used by Tor. `~/.gnupg/dirmngr.conf` specifies a .onion address as a `keyserver`, so this makes sense now. Moving on!

    # My ~/.gnupg/dirmngr.conf
    keyserver hkp://jirk5u4osbsr34t5.onion
    keyserver hkp://keys.gnupg.net

## Debug the dirmngr program

Killing off existing instances `pkill dirmngr`. Then starting `gpg` with `--debug` to get some verbose output.

    $ gpg --debug 1024 --refresh
    gpg: reading options from '/home/user/.gnupg/gpg.conf'
    gpg: enabled debug flags: ipc
    gpg: DBG: chan_3 <- # Home: /home/user/.gnupg
    gpg: DBG: chan_3 <- # Config: /home/user/.gnupg/dirmngr.conf
    gpg: DBG: chan_3 <- OK Dirmngr 2.1.19 at your service
    gpg: DBG: connection to the dirmngr established
    gpg: DBG: chan_3 -> GETINFO version
    gpg: DBG: chan_3 <- D 2.1.19
    gpg: DBG: chan_3 <- OK
    gpg: DBG: chan_3 -> KEYSERVER
    gpg: DBG: chan_3 <- S KEYSERVER hkp://keys.gnupg.net
    gpg: DBG: chan_3 <- OK
    gpg: refreshing 11 keys from hkp://keys.gnupg.net
    gpg: DBG: chan_3 -> KS_GET -- <a list of keys>
    gpg: DBG: chan_3 <- ERR 167772379 Server indicated a failure <Dirmngr>
    gpg: keyserver refresh failed: Server indicated a failure
    gpg: DBG: chan_3 -> BYE
    gpg: secmem usage: 0/32768 bytes in 0 blocks

After searching around for the number after `ERR 167772379` and _gpg_, I came to [a bug posted to Arch Linux][arch] describing
a proposed work-around to start `dirmngr` with `--standard-resolver`.

> **--standard-resolver** This  option  forces the use of the system's standard DNS resolver code.

    echo "standard-resolver" >> ~/.gnupg/dirmngr.conf

Kill the `dirmngr` process if it's running to make sure new config is used

    $ pkill dirmngr

Now retrying

    $ pkill dirmngr
    $ gpg --refresh-keys
    gpg: refreshing 42 keys from hkp://keys.gnupg.net
    [output redacted]
    gpg: Total number processed: 42
    gpg:              unchanged: 21

Success!

## Thoughts

I suspect this problem may be due to [my strict iptables] rules as well. Will maybe get back to this in the future with a different fix.

## Resources

  - [Unable to connect to keyservers for any action with 2.1.17](https://bugs.archlinux.org/task/52234)
  - [dirmngr fails when reverse DNS lookups do not work](https://bugs.debian.org/cgi-bin/bugreport.cgi?bug=854359)