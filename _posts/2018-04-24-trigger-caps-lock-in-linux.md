---
layout: post
title:  "Trigger Caps Lock in Linux from the terminal"
date:   2018-04-24 15:30:42 +0200
categories: linux x11
---

I have rebound my Caps Lock key using `xmodmap` so that it can be used as a modifier key to trigger Norwegian characters.

However, when I connect an external keyboard, the keymap changes back to default again and I have to re-run `xmodmap ~/.xmodmap` in order to reload the custom key bindings.
If I have enabled Caps Lock before I reload, I can no longer use that button to disable it again. This is where `xdotool` come into the picture.

## Send keystroke to X11

> xdotool - command-line X11 automation tool

One of the available commands is `key`, which will send a key to the current window.
I don't really care what window I send Caps Lock to, so I can just send it directly.

```terminal
$ xdotool key Caps_Lock
```

And now my Caps Lock is disabled again.

## References
- `man xdotool`
