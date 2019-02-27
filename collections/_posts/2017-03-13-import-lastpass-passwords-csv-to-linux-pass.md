---
layout: post
title: "Import LastPass passwords csv to linux pass"
date: 2017-03-13 03:13:35 +0100
categories: pass lastpass linux csv
redirect_from:
  - /post/import-lastpass-passwords-csv-to-linux-pass
---

I finally made the move. Exported all my LastPass passwords to CSV, then made a tool to import it to `pass`.

    $ git clone https://github.com/stigok/lastpass-to-pass
    $ cd lastpass-to-pass
    $ npm install
    $ cat ~/lastpass.csv | node index.js

Check repo for more info <https://github.com/stigok/lastpass-to-pass>