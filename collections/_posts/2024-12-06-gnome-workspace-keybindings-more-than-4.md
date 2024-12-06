---
layout: post
title:  "Configure keybindings for more than 4 workspaces in Gnome"
date:   2024-12-06 09:12:35 +0100
categories: gnome, linux
excerpt: The settings UI in Gnome only supports configuring keybindings for up to four workspaces. I would like to configure at least five.
#proccessors: pymd
---

We can Use `dconf` to read and write settings for Gnome. Many things you can
tune in the native settings UI is actually persisted in the dconf database.

I have 5 workspaces, instead of the 4 standard ones. In the keyboard shortcut
UI, it is not possible to configure key bindings more than four. Thankfully,
we can use the `dconf` tool directly to manipulate the settings.

I want to use the <kbd>Windows Logo</kbd> key on my keyboard, combined with a number,
to switch the currently active workspace. I first went into the Gnome settings
to configure the first four, so that I know the format to use for the fifth.

```
$ dconf read /org/gnome/desktop/wm/keybindings/switch-to-workspace-4
['<Super>4']

$ dconf write /org/gnome/desktop/wm/keybindings/switch-to-workspace-5 "['<Super>5']"
```

I also want to use <kbd>Shift</kbd>+<kbd>Windows Logo</kbd>+<kbd>5</kbd> to
move the active window to the fifth workspace.

```
$ dconf read /org/gnome/desktop/wm/keybindings/move-to-workspace-5
['<Shift><Super>4']

$ dconf write /org/gnome/desktop/wm/keybindings/move-to-workspace-5 "['<Shift><Super>5']"
```

This was performed on a Fedora 40 workstation, using Gnome Shell v47.
But this should work for all of Gnome released within the past 5 years, at least.

## References
- `man dconf`
