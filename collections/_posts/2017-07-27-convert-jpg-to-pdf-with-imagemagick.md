---
layout: post
title: "Convert JPG to PDF with ImageMagick"
date: 2017-07-27 20:10:27 +0200
categories: imagemagick convert pdf arch linux
redirect_from:
  - /post/convert-jpg-to-pdf-with-imagemagick
---

The `convert` binary is supplied by ImageMagick, so if you haven't already installed it, do so first. `ghostscript` is also necessary to create PDF's.

    # pacman -S imagemagick ghostscript

Convert all JPG files in current directory to PDF with `.pdf` suffix

    $ for f in ./*; do convert $f $f.pdf; done

Combine all pictures in a single PDF file

    $ convert "*.jpg" combined.pdf

## Problems
### ImageMagick produces invalid PDF's

If the resulting PDF's are corrupted or are otherwise not created as expected, make sure you actually installed `ghostscript`.