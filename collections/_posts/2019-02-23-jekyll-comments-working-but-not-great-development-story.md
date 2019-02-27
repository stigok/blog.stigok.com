---
layout: post
title:  "The working, but not great, Jekyll comment plugin development story"
date:   2019-02-23 17:03:49 +0100
categories: jekyll ruby deployment
---

I was building a comment API to accept comments to my Jekyll blog (this blog).
I wanted to be able to accept comments and get them back out again.
This had me planning for two API endpoints:

- `PUT` /comments
- `GET` /comments/:post_id

This seemed simple enough, so I started out making the internal API
to accept comments, set up persistent storage with SQLite then went
on and wrote out the HTTP bindings. This project would be written in
Node.js, because I wanted to write some pretty JavaScript.
I was very happy with my efforts so far, and I even had
tests for the pretty little thing.

When I was done and things worked out alright, it was time to put into production
on my web server. First issue I stumbled into, (which is not the first time, at all)
was to have the Node.js version available on my server. My server, however did not
have any modern versions of Node.js available in the standard repositories.
But I didn't want to pollute the server with new repos as I did not want
to burden my application with too much deployment steps.

I ended up packing it into a Dockerfile, as I already had Docker running
on the server. Then I had to expose a local directory to have a place
to store the sqlite database file. I didn't really like this too much either.
Something felt wrong.

When everything was up and running, and the Jekyll generator plugin worked as it should
too, I realized this wasn't what I wanted at all. I was indeed very happy with finishing
this small project to the end, but I was somewhat surprised at how wrong everything
about this solution felt.

I had a static site whos strength was that it was generated and thereby only consisting of
plain files stuffed into a git repo. This was beautiful. But now, I had introduced
an external dependency, of which could not live in the same repo.
This was indeed not what I wanted.

Now I'm starting over again with a Python 3.5 HTTP back-end which spits out
comment files as they come in, then adds them to the repository.
3.5 because it seems to be included in most of the mainstream distros these days.
I imagine things will turn out better this time around. Let's see...

## Links

- [The code to the Node.js SQLite solution](https://github.com/stigok/express-jekyll-comments/)
- [Jekyll comments generator plugin](https://github.com/stigok/blog.stigok.com/blob/6e3e48916fb9d9c8010153433f62714bdcf19f57/_plugins/comments.rb) that uses `curl` to download comments

The code for the Python 3.5 plain file solution may just pop up some day

