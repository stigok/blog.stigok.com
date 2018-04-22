---
layout: post
title: "Python 3 assert statements"
date: 2018-02-10 21:36:16 +0100
categories: python
redirect_from:
  - /post/python-3-assert-statements
---

The syntax is a bit differenf from other Python functions;

    num = 2
    assert num == 1, "Message to print when assertion failed"

I am trying out a new way for me to use assertion statements in normal code flow;

    try:
        assert len(address) > 0, "Address is required"
        assert type(zip) is num, "Zip code is not a number"
        assert valid_users.get(name) is not None, "User does not exist"
    except AssertionError as e:
        log.error(e)


## References
- https://docs.python.org/3/reference/simple_stmts.html#assert