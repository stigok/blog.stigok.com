---
layout: post
title:  "Custom keyboard layout in Wayland like Xmodmap in X11"
date:   2020-10-27 15:42:50 +0100
categories: x11 wayland gnome
excerpt: Moving from Xmodmap in X11 to xkb in Wayland
---

## Preface

Earlier I used *.Xmodmap* to define my custom keyboard layout for Norwegian
characters behind Caps-Lock. It was surprisingly easy, as described [in an
earlier blog post][1]. ([gist][2])

However, moving to Wayland, this no longer works. It was actually
surprisingly hard to get right.

As of writing I'm on GDM (Gnome Display Manager) in gnome-shell on Arch Linux,
but this *may* apply to other Display Managers under Wayland as well.

## Configuring a custom keyboard layout in Wayland

My old *.Xmodmap*, it looked like this:

```
clear lock
!Maps Caps-Lock as Level3 Shift
keycode 66 = Mode_switch ISO_Level3_Shift

!Norwegian chars ÆØÅ
keycode 47 = semicolon colon oslash Oslash
keycode 48 = apostrophe quotedbl ae AE
keycode 34 = bracketleft braceleft aring Aring
```

This gives me Norwegian characters when I hold the Caps-Lock key and press the
bounded keys. The keycodes used above can be identified using `xev`.

For Wayland it has to be split up into two portions. One new file, containing
the key bindings goes into */usr/share/X11/xkb/symbols/<layout_name>*.

To see how other layouts are configured, you can look through the files in
the same directory. I am using the *English (US, euro on 5)* as my reference
since it is the base layout I'm interested in overriding.

The contents of my *stigok* keyboard layout looks like this:

```
default partial alphanumeric_keys

// Overriding the existing us(euro) symbols
// May not be necessary, but it works
xkb_symbols "euro" {

    // Contents of original us(euro)
    include "us(basic)"
    name[Group1]= "STIGOK US English (Norwegian chars behind Caps-Lock)";

    include "eurosign(5)"
    include "level3(ralt_switch)"
    // End contents of original us(euro)

    // Norwegian keys hiding in Level 3
    key <AC10> { [ semicolon, colon, oslash, Oslash ] };
    key <AC11> { [ apostrophe, quotedbl, ae, AE ] };
    key <AD11> { [ bracketleft, braceleft, aring, Aring ] };

    // Toggle Level 3 with Caps-Lock
    key <CAPS> { [ ISO_Level3_Shift ] };

};
```

Now, to get this visible for the Display Manager, I am adding a new node to the
*/usr/share/X11/xkb/rules/evdev.xml* file under the Xpath */xkbConfigRegistry/layoutList*:

```xml
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE xkbConfigRegistry SYSTEM "xkb.dtd">
<xkbConfigRegistry version="1.1">
  <modelList>
    <!--[redacted]-->
  </modelList>
  <layoutList>
    <layout>
      <configItem>
        <name>stigok</name>
        <shortDescription>en</shortDescription>
        <description>STIGOK English (US)</description>
        <languageList>
          <iso639Id>eng</iso639Id>
        </languageList>
      </configItem>
      <!--[redacted]-->
    </layout>
<!--[redacted]-->
```

I.e. the only portion you need to add is the `<layout>...</layout>` under `<layoutList>`.
- The value for `<name>` is the name of the symbols file we created earlier.
- The `<description>` is what shows up in the Settings dialog in gnome-shell
- The other values are mirrored from other similar configurations

Now you have to restart Gnome. In Wayland it can't be done using *Alt+F2 -> R*,
so I'm restarting by either killing all processes named something with Wayland,
or I restart the GDM service `systemctl restart gdm.service`.

Now you should have something like the below images.

![Adding new keyboard layout GDM gnome-shell searching by name](https://public.stigok.com/img/2020-10-27-161221.png)

![GDM gnome-shell settings custom keyboard layout](https://public.stigok.com/img/2020-10-27-161031.png)

## Improvements

### Putting the configuration in home directory dotfiles

I tried to configure this in my *$HOME/.config/xkb* instead, but I was unable
to get it working properly. The first hurdle was to make sure the `$XDG_CONFIG_HOME`
was actually set to `~/.config`. That was solved by setting the variable in
*/etc/security/pam_env.conf*, since it's no longer set automatically in Arch Linux.

```bash
XDG_CONFIG_HOME DEFAULT=${HOME}/.config
```

I think I have to import complete layouts there to get it working. Something
for later, I guess.

**Update:** as the comments in the  *pam_env.conf* file says, `$HOME` might not
always be set for the calling applications. This seemed to set the path for `xkb`
correctly, but not for `lpass` (of all things), which had the `XDF_CONFIG_HOME`
set as `/.config/lpass` instead of `~/.config/lpass`. This need further investigation
to get working properly.

## References

- Thanks to [Håvard Moen](https://github.com/umglurf/) for helping out with configuration
- <https://askubuntu.com/a/100228/20164>
- <https://superuser.com/a/848050/5075>
- <https://askubuntu.com/a/1274619/20164>

[1]: https://blog.stigok.com/2017/07/10/norwegian-keys-on-us-keyboard-layout-with-xmodmap.html
[2]: https://gist.github.com/stigok/2aabc44c3a151fc3d3fed119f2c1f44f/
