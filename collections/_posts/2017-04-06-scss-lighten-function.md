---
layout: post
title: "SCSS lighten function"
date: 2017-04-06 18:55:14 +0200
categories: sass css
redirect_from:
  - /post/scss-lighten-function
---

SASS has a lighten function along with other color manipulation utilities (e.g. darken) which are really easy to use.

**pink-website.scss**

    $primary: #ff0000;

    body {
      background-color: lighten($primary, 30%);
      color: darken($primary, 20%);
    }

Which compiles into the following

    body {
      background-color: #ff9999;
      color: #990000;
    }


It takes two arguments; `color` and `amount`.

## Reference
- <http://sass-lang.com/documentation/Sass/Script/Functions.html#lighten-instance_method>
- my brainz