---
layout: post
title:  "Verify TLS certificates for DNS over TLS connections in unbound"
date:   2018-06-19 16:31:32 +0200
categories: unbound dns tls security
---

## Introduction

I was just notified by the operator of [uncensoreddns.org][] that I should be
pinning the public key of the DNS server's TLS certificate.

I am using unbound as my local DNS resolver and [uncensoreddns.org][] as my
"SSL upstream" forwarding zone.

I was surprised to learn that *unbound* did not perform any verification
by itself, and that I have been open to MITM-attacks just as easily as with
plaintext DNS over port 53 for a long ass time.

## Enable TLS certificate verification

Use `forward-tls-upsteam` option to use DNS over TLS. However, without
combining with `tls-cert-bundle`, no TLS certificate authentication will be
performed.

Here is a working example *unbound.conf* that performs validation that the
hostname matches the DNS hostname of the certificate:

```
server:
  use-syslog: yes
  username: "unbound"
  directory: "/etc/unbound"
  trust-anchor-file: trusted-key.key
  tls-cert-bundle: /etc/ssl/certs/ca-certificates.crt

forward-zone:
  name: "."
  forward-tls-upstream: yes
  # The below end-of-line comments (without spaces) are used for hostname
  # verification of the TLS certificate
  forward-addr: 2a01:3a0:53:53::@853#anycast.censurfridns.dk
  forward-addr: 89.223.43.71@853#anycast.censurfridns.dk
```

To verify that the hostname check is actually performed, try changing the
hostname that has been suffixed to the `forward-addr` lines into `#example.com`
and see what happens when requesting a lookup with
`dig @127.0.0.1 blog.stigok.com`:

```
systemd[1]: Started Unbound DNS Resolver.
unbound[1221]: [1221:0] notice: init module 0: validator
unbound[1221]: [1221:0] notice: init module 1: iterator
unbound[1221]: [1221:0] info: start of service (unbound 1.7.2).
unbound[1221]: [1221:0] error: ssl handshake failed crypto error:1416F086:SSL routines:tls_process_server_certificate:certificate verify failed
unbound[1221]: [1221:0] notice: ssl handshake failed 2a01:3a0:53:53:: port 853
unbound[1221]: [1221:0] error: ssl handshake failed crypto error:1416F086:SSL routines:tls_process_server_certificate:certificate verify failed
unbound[1221]: [1221:0] notice: ssl handshake failed 2a01:3a0:53:53:: port 853
```

## References
- https://blog.uncensoreddns.org/
- https://www.nlnetlabs.nl/bugs-script/show_bug.cgi?id=658#c9

[uncensoreddns.org]: https://uncensoreddns.org
