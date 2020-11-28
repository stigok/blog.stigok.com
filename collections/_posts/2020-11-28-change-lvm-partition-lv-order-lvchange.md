---
layout: post
title:  "Change partition or logical volume (LV) order in LVM"
date:   2020-11-28 14:41:47 +0100
categories: lvm linux
excerpt: Changing the order of logical volumes inside LVM using lvchange
#proccessors: pymd
---

## Preface

After creating logical volumes in LVM I noticed I had created them in the wrong order.
If it makes a difference or not is another question, but I wanted my *root* volume to appear
before my *swap* volume.

This can be done by chaning the minor decice number of the logical volume. Lower numbers
comes before higher ones, hence, I want to make *root* a lower number than *swap*.

## Changing minor device numbers of logical volumes in LVM

The manpage (`man lvchange`) describes the options well enough with the
description

> Make the minor device number persistent for an LV. [...]

Let's take a look at the layout of my volume group (VG), **Main**:

```
$ lsblk /dev/nvme0n1p3
NAME          MAJ:MIN RM   SIZE RO TYPE  MOUNTPOINT
nvme0n1p3     259:3    0 475.4G  0 part
└─lvm         254:0    0 475.4G  0 crypt
  ├─Main-swap 254:1    0    16G  0 lvm   [SWAP]
  ├─Main-root 254:2    0   100G  0 lvm   /
  └─Main-home 254:3    0   310G  0 lvm   /home
```

Here we can see that the *major* device number for the VG is 254, while the
minor numbers increment by one for each of the LV's.

I want `Main-root` to have number 1, and swap to be number 2. To do that, I first
have to move `Main-swap` to a free number, e.g. **7**.
If any of the volumes are currently active, you will be asked if it's okay to deactivate them temporarily.

```
# lvchange --persistent y --minor 7 /dev/Main/swap
```

Then set `Main-root` to number 1 now that it should be available

```
# lvchange --persistent y --minor 1 /dev/Main/root
```

Set `Main-swap` again to number 2

```
# lvchange --persistent y --minor 2 /dev/Main/swap
```

Now, since the volumes were deactivated, activate them all again

```
# lvchange --activate ay /dev/Main
```

`ay` here means the following:

> ay specifies autoactivation, in which case an LV
> is activated only if it matches an item in lvm.conf
> activation/auto_activation_volume_list.  If the list is not set, all
> LVs are considered to match, and if if the list is set but
> empty, no LVs match

If this doesn't activate your volumes, you may want `y` instead.

## References
- `man lvchange`
