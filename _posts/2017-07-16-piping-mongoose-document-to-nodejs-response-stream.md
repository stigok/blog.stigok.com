---
layout: post
title: "Piping mongoose document to Node.js response stream"
date: 2017-07-16 05:11:04 +0200
categories: mongoose mongodb express javascript node.js
redirect_from:
  - /post/piping-mongoose-document-to-nodejs-response-stream
---

You can keep the mongoose documents as streams by calling `cursor()` and piping directly to the response stream of the HTTP server. In this example I'm using it in an Express.js router route definition. Since Node.js streams expect stream data to be of type `String` or `Buffer`, the JavaScript object needs to be stringified first using a `transform` function. `JSON.stringify` can be used directly which helps in making this a nice pipeline.

    router.get('/', (req, res, next) => {
      Model.find()
        .cursor({transform: JSON.stringify})
        .pipe(res.type('json'))
    })

## References
- https://stackoverflow.com/questions/20058614/stream-from-a-mongodb-cursor-to-express-response-in-node-js/45124486#45124486