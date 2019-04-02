---
layout: post
title: "Useful commands to Cisco 3560/3650E"
date: 2017-07-07 15:45:30 +0200
categories: cisco network
redirect_from:
  - /post/useful-commands-to-cisco-35603650e
---

![Cisco 3560E shell](https://public.stigok.com/img/1499247045993464200.png)

> All commands will be resolved to their longer equivalents as long as they are unambiguous. E.g. `show ru`
> will be enough to execute `show running-configuration`, but `show r` *would* be amigous with `show route`
> and will fail.

Connect with SSH (see bottom of post if you're running into "No matching key exchange method found")

    $ ssh admin@switch.domain.tld

> The rest of the commands in this section are executed from within the Cisco SSH shell.
> `$(conf)` signifies that the shell is in configure mode.

Dump current configuration

    $ show running-configuration

Enter configuration mode

    $ conf t

> When in configuration mode, a different subset of commands are available. To run commands normal commands
> when in this mode, prefix the command line with `do`. E.g. `$(conf) do show vlan`

Configure a SPI interface

    $(conf) interface vlan777
    $(conf-if) 

Configure a port

    $(conf) interface GigabitEthernet 0/3
    $(conf if) vlan 777
    $(conf if) mode access

Configure an ipv6 route. This command will forward the whole /48 subnet to a
single host at ::2121/64.

    $(conf) ipv6 route 1234:5678:abcd::/48 1234:5678::2121

Pipe to a `grep`-like filter. The below command will print all running configuration lines which contains the
string *ipv6*. (`inc` is short for `include`)

    $(conf) do show run | inc ipv6

End configuration mode and go back to normal

    $(conf) end
    $

List files in file memory

    $ dir

Open up a file for viewing. `more` is the default pager, but reads named files too

    $ more config.txt

Save configuration to persistent memory (this is done outside of configuration mode)

    $ write mem

Show interface configuration (`int` is short for `interfaces`)

    $ show int

Show MAC address table

    $ show mac address-table

## Issues encountered
### No matching key exchange method found

When attempting to connect using SSH, a key error occurs, preventing connection

    $ ssh admin@switch.example.com
    Unable to negotiate with 10.0.23.16 port 22: no matching key exchange method found. Their offer: diffie-hellman-group1-sha1

Excplicitly pass a deprecated option to `ssh`

    $ ssh -o "KexAlgorithms diffie-hellman-group1-sha1" admin@switch.example.com
    (cisco)$

### No matching cipher found

When attempting to connect using SSH, a key error occurs, preventing connection

    $ ssh admin@switch.example.com
    Unable to negotiate with 10.0.23.16 port 22: no matching cipher found. Their offer: aes128-cbc,3des-cbc,aes192-cbc,aes256-cbc

Excplicitly pass a deprecated cipher from the list above to `ssh`

    $ ssh -o "Ciphers aes128-cbc,3des-cbc,aes192-cbc,aes256-cbc" admin@switch.example.com
    (cisco)$
