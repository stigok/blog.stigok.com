---
layout: post
title: "Norwegian keys on US keyboard layout with .Xmodmap"
date: 2017-07-10 15:48:39 +0200
categories: keyboard-layout norwegian xmodmap xorg
redirect_from:
  - /post/norwegian-keys-on-us-keyboard-layout-with-xmodmap
---

Was first looking at how to configure using xkbcomp and setxkbmap, but just went with .Xmodmap for simplicity. The [Arch Linux wiki page on Xmodmap](https://wiki.archlinux.org/index.php/Xmodmap) states that it's not really the best way to go, but no way tends to be.

> xmodmap settings are reset by setxkbmap, which not only alters the alphanumeric keys to the values given in the map, but also resets all other keys to the startup default.

The following `.Xmodmap` configuration binds Caps Lock as the modifier key to trigger the Norwegian letters on a US keyboard layout.

    ! ~/.Xmodmap
    clear lock
    !Maps Caps-Lock as Level3 Shift
    keycode 66 = Mode_switch ISO_Level3_Shift
    !Norwegian alpha chars ÆØÅ
    keycode 47 = semicolon colon oslash Oslash
    keycode 48 = apostrophe quotedbl ae AE
    keycode 34 = bracketleft braceleft aring Aring

## Usage

[[Caps-Lock]] + [[\']] = æ  
[[Caps-Lock]] + [[;]] = ø  
[[Caps-Lock]] + [[&#91;]] = å  
[[Caps-Lock]] + [[Shift]] + [[\']] = Æ  
[[Caps-Lock]] + [[Shift]] + [[;]] = Ø  
[[Caps-Lock]] + [[Shift]] + [[&#91;]] = Å

## References
- https://wiki.archlinux.org/index.php/Xmodmap
- https://wiki.archlinux.org/index.php/Keyboard_configuration_in_Xorg
- https://github.com/gulrotkake/Ubuntu-MacBook-Norwegian-keyboard
- http://wiki.linuxquestions.org/wiki/List_of_Keysyms_Recognised_by_Xmodmap