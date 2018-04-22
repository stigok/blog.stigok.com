---
layout: post
title: "Notes on cryptsetup LUKS2 format"
date: 2017-12-13 10:50:37 +0100
categories: cryptsetup luks security encryption
redirect_from:
  - /post/notes-on-cryptsetup-luks2-format
---

> This is **not** an exhaustive list of all new features or changes in cryptsetup 2.0.0, just my own notes on the things that matters to me right now. A link to the full upstream changelog in the reference list at the bottom of this page.

- Backwards compatible with *luks1* forever
- Integrates with dm-integritycheck for integrity protection with `--integrity`
  - Will write a new, random IV, on every write to the disk per-sector metadata as a nonce to the authenticated encryption algorithm (AEAD)
  - Journaling will decrease available disk space
  - For e.g. SSD's with integrated integrity checking this feature can be disabled with `--integrity-no-journal`

- For AEAD:
  - Supports *aes-xts-plain64* with *hmac-sha256* or *hmac-sha512* for authentication tag
  - *aes-gcm-random* (native AEAD mode) (**Not ready for production due to obvious security flaws**)
  - ChaCha20 with Poly1305 authenticator

- New *memory-hard* PBKDF to increase memory cost for known attacks now that it is feasible to run dictionary and brute-force attacks in parallel on GPUs
  - Support for Argon2i and Argon2id with the latter preffered for side channel and GPU cracking resistance
  - Argon2 is dual licenced under Creative Commons CC0 1.0 and Apache Public License 2.0

- Can use kernel keyring
  - to store volume key for dm-crypt where it avoids sending volume key in every device-mapper ioctl structure
  - to automatically unlock LUKS device if a passphrase is put into kernel keyring and proper keyring token is configured
  - can be used together with TPM

## References
- [cryptsetup 2.0.0 Release Notes](https://gitlab.com/cryptsetup/cryptsetup/blob/master/docs/v2.0.0-ReleaseNotes)
- [Argon2](https://en.wikipedia.org/wiki/Argon2)