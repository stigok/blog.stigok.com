---
layout: post
title:  "Booting with systemd-boot hangs on black screen with white caret"
date:   2019-02-21 18:22:28 +0100
categories: systemd
---

> **TL;DR**: a black screen with a blank caret with systemd-boot is a symptom of a malformed mkinitcpio.conf file, resulting in an incomplete initramfs image.

I was updating my */etc/mkinitcpio.conf* to enable resume from hibernation, when I noticed
there was a .pacnew version of the config file beside it. I went ahead and hand-edited the
new changes, added `resume` to the list of `HOOKS`, closed the file and ran
`mkinitcpio -p linux`.

Things seemed nice and dandy until I updated my kernel again. After a reboot, directly after
UEFI POST sequence, I was dumped to a black screen with a single white caret, hanging indefinitely.
No keyboard combinations worked, and even waiting two hours did nothing.


![systemd-boot black screen with white caret](/assets/img/IMG_20190221_174722-small.jpg)

I booted the system up again with the Arch Linux installation medium and mounted my main disk.
I then ran a diff of the two mkinitcpio.conf versions to try and find the errors.
It appeared that in both my `MODULES` and my `HOOKS` section, I had surrounded the items in both
parentheses *and* a single pair of double quotes.

```bash
MODULES=("intel_agp i915")
HOOKS=("base udev keyboard keymap consolefont autodetect modconf block efiverify encrypt lvm2 resume decryption-keys filesystems fsck")
```

The main upstream package change of the mkinitcpio.conf file was to switch from a single set of
double quotes to signify lists, into using bash array notation.

```bash
# Old format
list="a b c d e f g"

# New format (bash array notation)
list=(a b c d e f g)

# What I did wrong in my mkinitcpio.conf file
list=("a b c d e f g")
```

The last example resulted in the `list` variable containing an array with a single item,
hence early boot process unable to load required modules and run initial hooks.

I suspect I must've missed some warnings from pacman hooks after the system upgrade when
it generated the new initramfs images.

But the conclusion is that a black screen with a blank caret with `systemd-boot` is a symptom
of a malformed *mkinitcpio.conf* file, resulting in an incomplete initramfs image.

## References
- man mkinitcpio.conf
