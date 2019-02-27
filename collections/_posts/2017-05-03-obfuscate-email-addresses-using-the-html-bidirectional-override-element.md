---
layout: post
title: "Obfuscate email addresses using the HTML bidirectional override element"
date: 2017-05-03 14:49:07 +0200
categories: html javascript obfuscation email
redirect_from:
  - /post/obfuscate-email-addresses-using-the-html-bidirectional-override-element
---

Starting with a reversed e-mail address

    'sshow@example.com'.split('').reverse().join('')
    // "moc.elpmaxe@wohss"

Putting it in a `<bdo>` element to make it readable again.

    <bdo dir="rtl">moc.elpmaxe@wohss</bdo>
    // sshow@example.com

The `dir` attribute specifies the text direction of the containing text. Either `rtl` or `ltr`.

I remember reading an article about different ways to obfuscate email addresses in websites to avoid feeding spambots and crawlers. Well, here's another one.

... whether this is a good idea or not is a discussion of its own.

Full example

    <!DOCTYPE html>
    <html>
      <head>
        <meta charset="utf-8">
        <script>
          function reverse(str) {
            return str.split('').reverse().join('')
          }
          function load () {
            var el = document.getElementById('email')
            el.innerHTML = reverse('sshow@example.com')
          }
        </script>
      </head>
      <body onload="load()">
        <bdo id="email" dir="rtl"></bdo>
      </body>
    </html>