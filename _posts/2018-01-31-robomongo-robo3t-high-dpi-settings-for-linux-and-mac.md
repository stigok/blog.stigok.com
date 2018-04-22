---
layout: post
title: "Robomongo (Robo-3T) high DPI settings for Linux and Mac"
date: 2018-01-31 15:13:10 +0100
categories: qt highdpi
redirect_from:
  - /post/robomongo-robo3t-high-dpi-settings-for-linux-and-mac
---

My monitor is high DPI, and when starting robomongo, everything was a bit too small.

Downloading an unzipping gave me this tree

    robo3t-1.1.1-linux-x86_64-c93c6b0 $ tree
    .
    ├── bin
    │   ├── qt.conf
    │   └── robo3t
    ├── lib
    │   ├── [...]
    ├── CHANGELOG
    ├── COPYRIGHT
    ├── DESCRIPTION
    └── LICENSE

Which gave a hint this application is using Qt. Looking in the [git repo](https://github.com/Studio3T/robomongo/tree/v1.1.1), I found that they are using Qt5, and [they support HIGH_DPI mode](https://doc.qt.io/qt-5/highdpi.html).

![Normal mode](https://public.stigok.com/img/1517407987621419488.png)

Before starting the binary, I can set the environment variable 'QT_SCALE_FACTOR=2' to be able to get a usable user interface.

    $ QT_SCALE_FACTOR=2 ./bin/robo3t

![Scaled DPI](https://public.stigok.com/img/1517408005510610770.png)