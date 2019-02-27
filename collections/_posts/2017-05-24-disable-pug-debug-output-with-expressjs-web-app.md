---
layout: post
title: "Disable pug debug output with express.js web app"
date: 2017-05-24 15:24:09 +0200
categories: pug express.js javascript
redirect_from:
  - /post/disable-pug-debug-output-with-expressjs-web-app
---

While developing an express app and using pug, as long as `NODE_ENV !== 'production'` pug outputs function bodies to stdout. I find that pretty annoying, and I only want it to output stuff to console when theres actually an error.

Can be done by setting a custom engine function for pug

    app.engine('pug', (path, options, fn) => {
      options.debug = false
      return pug.__express.call(null, path, options, fn)
    })
    app.set('view engine', 'pug')

## References
- [source for `pug.__express`](https://github.com/pugjs/pug/blob/897c7779fb53b7281be6cc9e61991281ee4be443/packages/pug/lib/index.js#L461-L466)
- http://expressjs.com/en/4x/api.html#app.engine
- https://pugjs.org/api/reference.html