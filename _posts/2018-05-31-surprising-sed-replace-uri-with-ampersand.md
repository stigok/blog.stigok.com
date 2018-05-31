---
layout: post
title: "Surprised when sed replacing a token with a URI"
date: 2018-05-31 15:07:42 +0200
categories: linux sed perl
---

As I was replacing tokens with secrets in an init container for a Kubernetes
Pod, I was not getting it to work correctly. Not until I actually took a look
at the resulting configuration file, I was getting some hints as to what was
wrong. I had used sed before when replacing tokens with URI's, but apparently
not URI's with multiple query-string parameters which are joined by
ampersands (`&`).

{{ raw }}
Let's replace the token `%link%` with a URI
{{ endraw }}

{{ raw }}
    $ URI='https://example.com/?num=7'
    $ echo 'my_link = %link%' | sed "s|%link%|${URI}|"
    my_link = https://example.com/?num=7
{{ endraw }}

Seems good. Let's add another parameter to the URI

{{ raw }}
    $ URI='https://example.com/?num=42&type=magic'
    $ echo 'my_link = %link%' | sed "s|%link%|${URI}|"
    my_link = https://example.com/?num=42%link%type=magic
{{ endraw }}

Surprise. What I expected the link to look like was
`https://example.com/?num=42&type=magic`

A quick [search][search] revealed [the culprit][answer];

> The REPLACEMENT can contain `\N` (N being a number from 1 to 9,
>  inclusive) references, which refer to the portion of the match which
>  is contained between the Nth `\(` and its matching `\)`.  Also, the
>  REPLACEMENT can contain unescaped `&` characters which reference the
>  whole matched portion of the pattern space

What I could do here is to first sanitise the URL by backslashing all
occurences of `&` before using it as the replacement

{{ raw }}
    URI='https://example.com/?num=42&type=magic'
    $ echo 'my_link = %link%' | sed "s|%link%|$(echo "$URI" | sed 's|&|\\&|g')|"
    my_link = https://example.com/?num=42&type=magic
{{ endraw }}

But to me, that seems clunky. I would rather switch to a different replacer
to avoid corner cases like this and excessive code in the init containers.
Trying with a `perl` script instead

{{ raw }}
    URI='https://example.com/?num=42&type=magic'
    $ echo 'my_link = %link%' | perl -pi -e "s|%link%|${URI}|"
    my_link = https://example.com/?num=42&type=magic
{{ endraw }}

This works like I want it to. It brings a dependency on perl, but I am
handling a stock debian image right now, so I don't really care.

## References
- https://unix.stackexchange.com/questions/341644/using-perl-to-replace-a-string-with-contents-from-file-found-in-an-array
- https://unix.stackexchange.com/questions/204156/perl-sed-replacement
- `perl -h`
- https://stackoverflow.com/questions/32750591/sed-behavior-with-ampersand#32750675

[answer]: https://stackoverflow.com/questions/32750591/sed-behavior-with-ampersand#32750675
[search]: https://duckduckgo.com/?q=sed+ampersand

