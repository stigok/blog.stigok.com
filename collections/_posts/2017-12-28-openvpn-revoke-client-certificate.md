---
layout: post
title: "OpenVPN revoke client certificate"
date: 2017-12-28 01:12:21 +0100
categories: openvpn linux security
redirect_from:
  - /post/openvpn-revoke-client-certificate
---

I transfered a VPN configuration file over an insecure connection by a mistake and I had to revoke the client certificate. Here are the steps to revoke a client certificate and deny access for the client to the OpenVPN server.

For revocations to have any effects, the OpenVPN server instance should be configured with `crl-verify`. This option takes a path to a file which has to be world readable (0644) since the server reads this file *after* root privileges has been dropped.

Put the following line into the OpenVPN server configuration file

    crl-verify /usr/share/openvpn/crl.pem

Make sure the file exists and is readable

    # touch /usr/share/openvpn/crl.pem
    # chmod 0644 /usr/share/openvpn/crl.pem

I've set up my PKI with `easy-rsa` into `/ca/openvpn/`. Source the vars and revoke a client by its common name (CN). In this example, the client is *stigok*.

    # cd /ca/openvpn
    # source ./vars
    # revoke-full stigok 

You will be asked for the passphrase for the CA private key (if you've set one) twice. First time to revoke it, second time to attempt to verify the client certificate. This command is expected to return with an error:

    error 23 at 0 depth lookup:certificate revoked

If this is the first revocation for the CA, a new file `crl.pem` should have been created in the current working directory. Copy it to the world readable location so OpenVPN knows about it.

    # cp crl.pem /usr/share/openvpn/

The OpenVPN server will re-read this file upon new client connections and also on every TLS renegotiation. Alternatively, updates can be propagated immediately by signaling the server process. This will, however, reset all active connections.

    # pkill -HUP openvpn

Attempts to connect with the revoked certificate should fail, however silently on the client side, as connections will appear to simply just not complete. Server side it will be declined and logged with a message describing which certificate was attempted, originating IP-address, and note that it has been revoked:

    11.22.33.44:58635 TLS: Initial packet from [AF_INET]11.22.33.44:58635, sid=fc7be805 1e8467be
    11.22.33.44:58635 CRL CHECK OK: C=NO, ST=Norway, L=Oslo, O=example.com, OU=vpn.example.com, CN=example.com CA, name=OpenVPN, emailAddress=vpn@example.com
    11.22.33.44:58635 VERIFY OK: depth=1, C=NO, ST=Norway, L=Oslo, O=example.com, OU=vpn.example.com, CN=example.com CA, name=OpenVPN, emailAddress=vpn@example.com
    11.22.33.44:58635 CRL CHECK FAILED: C=NO, ST=Norway, L=Oslo, O=example.com, OU=vpn.example.com, CN=stigok, name=OpenVPN, emailAddress=vpn@example.com is REVOKED
    11.22.33.44:58635 TLS_ERROR: BIO read tls_read_plaintext error: error:140890B2:SSL routines:SSL3_GET_CLIENT_CERTIFICATE:no certificate returned
    11.22.33.44:58635 TLS Error: TLS object -> incoming plaintext read error
    11.22.33.44:58635 TLS Error: TLS handshake failed

## References
- https://openvpn.net/index.php/open-source/documentation/howto.html#revoke