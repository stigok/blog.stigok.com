---
layout: post
title: "Remove package with dependencies arch pacman yaourt"
date: 2017-05-12 17:12:56 +0200
categories: linux arch pacman yaourt
redirect_from:
  - /post/remove-package-with-dependencies-arch-pacman-yaourt
---

I mostly use yaourt instead of pacman to make my life easier when installing AUR packages. But they have the same syntax for this specific task.

Remove a package along with its installed dependencies

    $ yaourt -Rcs <package>

- `R` remove
- `c` Remove packages that are no longer installed from the cache as well as currently unused sync databases
- `s` Remove each target specified including all of their dependencies, provided that (A) they are not required by other packages; and (B) they were not explicitly installed by the user.

## References

- <https://www.reddit.com/r/archlinux/comments/2dsvp5/how_can_i_remove_all_dependencies_of_an_aur/>
- `man yaourt`
- `man pacman`