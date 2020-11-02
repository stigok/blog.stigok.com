---
layout: post
title: "Python 3 assert statements"
date: 2018-02-10 21:36:16 +0100
categories: python
redirect_from:
  - /post/python-3-assert-statements
---

The syntax is a bit different from standard Python functions, as there are no
parens `()` here.

```
assert <expression1>[, <expression2>]
```

The basics are
- Raise an `AssertionError` if `expression1` is *falsy*
  - Pass `expression2` to `AssertionError` if present (`AssertionError(expression2)`)

I am testing to see if this is a good way to assert expectations during
normal program execution.

```python
try:
    assert shouldProcess
    assert len(address) > 0, "Address is required"
    assert type(zip) is num, "Zip code is not a number"
    assert valid_users.get(name) is not None, "User does not exist"
except AssertionError as e:
    logging.error(e)
```

Note that assertion statements **will not run** when Python optimisation
flag is enabled. I.e. `python -O` or `PYTHONOPTIMIZE=1`, so it might be
safer to go about this in a different way for my example use-case.


## References
- <https://docs.python.org/3/reference/simple_stmts.html#assert>
