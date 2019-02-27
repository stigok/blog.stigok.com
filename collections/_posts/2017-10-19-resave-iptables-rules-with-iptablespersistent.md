---
layout: post
title: "Re-save iptables rules with iptables-persistent"
date: 2017-10-19 01:28:04 +0200
categories: iptables ubuntu debian apt dpkg
redirect_from:
  - /post/resave-iptables-rules-with-iptablespersistent
---

When installing `iptables-persistent` with `apt`

    # apt install iptables-persistent

the shell takes me into a wizard which asks me if I want to save the current IPv4 and IPv6 rules

![iptables-persistent wizard step 1](https://public.stigok.com/img/1508369020334500434.png)

So then, later, when I've made some more changes to my iptables ruleset, instead of manually putting the rulesets into files with `iptables-save` and `ip6tables-save`, I can re-run the wizard `iptables-persistent` which does it all for me:

    # dpkg-reconfigure iptables-persistent