---
layout: post
title:  "Adjust auto close hours for Discourse post topics"
date:   2019-04-14 12:09:44 +0200
categories: discourse
---

**TL;DR:** go into settings for a specific topic, and set *Auto close hours*.

I am administering a Discourse forum where old posts were getting automatically
closed after 30 days. I could not find out where these settings were hidden.
In the forum settings I was searching for a diversity of settings by keyword
without any luck: *old, cold, auto, close, lock, 30*.

Eventually, skimming through the logs, I could see that during setup, another
admin had set these settings on a per-topic basis. So this is in fact a setting
applied to individual topics, not to the forum globally. I guess it makes
sense when I see it now, it was just so hard for me to find.
