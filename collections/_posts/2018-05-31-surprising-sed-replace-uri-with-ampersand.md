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

Let's replace the token `%link%` with a URI

    $ URI='https://example.com/?num=7'
    $ echo 'my_link = %link%' | sed "s|%link%|${URI}|"
    my_link = https://example.com/?num=7

Seems good. Let's add another parameter to the URI

    $ URI='https://example.com/?num=42&type=magic'
    $ echo 'my_link = %link%' | sed "s|%link%|${URI}|"
    my_link = https://example.com/?num=42%link%type=magic

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

    URI='https://example.com/?num=42&type=magic'
    $ echo 'my_link = %link%' | sed "s|%link%|$(echo "$URI" | sed 's|&|\\&|g')|"
    my_link = https://example.com/?num=42&type=magic

But to me, that seems clunky. I would rather switch to a different replacer
to avoid corner cases like this and excessive code in the init containers.
Trying with a `perl` script instead

    URI='https://example.com/?num=42&type=magic'
    $ echo 'my_link = %link%' | perl -p -e "s|%link%|${URI}|"
    my_link = https://example.com/?num=42&type=magic

This seems to be okay, but let's add login credentials and see what happens

    URI='https://user:pass@example.com/?num=42&type=magic'
    $ echo 'my_link = %link%' | perl -p -e "s|%link%|${URI}|"
    my_link = https://user:pass.com/?num=42&type=magic

So this time around perl uses `@` for named capture groups. In which case
I first have to escape them like `\@`. So I'm where I started out again.

Let's try with `awk`

    URI='https://user:pass@example.com/?num=42&type=magic'
    $ echo 'my_link = %link%' | awk "{ gsub(/%link%/, \"$URI\"); print }"
    my_link = https://user:pass@example.com/?num=42%link%type=magic

Same problem as with `sed` in that it replaces `&` with the matching string.

I yield. I found a [nice script][script] written by Ed Morton which sanitises
both the search and the replacement string before running `sed` one last time.
I added an option to do in-place file replacements.

    #!/bin/sh
    # Modified version of Ed Morton's sedstr:
    # https://stackoverflow.com/a/29626460/90674
    # stigok, 2018

    # In-place option
    opts=''
    if [ "$1" = '-i' ]; then
      opts='-i ' # note trailing space
      shift
    fi

    old="$1"
    new="$2"
    file="${3:--}"
    escOld=$(sed 's/[^^]/[&]/g; s/\^/\\^/g' <<< "$old")
    escNew=$(sed 's/[&/\]/\\&/g' <<< "$new")
    sed ${opts}"s/$escOld/$escNew/g" "$file"

This is finally a working solution for me with all the possible special
characters of the URI in question, without breaking or triggering unwanted
features and variable expansions in sed and bash.

Cheers!

## References
- https://unix.stackexchange.com/questions/341644/using-perl-to-replace-a-string-with-contents-from-file-found-in-an-array
- https://unix.stackexchange.com/questions/204156/perl-sed-replacement
- `perl -h`
- https://stackoverflow.com/questions/32750591/sed-behavior-with-ampersand#32750675

[answer]: https://stackoverflow.com/questions/32750591/sed-behavior-with-ampersand#32750675
[search]: https://duckduckgo.com/?q=sed+ampersand
[script]: https://stackoverflow.com/questions/29613304/is-it-possible-to-escape-regex-metacharacters-reliably-with-sed/29626460#29626460
