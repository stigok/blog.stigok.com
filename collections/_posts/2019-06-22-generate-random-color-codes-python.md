---
layout: post
title:  "Generating random color codes in Python 3"
date:   2019-06-22 11:30:39 +0200
categories: python
---

I was looking at question from Stack Overflow asking
[How to use `random` to choose colors][so-question].
I wanted to take a shot at it before reading the answers again.

- I was thinking of a 24-bit RGB color in hex representation, e.g. `#ff9900`
  (I come from HTML and CSS)
- But hex is just a representation/format, like octal and decimal.
  In decimal, a 24-bit color is a number between 0 and 16777216. (`16**6`
  which means `0-F` 6 times)
- The web colors consist of red, green and blue (RGB) channels. Maybe I should
  return a byte for each channel instead of all at once. Maybe it matters in
  terms of distribution?
- Can I do this with `str.format` or `%` formatting?

At first I made a function to return a hexadecimal string

```python
import random

def rand_web_color_hex():
    rgb = ""
    for _ in "RGB":
        i = random.randrange(0, 2**8)
        rgb += i.to_bytes(1, "big").hex()
    return rgb

print(rand_web_color_hex()) # 2f04d8
print(rand_web_color_hex()) # 8dbc53
print(rand_web_color_hex()) # 022632
```

Then I thought I might make a function to return a decimal number, too

```python
import random

def rand_web_color_dec():
    return random.randrange(0, 16**6)

print(rand_web_color_dec()) # 7431420
print(rand_web_color_dec()) # 12678862
print(rand_web_color_dec()) # 582912
```

Naturally, a lot smaller. A simple call to `randrange()` to return a 24-bit
integer. Makes me want to revisit the hex function to simplify it. Hex is
simply a number formatted in base16, and for web colors, they are zero-padded
to maintain their length of 6 hexadecimal chars. I know there's a way
to return numbers in hex notation using `str.format` or `%` formatting.

```python
def rand_web_color_hex():
    return "%06x" % random.randrange(0, 16**6)

print(rand_web_color_hex()) # 06dffd
print(rand_web_color_hex()) # 6f0e2a
print(rand_web_color_hex()) # 1d3df2
```

In the above code, I'm using the `%` string formatting operator from [PEP 3101][],
where `%x` will format a number user lower-case hex notation, `6` pads the string
with spaces to make it consist of at least six characters, and lastly, the `0` will
use zeros for padding instead of spaces.

Now, if I want to return a color with separated channels, like `rgb(r, g, b)` format, I must add another function. I can use `bytearray.fromhex()` to convert from hex
to a byte array.

```python
def rand_web_color_hex():
    return "%06x" % random.randrange(0, 16**6)

def to_rgb(color_str):
    barr = bytearray.fromhex(color_str)
    return (barr[0], barr[1], barr[2])

hex_color = rand_web_color_hex()
print(hex_color, to_rgb(hex_color)) # 958c3b (149, 140, 59)
```

Cleaning up and refactoring a bit, yields the following code.

```python
import random

def rand_24_bit():
    """Returns a random 24-bit integer"""
    return random.randrange(0, 16**6)

def color_dec():
    """Alias of rand_24 bit()"""
    return rand_24_bit()

def color_hex(num=rand_24_bit()):
    """Returns a 24-bit int in hex"""
    return "%06x" % num

def color_rgb(num=rand_24_bit()):
    """Returns three 8-bit numbers, one for each channel in RGB"""
    hx = color_hex(num)
    barr = bytearray.fromhex(hx)
    return (barr[0], barr[1], barr[2])

c = color_dec()
print("%8d #%6s rgb%s" % ( c, color_hex(c), color_rgb(c) ))
```

This was a good exercise for maths and formatting in Python. Hope you liked it
too!

## Side notes

For the formatting of the decimal notation, I wanted to figure out the maximum
number of digits a 24-bit number could consist of in base10. <i>@capitol</i> of
[Hackeriet][] helped me out with the formula. Thanks!

```
ceil( 24 * (log(2) / log(10) ) = 8
```

## References
- <https://stackoverflow.com/questions/56606884/how-to-use-random-to-choose-colors>
- <https://stackoverflow.com/questions/5661725/format-ints-into-string-of-hex/19996754#19996754>
- <https://en.wikipedia.org/wiki/Web_colors>
- <https://docs.python.org/3/library/string.html#string.Formatter>
- <https://stackoverflow.com/questions/9641440/convert-from-ascii-string-encoded-in-hex-to-plain-ascii/27519487#27519487>

[Hackeriet]: https://hackeriet.no
[PEP 3101]: https://www.python.org/dev/peps/pep-3101/
[so-question]: https://stackoverflow.com/questions/56606884/how-to-use-random-to-choose-colors
