---
layout: post
title:  "A Python object type that supports both dict and dot value getter"
date:   2020-09-20 23:07:10 +0200
categories: python
excerpt: I wanted to be able to access dict values using dot notation.
processors: pymd
---

## Preface

Having a dict, I am unable to access its keys using dot notation

```python
d = dict(a="one", b="two")
print(d.a)
#eval
```

And to get strict dot notation, I could do something like this, using a class.
But now I can't access the values using the normal dict value getter:

```python
class DictToDotNotation:
  def __init__(self, d):
    for k, v in d.items():
      setattr(self, k, v)

d = dict(a="one", b="two")
obj = DictToDotNotation(d)
print(obj.a)
print(obj.b)
print(obj["a"])
#eval
```

So, I want something in between. A hybrid type of sorts.

## Overriding default behavior

Python exposes most its internal workings to the developer. All *unders* and
*dunders* are free to override in subclasses, which will help out here quite
a bit.

```python
class DictHybrid(object):
    """
    A dict-like object used to allow both dot and array notation to access its keys.
    Accessing an undefined key using dot notation returns `None`, while using
    traditional notation `[key]` will raise a `KeyError`.
    """

    # TODO: Can add default_data here too as a ChainMap to allow for values
    # overridden by parent templates
    def __init__(self, **kwargs):
        self._dict = dict(**kwargs)

    def __getattr__(self, name):
        """Support dot notation value access"""
        return self._dict.get(name)

    def __getitem__(self, name):
        """Support dict style value access"""
        return self._dict[name]


d = DictHybrid(a="foo")
print(d.a)
print(d.b)
print(d["b"])
#eval
```

If you want to define other behaviors for dot and dict notation, it's very straight
forward to do so. E.g. if you want to raise `KeyError`s for dot notation as well.

When I figured this out, I finally closed the gap between `dict` and
`namedtuple`, something that has been bothering me for a long time.

Please note that I have not implemented *setters* in my implementation.

## Existing solutions

After looking around for similar implementations, I found a great [post on
StackOverflow][1]. I advised you to go there for some great alternatives,
including setters.


## References
- <https://stackoverflow.com/questions/2405590/how-do-i-override-getattr-in-python-without-breaking-the-default-behavior>
- <https://github.com/python/cpython/blob/master/Lib/collections/__init__.py>
- <https://stackoverflow.com/questions/2352181/how-to-use-a-dot-to-access-members-of-dictionary>

[1]: https://stackoverflow.com/questions/2352181/how-to-use-a-dot-to-access-members-of-dictionary
