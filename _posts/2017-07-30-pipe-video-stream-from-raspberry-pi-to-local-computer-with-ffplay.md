---
layout: post
title: "Pipe video stream from raspberry pi to local computer with ffplay"
date: 2017-07-30 00:32:09 +0200
categories: arch linux arm raspberrypi ffmpeg ffplay video
redirect_from:
  - /post/pipe-video-stream-from-raspberry-pi-to-local-computer-with-ffplay
---

I use this to get a live video stream from my Raspberry Pi with Camera attached

Execute this on the Pi, where `TARGET_IP` is my local computer where I will watch the stream, and `PORT` is an arbitrary port number.

    $ raspivid -t 999999 -o - | nc -u $TARGET_IP $PORT

Execute this on the local computer where you will watch the video stream

    $ nc -ul $PORT | ffplay -

## References

- `man ffplay`
- `man nc`
- https://www.raspberrypi.org/blog/camera-board-available-for-sale/
- https://blog.philippklaus.de/2013/06/using-the-raspberry-pi-camera-board-on-arch-linux-arm/