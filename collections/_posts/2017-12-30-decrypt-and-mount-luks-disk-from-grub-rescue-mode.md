---
layout: post
title: "Decrypt and mount LUKS disk from GRUB rescue mode"
date: 2017-12-30 00:42:00 +0100
categories: grub luks security recovery
redirect_from:
  - /post/decrypt-and-mount-luks-disk-from-grub-rescue-mode
---

I am running a Linux installation with an encrypted boot partition using LUKS and GRUB. From time to time I enter the wrong password to GRUB which dumps me into grub rescue mode. There is no help, nor any tab completion to get when in that prompt, so you better know your way around. The only other option, which may actually be faster, is to reboot to get another shot. Anyway, here's how to open a LUKS volume from grub rescue mode and continue booting without rebooting.

Start out by entering an invalid password to GRUB bootloader

    Welcome to GRUB!

    Attempting to decrypt master key...
    Enter passphrase for hd0,gpt2 (<disk uuid>):
    error: access denied
    error: no such cryptodisk found.
    error: disk `cryptouuid/<disk uuid>` not found.
    Entering rescue mode...

List all devices found (out of curiosity)

    grub rescue> ls

Mount the encrypted */boot* partition (as attempted from the start). I know that partition number 2 on the first (and only) disk is mine, hence `(hd0,gpt2)`.

    grub rescue> cryptomount (hd0,gpt2)

If you want to specify a UUID instead, use the `-u` option to [cryptomount][].

Try to enter the correct passphrase

    Attempting to decrypt master key...
    Enter passphrase for hd0,gpt2 (<disk uuid>):

Output on success:

    Slot 3 opened

Load the module for normal boot

    grub rescue> insmod normal

Boot normally as GRUB tried to do at first

    grub rescue> normal

Now you should be taken to the next step in the boot process, which in my case is the GRUB OS selection menu. Works for me!

## Other commands

The rescue shell does not have a whole lot of commands. It doesn't even have a
`help` command. For some more information about the available functions, consult
the GRUB manual section for "*[Commands][1]*" and
the subsection "*[GRUB only offers a rescue shell][2]*".

## References
- <https://www.coreboot.org/GRUB2>
- <https://www.gnu.org/software/grub/>

[cryptomount]: https://www.gnu.org/software/grub/manual/grub/html_node/cryptomount.html#cryptomount
[1]: https://www.gnu.org/software/grub/manual/grub/html_node/Commands.html#Commands
[2]: https://www.gnu.org/software/grub/manual/grub/html_node/GRUB-only-offers-a-rescue-shell.html#GRUB-only-offers-a-rescue-shell
