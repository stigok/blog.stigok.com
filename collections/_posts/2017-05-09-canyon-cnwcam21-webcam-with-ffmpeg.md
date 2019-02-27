---
layout: post
title: "Canyon CN-WCAM21 webcam with ffmpeg"
date: 2017-05-09 21:30:01 +0200
categories: ffmpeg linux video
redirect_from:
  - /post/canyon-cnwcam21-webcam-with-ffmpeg
---

Trying to take a screenshot with my webcam using ffmpeg. Using ffmpeg is making it possible to turn the camera on, wait for brightness/gamma correction, then take a snapshot when I deem it to be ready. For me I assume 5 seconds should be ready.

My initial try which worked fine with my built-in webcam, threw an error at me when trying with my (old!!!) CN-WCAM21

    $ ffmpeg -f v4l2 -i /dev/video1 -f image2 /tmp/screen.jpeg                                        
    ffmpeg version 3.3 Copyright (c) 2000-2017 the FFmpeg developers
      built with gcc 6.3.1 (GCC) 20170306
      configuration: --prefix=/usr --disable-debug --disable-static --disable-stripping --enable-avisynth --enable-avresample --enable-fontconfig --enable-gmp --enable-gnutls --enable-gpl --enable-ladspa --enable-libass --enable-libbluray --enable-libfreetype --enable-libfribidi --enable-libgsm --enable-libiec61883 --enable-libmodplug --enable-libmp3lame --enable-libopencore_amrnb --enable-libopencore_amrwb --enable-libopenjpeg --enable-libopus --enable-libpulse --enable-libschroedinger --enable-libsoxr --enable-libspeex --enable-libssh --enable-libtheora --enable-libv4l2 --enable-libvidstab --enable-libvorbis --enable-libvpx --enable-libwebp --enable-libx264 --enable-libx265 --enable-libxcb --enable-libxvid --enable-netcdf --enable-shared --enable-version3
      libavutil      55. 58.100 / 55. 58.100
      libavcodec     57. 89.100 / 57. 89.100
      libavformat    57. 71.100 / 57. 71.100
      libavdevice    57.  6.100 / 57.  6.100
      libavfilter     6. 82.100 /  6. 82.100
      libavresample   3.  5.  0 /  3.  5.  0
      libswscale      4.  6.100 /  4.  6.100
      libswresample   2.  7.100 /  2.  7.100
      libpostproc    54.  5.100 / 54.  5.100
    [video4linux2,v4l2 @ 0x55d662616560] Cannot find a proper format for codec 'none' (id 0), pixel format 'none' (id -1)
    Assertion *codec_id != AV_CODEC_ID_NONE failed at libavdevice/v4l2.c:808

Got it to work by setting LD_PRELOAD environment variable loading video4linux1 driver

    LD_PRELOAD=/usr/lib/libv4l/v4l1compat.so ffmpeg -f v4l2 -t 6 -i /dev/video1 -r 1 -f image2 /tmp/screen%03d.jpeg

The above command starts the video stream, takes a screenshot (`-f image2`) every second (`-r 1`), then stops after 6 seconds (`-t`)

## References

- <https://wiki.archlinux.org/index.php/webcam_setup#Get_software_to_use_your_webcam>