---
layout: post
title: "Polybar jsoncpp cannot open shared object file"
date: 2018-02-09 13:27:27 +0100
categories: arch pacman polybar
redirect_from:
  - /post/polybar-jsoncpp-cannot-open-shared-object-file
---

After an update of jsoncpp (1.8.4-1 -> 1.8.4-2), the library object file path that polybar looks for has changed filename.

    polybar: error while loading shared libraries: libjsoncpp.so.19: cannot open shared object file: No such file or directory

After the update, the filename version suffix was bumped up, chaning it to `libjsoncpp.so.20`. This path is set when polybar is compiled, so re-compiling will make it work properly again.

I am using `yaourt` to handle my AUR packages, so I can trigger a compile and reinstall with

    # yaourt -S polybar

If you think this is too much of an hassle to do, you can also just symlink the new file back to the old filename.

    # ln -s /usr/lib/libjsoncpp.so.20 /usr/lib/libjsoncpp.so.19

There's [an issue filed in polybar upstream](https://github.com/jaagr/polybar/issues/987) that discusses this issue in more detail.

Additionally, an issue regarding [inconsistencies in .so file versioning](https://github.com/open-source-parsers/jsoncpp/issues/734) during the build process of jsoncpp itself may be related.

## References
- https://github.com/jaagr/polybar/issues/987
- https://www.archlinux.org/packages/extra/x86_64/jsoncpp/
- https://github.com/open-source-parsers/jsoncpp/issues/734