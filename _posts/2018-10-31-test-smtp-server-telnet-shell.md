---
layout: post
title:  "Test an SMTP server from commandline using telnet"
date:   2018-10-31 12:09:35 +0100
categories: email
---

## Introduction

Whenever I look at protocols I always get positively surprised at how easy
they are to understand once you take an actual look at it.

This time around I wanted to test an SMTP server without the need to set up
some random test client software or putting something together in Python or
Node.js. This time, I wanted something plain simple.

## telnet

telnet is a tool that is very handy in talking with servers. No magic goes
on under the hood. Perform simple conversations from a client to a server.

To connect to a server on the default plain-text SMTP port (25) you can
use `telnet target.host.example.com smtp` and start the conversation.

Below is an example transmission of sending an email over SMTP. All lines
starting with three numbers are responses from the server. `HELO`, `MAIL`,
`RCPT`, `DATA`, and `QUIT` are commands written and sent from the client
to the server.

```bash
$ telnet mail.example.com smtp
Trying 127.0.0.1...
Connected to localhost.localdomain (127.0.0.1).
Escape character is '^]'.
220 mail.example.com ESMTP Sendmail 8.13.8/8.13.8; Tue, 22 Oct 2013 05:05:59 -0400
HELO myserver.no
250 myclient.hostname.example.no Hello myclient.hostname.example.no [127.0.0.1], pleased to meet you
MAIL from: sender@myclient.hostname.example.no
250 2.1.0 sender@myclient.hostname.example.no... Sender ok
RCPT to: mom@example.com
250 2.1.5 mom@example.com... Recipient ok
DATA
354 Enter mail, end with "." on a line by itself
Hey
This is test email only

Thanks
.
250 2.0.0 r9M95xgc014513 Message accepted for delivery
QUIT
221 2.0.0 mail.example.com closing connection
Connection closed by foreign host.
```

The same as above with server replies and telnet output omitted:

```bash
HELO myserver.no
MAIL from: sender@myclient.hostname.example.no
RCPT to: mom@example.com
DATA
Hey
This is test email only

Thanks
.
QUIT
```

Now it's easy to see how simple the SMTP protocol is.

## References
- https://tecadmin.net/ways-to-send-email-from-linux-command-line/
