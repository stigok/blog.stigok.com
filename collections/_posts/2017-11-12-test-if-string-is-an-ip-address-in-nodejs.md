---
layout: post
title: "Test if string is an IP address in Node.js"
date: 2017-11-12 14:12:58 +0100
categories: node.js
redirect_from:
  - /post/test-if-string-is-an-ip-address-in-nodejs
---

I just found some new (since v0.3.0) functions in the Node.js `net` module. These functions enables you to check if you have a valid IP address regardless of family, or if it is specifically an IPv4 or IPv6 address.

    const net = require('net')
    const assert = require('assert')
    
    assert.strictEqual(net.isIP('127.0.0.1'), 4);
    assert.strictEqual(net.isIP('x127.0.0.1'), 0);
    assert.strictEqual(net.isIP('example.com'), 0);
    assert.strictEqual(net.isIP('0000:0000:0000:0000:0000:0000:0000:0000'), 6);

    assert.strictEqual(net.isIPv4('127.0.0.1'), true);
    assert.strictEqual(net.isIPv4('example.com'), false);

    assert.strictEqual(net.isIPv6('127.0.0.1'), false);
    assert.strictEqual(net.isIPv6('2001:252:0:1::2008:6'), true);

The above tests are taken from [the upstream repository](https://github.com/nodejs/node/blob/98e54b0/test/parallel/test-net-isip.js).

## References
- https://nodejs.org/dist/latest/docs/api/net.html#net_net_isip_input
- https://github.com/nodejs/node/blob/98e54b0/test/parallel/test-net-isip.js