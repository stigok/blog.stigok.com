---
layout: post
title: "Express.js can't set headers after they are sent"
date: 2017-06-26 00:17:39 +0200
categories: expressjs nodejs
redirect_from:
  - /post/expressjs-cant-set-headers-after-they-are-sent
---

> `Error: Can't set headers after they are sent.`

> **NOTE** If you are not using `res.format()` within the route function, this
> solution may not match your current issues

The error message is pretty clear, and the stack trace thrown with it gives you
the exact line number of the error. But the code inducing it might not be so
easy to make sense of. What's wrong with the below code?

    app.get('/foo', (req, res, next) => {
      return res.format({
        text: res.send('bar')
        html: res.send('<strong>bar</strong>')
      })
    })

When the route function is executed (`app.get()`), the object literal within the
`res.format` parenthesis are invoked before `res.format` itself. This means the
response is in fact sent before `res.format` has even gotten started with
setting the response headers. The correct way to go about this is to use
functions as property values for the supplied object to `res.format`:

    app.get('/foo', (req, res, next) => {
      return res.format({
        text: () => res.send('bar')
        html: () => res.send('<strong>bar</strong>')
      })
    })