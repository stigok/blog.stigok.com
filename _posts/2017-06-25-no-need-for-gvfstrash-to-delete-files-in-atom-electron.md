---
layout: post
title: "No need for gvfs-trash to delete files in Atom (Electron)"
date: 2017-06-25 01:05:49 +0200
categories: atom electron trash
redirect_from:
  - /post/no-need-for-gvfstrash-to-delete-files-in-atom-electron
---

I don't want to install `gvfs` with all its friends just to be able to delete files from within Atom/Electron. However, the environment variable `ELECTRON_TRASH` can set another trash handler, e.g. `trash-cli`, which is dependency free.

So I append the variable to my `.zshenv` and live on as a happy man.

    # pacman -S trash-cli --noconfirm > /dev/null
    $ echo "export ELECTRON_TRASH=trash-cli" >> ~/.zshenv

Restart Atom and continue to enjoy life

## References

- https://wiki.archlinux.org/index.php/Atom#Unable_to_delete_files
- https://github.com/electron/electron/pull/7178