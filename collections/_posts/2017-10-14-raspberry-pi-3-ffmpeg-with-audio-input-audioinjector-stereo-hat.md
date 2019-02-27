---
layout: post
title: "Raspberry Pi 3 ffmpeg with audio input audioinjector stereo hat"
date: 2017-10-14 23:08:12 +0200
categories: draft raspberrypi audioinjector audio ffmpeg alsa
redirect_from:
  - /post/raspberry-pi-3-ffmpeg-with-audio-input-audioinjector-stereo-hat
---

![Audio Injector Raspberry Pi hat](http://www.audioinjector.net/images/audioinjector.pi.soundcard.jpg)

I'm using this with `ffmpeg` in order to stream an audio source to a RTMP endpoint.

Install the hat physically on the Pi (picture coming)

## Install

> Follow [installation instructions](http://forum.audioinjector.net/viewtopic.php?f=5&t=3#p3) on the official forum, or just go along with mine

Open up alsamixer and enable capture mode on the input devices by selecting the device and pressing the spacebar
![alsamixer enable capture mode](https://public.stigok.com/img/1508014360377760233.png)

Change the input muxer in alsamixer to line-in in order to capture audio from the phono jacks instead of the on-board mic. Do this by selecting the device, then press arrow up/down. Look in the upper left corner to see what's currently selected.
![alsamixer select input mux device](https://public.stigok.com/img/1508013757462968285.png)

## References
- http://forum.audioinjector.net/viewtopic.php?f=5&t=3
- https://trac.ffmpeg.org/wiki/Capture/ALSA
- https://superuser.com/questions/253467/convert-an-mp3-from-48000-to-44100-hz