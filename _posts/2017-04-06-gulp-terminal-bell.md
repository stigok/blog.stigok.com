---
layout: post
title: "Gulp terminal bell"
date: 2017-04-06 14:33:03 +0200
categories: gulp node ascii
redirect_from:
  - /post/gulp-terminal-bell
---

After being enlightened by the eventual solution, the title could also be

> How to trigger system bell with Node.js

But *Gulp terminal bell* is what I initially searched for, so here it goes.

It's very much like triggering a bell in any other programming language. When STDOUT or STDERR is directed to a terminal or terminal emulator, echoing a `BEL` character will trigger the system bell.

So with Gulp (Node.js) it will be as simple as writing the character to the console in one of the following notations:

    console.log('\7')
    console.log('\x07')
    console.log('\u0007')

You can also write to stdout with `process.stdout.write` (`console.log` does exactly that).

I have not gone deep and checked all the different notations on different platforms, but they all work in my Terminator terminal with node v6.10.0.

## Using it in Gulp

Right now I'm doing a lot of web development and using gulp to automatically build on changes with livereload, but when the build fails I need to be notified somehow. And when I'm working on Windows, I want the taskbar item to lighten up on build errors. Using the solution above does exactly what I want in an uncomplicated manner.

    // Print error and a bell character
    function printAndSwallowError (err) {
      console.error(err.toString())
      console.error('\x07')
      this.emit('end')
    }

    pipes.buildScripts = () => gulp.src(src.js)
      .pipe(plugins.babel())
      .on('error', printAndSwallowError)
      .pipe(plugins.concat('bundle.min.js'))
      .pipe(gulp.dest(dist.js))

In this example, I'm also *swallowing* the error by simply listening for it. If no event listeners are attached, the error will bubble up and stop execution of the gulp process. With this setup however, I get notified when I have build errors, fix the errors, the re-save and continue on without having to restart any processes.

Other people use [gulp-plumber](https://github.com/floatdrop/gulp-plumber) for the same thing, but let's keep it simple. 

## References
- <http://www.asciitable.com/>
- <https://gist.github.com/taterbase/3154646>