---
layout: post
title: "Generate password hash for ASP.NET Identity authentication"
date: 2018-01-03 01:06:40 +0100
categories: C# asp.net security
redirect_from:
  - /post/generate-password-hash-for-aspnet-identity-authentication
---

I had to create a password hash to update a password directly in a MSSQL server instance, but did not have Visual Studio to launch the Auth Web UI to administer the user accounts, nor did I have a C# IDE or compiler installed. I wanted to solve it with JavaScript anyway, so there goes.

This code, written for Node.js 6.x, generates a password hash to use with ASP.NET Identity version 2.

The resulting base64-encoded string consist of a single byte `0` signifying this is a version 2 hash. Then comes a 16 bytes long salt, followed by 32 bytes long SHA-1 hashed (1000 rounds) password.

    const crypto = require('crypto')

    const pass = 'foobar'
    const version = Buffer.alloc(1) // 1 zero-filled byte
    const salt = crypto.randomBytes(16)
    const hash = crypto.pbkdf2Sync(pass, salt, 1000, 32, 'sha1')
    const b64 = Buffer.concat([
      version,
      salt,
      hash
    ]).toString('base64')
    
    console.log(b64)

## References
- http://www.blinkingcaret.com/2017/11/29/asp-net-identity-passwordhash/