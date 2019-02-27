---
layout: post
title: "Install a package from local directory in Atom"
date: 2017-02-26 00:17:03 +0100
categories: atom
redirect_from:
  - /post/install-a-package-from-local-directory-in-atom
---

# Manually install a package in Atom

I had to do this while developing a custom package.

    $ apm link <directory>

And you're done!

## What is that command really doing?

It is simply linking the specified directory to the Atom packages folder,
typically residing in `~/.atom/packages`.

## References

- https://discuss.atom.io/t/manually-install-package/9251/2
