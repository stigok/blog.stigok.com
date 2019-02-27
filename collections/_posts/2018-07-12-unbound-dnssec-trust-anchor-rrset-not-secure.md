---
layout: post
title:  "DNS lookups times out and unbound trust anchors DNSKEY rrset is not secure"
date:   2018-07-12 12:20:54 +0200
categories: dns dnssec unbound
---

DNS hostname lookups timed out, and by looking at the logs, unbound was giving
errors when restarting

    juli 12 11:53:12 unbound[1361]: [1361:0] info: generate keytag query _ta-3d98-4a5c-4f66. NULL IN
    juli 12 11:53:12 unbound[1361]: [1361:0] info: failed to prime trust anchor -- DNSKEY rrset is not secure . DNSKEY IN
    juli 12 11:53:12 unbound[1361]: [1361:0] info: failed to prime trust anchor -- DNSKEY rrset is not secure . DNSKEY IN

It was configured to use the following *trusted-key.key* trust anchor file

    . IN DNSKEY 257 3 8 AwEAAaz/tAm8yTn4Mfeh5eyI96WSVexTBAvkMgJzkKTOiW1vkIbzxeF3+/4RgWOq7HrxRixHlFlExOLAJr5emLvN7SWXgnLh4+B5xQlNVz8Og8kvArMtNROxVQuCaSnIDdD5LKyWbRd2n9WGe2R8PzgCmr3EgVLrjyBxWezF0jLHwVN8efS3rCj/EWgvIWgb9tarpVUDK/b58Da+sqqls3eNbuv7pr+eoZG+SrDK6nWeL3c6H5Apxz7LjVc1uTIdsIXxuOLYA4/ilBmSVIzuDWfdRUfhHdY6+cn8HFRm+2hM8AnXGXws9555KrUB5qihylGa8subX2Nn6UwNR1AkUTV74bU=
    . IN DNSKEY 256 3 8 AwEAAYvxrQOOujKdZz+37P+oL4l7e35/0diH/mZITGjlp4f81ZGQK42HNxSfkiSahinPR3t0YQhjC393NX4TorSiTJy76TBWddNOkC/IaGqcb4erU+nQ75k2Lf0oIpA7qTCk3UkzYBqhKDHHAr2UditE7uFLDcoX4nBLCoaH5FtfxhUqyTlRu0RBXAEuKO+rORTFP0XgA5vlzVmXtwCkb9G8GknHuO1jVAwu3syPRVHErIbaXs1+jahvWWL+Do4wd+lA+TL3+pUk+zKTD2ncq7ZbJBZddo9T7PZjvntWJUzIHIMWZRFAjpi+V7pgh0o1KYXZgDUbiA1s9oLAL1KLSdmoIYM=
    . IN DNSKEY 257 3 8 AwEAAagAIKlVZrpC6Ia7gEzahOR+9W29euxhJhVVLOyQbSEW0O8gcCjFFVQUTf6v58fLjwBd0YI0EzrAcQqBGCzh/RStIoO8g0NfnfL2MTJRkxoXbfDaUeVPQuYEhg37NZWAJQ9VnMVDxP/VHL496M/QZxkjf5/Efucp2gaDX6RS6CXpoY68LsvPVjR0ZSwzz1apAzvN9dlzEheX7ICJBBtuA6G3LQpzW5hOA2hzCTMjJPJ8LbqF6dsV6DoBQzgul0sGIcGOYl7OyQdXfZ57relSQageu+ipAdTTJ25AsRTAoub8ONGcLmqrAmRLKBP1dfwhYB4N7knNnulqQxA+Uk1ihz0=

The file was created exactly one year ago, may be very relevant

    $ stat /etc/unbound/trusted-key.key
      File: /etc/unbound/trusted-key.key
      Size: 1107      	Blocks: 8          IO Block: 4096   regular file
    Device: 19h/25d	Inode: 1713967     Links: 1
    Access: (0644/-rw-r--r--)  Uid: (    0/    root)   Gid: (    0/    root)
    Access: 2018-07-12 11:51:55.711871713 +0200
    Modify: 2017-07-12 04:06:51.000000000 +0200
    Change: 2017-12-22 20:06:08.430190698 +0100
     Birth: -


Looking around, I found a tool called [get-trust-anchor][] which is a *tool
for fetching/refreshing DNS Root Zone trust anchors* that also verifies
the S/MIME signatures of the files.
It worked great for downloading a new DNSKEY record and provided me with
a file called *ksk-as-dnskey.txt* which contained a single record that was
identical to the last record of my existing trust anchor file:

    . IN DNSKEY 257 3 8 AwEAAagAIKlVZrpC6Ia7gEzahOR+9W29euxhJhVVLOyQbSEW0O8gcCjFFVQUTf6v58fLjwBd0YI0EzrAcQqBGCzh/RStIoO8g0NfnfL2MTJRkxoXbfDaUeVPQuYEhg37NZWAJQ9VnMVDxP/VHL496M/QZxkjf5/Efucp2gaDX6RS6CXpoY68LsvPVjR0ZSwzz1apAzvN9dlzEheX7ICJBBtuA6G3LQpzW5hOA2hzCTMjJPJ8LbqF6dsV6DoBQzgul0sGIcGOYl7OyQdXfZ57relSQageu+ipAdTTJ25AsRTAoub8ONGcLmqrAmRLKBP1dfwhYB4N7knNnulqQxA+Uk1ihz0=

However, copying this to the unbound configuration directory
(*/etc/unbound/ksk-as-dnskey.txt*) and using it as the new anchor file worked
great. Maybe unbound was trying to use the first record of my original file?

My working unbound configuration now looks like this:

    server:
      use-syslog: yes
      username: "unbound"
      directory: "/etc/unbound"
      trust-anchor-file: /etc/unbound/ksk-as-dnskey.txt
      tls-cert-bundle: /etc/ssl/certs/ca-certificates.crt

    forward-zone:
      name: "."
      forward-tls-upstream: yes
      # End-of-line comments (no spaces!) are used for TLS hostname verification
      forward-addr: 2a01:3a0:53:53::@853#unicast.censurfridns.dk
      forward-addr: 89.223.43.71@853#unicast.censurfridns.dk

Remember to restart the service after configuration changes

    # systemctl restart unbound

## References
- [get-trust-anchor][]
- [https://blog.uncensoreddns.org/]()

[get-trust-anchor]: https://github.com/iana-org/get-trust-anchor
