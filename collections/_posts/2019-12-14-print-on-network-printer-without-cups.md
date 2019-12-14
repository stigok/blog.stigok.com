---
layout: post
title:  "Print a PDF on a printer without CUPS or other driver"
date:   2019-12-14 19:08:10 +0100
categories: linux
---

I think CUPS is too much. I just want to print simple files from the command
line. This is how it works for me with the printers I have been connecting too.

Depending on how your printer is connected, the target to which you can pipe
to will be either

- if connected through the **local network**, you can use `netcat`,
  `nc`, `lp`, or even just pipe to `/dev/tcp/<ip>/<port>` on some systems.
- if connected through **USB**, the printer might have been made available
  as a USB device in `/dev/usb/<something>`.

In the examples below I am using a network printer and connecting to it
using the following information:

```
$ export PRINTER_IP=192.168.0.5 PRINTER_PORT=9100
```

### Print plain text or raw PostScript

If you simply want to print plain text, you can pipe it straight through

```
$ echo "hello, world!" | nc -w 1 $PRINTER_IP $PRINTER_PORT
```

### Print a PDF

If you want to print a PDF, you most probably have to convert it to PostScript
first. I use `pdf2ps` for this. Using `-` as a target filename will output to
*stdout*.

```
$ pdf2ps ~/my-file.pdf - | nc -w 1 $PRINTER_IP $PRINTER_PORT
```

### Using `lp`

If you want to use `lp` for printing, use the `-h` argument. Specifying `-` as
input filename will make it read from *stdin*.

```
$ echo foobar | lp -h ${PRINTER_IP}:${PRINTER_PORT} -
```

## References
- <https://stackoverflow.com/questions/30130190/linux-print-directly-to-network-printer-that-is-not-installed/>
