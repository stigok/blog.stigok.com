---
layout: post
title: "This is now a Jekyll blog"
date: 2018-04-27 11:16:42 +0200
categories: jekyll
---

I moved my blog to Jekyll as I was getting tired of maintaining my homemade, over-engineered
blog program with a Node backend and mongodb database. It simply wasn't easy enough to create
new posts, and creating new features was getting increasingly dulling.

It has become apparent how important it is to have simple workflows for small tasks. Creating
a blog post should not have to be anything more than writing the content and pressing *Post*.

I wanted was to open vim, write some markdown and publish that post with a single command.

Jekyll was perfect match for this

  1. Write markdown
  2. Preview `bundle exec jekyll serve --limit_posts=5`
  2. Commit `git commit -am 'New post'`
  3. Push `git push`

And that's it (assuming my server is auto-pulling, which it is)!
