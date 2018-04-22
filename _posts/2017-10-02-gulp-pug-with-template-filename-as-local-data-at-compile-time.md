---
layout: post
title: "gulp pug with template filename as local data at compile time"
date: 2017-10-02 14:02:50 +0200
categories: gulp pug html
redirect_from:
  - /post/gulp-pug-with-template-filename-as-local-data-at-compile-time
---

I wanted to get the filename of the compiled template at compile-time to handle `active` link styles in my templates.

    const pug = require('pug')
    const path = require('path')

    const pugWithFilenameAsLocal = {
      compile: function (str, options) {
        options.data = Object.assign(options.data || {}, {
          filename: path.basename(options.filename).replace('.pug', '.html')
        })
        return pug.compile(str, options)
      }
    }

    gulp.task('html', () => gulp.src([paths.pug, '!**/_*.pug'])
      // Use custom gulp to have the filename available as a local in all compiled templates
      .pipe(plugins.pug({pug: pugWithFilenameAsLocal}))
      .pipe(gulp.dest(paths.dist))
    )

So now I can make conditional `active` link items like this in my layout template

    //- _navigation.pug
    -
      var links = [
        {href: 'projects.html', text: 'Prosjekter'},
        {href: 'things.html', text: 'Things and stuff'},
        {href: 'foobar.html', text: 'Foobar'}
      ]
    nav
      ul
        each link in links
          li #[a(href=link.href class=(filename == link.href ? 'active' : ''))= link.text]

And then include that in the individual sub templates

    //- project.pug
    doctype html
    html
      head
        title Website name
        link(rel='stylesheet' href='css/bundle.css')
      body
        include _nav
        h1 My projects
        p Content goes here, you know...

## References
- https://stackoverflow.com/q/12262844/90674