---
layout: post
title:  "Running WoW 1.12 on Arch Linux with wine and nvidia with a NVIDIA GeForce 9800 GTX"
date:   2019-04-30 23:15:00 +0200
categories: arch linux wine nvidia gaming
---

I wanted to run the WoW 1.12 on Arch Linux using Wine with the NVIDIA graphics driver instead of nouveau.

My first obstacle was getting X to start at all with my card. After installing the `nvidia-dkms` drivers,
X would die on start. The logs at */var/log/Xorg.0.log* read that the current driver did not support my card.
A [news entry](https://www.archlinux.org/news/nvidia-340xx-and-nvidia/) on Arch Linux website states the following:

> As NVIDIA dropped support for G8x, G9x, and GT2xx GPUs with the release of 343.22, there now is set of nvidia-340xx packages supporting those older GPUs. 340xx will receive support until the end of 2019 according to NVIDIA.
> Users of older GPUs should consider switching to nvidia-340xx. The nvidia-343.22 and nvidia-340xx-340.46 packages will be in testing for a few days.

As I have a GeForce 9800 GTX, I now need the 340xx drivers

```
$ aurman -S nvidia-340xx-dkms lib32-nvidia-340xx-utils
```

Then I ran wow successfully again

```
$ wine ~/games/wow/WoW.exe
```

## Aftermath

It works, however, I hoped it would save me some CPU time, but it didn't. I guess wine is doing a lot of work translating opengl/directx or something similar...(?)
The graphics is better, though!

## Troubleshooting

### Wine libGL error

When running Wine, it exists with a similar error message

```
libGL error: No matching fbConfigs or visuals found
libGL error: failed to load driver: swrast
X Error of failed request:  GLXBadContext
  Major opcode of failed request:  154 (GLX)
  Minor opcode of failed request:  6 (X_GLXIsDirect)
  Serial number of failed request:  173
  Current serial number in output stream:  172
  ```

  Make sure you have the lib32 utils installed as well. For the 340xx drivers,`lib32-nvidia-340xx-utils` is the package you want.

## References
- https://bbs.archlinux.org/viewtopic.php?id=139071
- https://askubuntu.com/questions/541343/problems-with-libgl-fbconfigs-swrast-through-each-update
- https://bbs.archlinux.org/viewtopic.php?id=231318
