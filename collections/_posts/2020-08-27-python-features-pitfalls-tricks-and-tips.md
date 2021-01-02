---
layout: post
title:  "Python 3: pitfalls, tips and tricks"
date:   2020-08-27 13:37:11 +0200
categories: python talk draft
processors: pymd
---

I wanted to write a post about some common things I see when reading
and writing Python scripts. I'm also including some thoughts on how to keep
the code clean and readable, and making the most of the tools already built
into the standard Python library.

This post will be continously updated as time passes by. I will eventually
present it as a talk.

## String Interpolation

We have some different methods to choose from when it comes to string
interpolation. Some times it might be nice enough to just concatenate:

### Concatenating Strings with Variables

```python
def html_list(strings):
    """Takes a list of strings and returns a HTML unordered list"""
    html = "<ul>"
    for s in strings:
        html += "<li>" + s + "</li>"
    html += "</ul>"
    return html

print(html_list(["light", "sun", "moon", "water"]))
#eval
```

but when more than two strings are concatenated with variables, it tends
to get messy, hence, decreased readability.

Building an HTML string, like in the example above, but with multiple
interpolation, will quickly get messy if we need to do it a lot. Creating
a string containing a HTML element with multiple properties might look
like something like this:

```python
params = dict(color="#FF9900", value="Orange", hover="A fruit that tastes good")
html = '<span style="display: inline-block; padding: .5em; background-color: ' \
       + params["color"] + ';" title="' + params["hover"] + '">' \
       + params["value"] + '</span>'
print(html)
#eval
```

Since this is a fairly long string, and we already have the params stored
in a `dict`, it's a perfect suit for `str.format()`.

### str.format

This is a good choice when you want to use the same argument more than once in
a single line, or when you want to apply special formatting to the variables.

#### Interpolating values in a `dict`

The same example as in plain concatenation, but now using `str.format`. This
function accepts both positional and keyword arguments, which lets us take
advantage of dictionary item *unpacking*, using the `**` syntax.

```python
params = dict(color="#FF9900", value="Orange", hover="A fruit that tastes good")
html = '<span style="display: inline-block; padding: .5em;' \
       'background-color: {color};" title="{hover}">{value}</span>'
print(html.format(**params))
#eval
```

You can see how this immediately makes this easier to follow by removing
the clutter with open/close quotes and `+` signs.

Other great uses for `str.format()` not covered here:
- binary and hex formatting
- prefixing and suffixing numerical types
- text alignment, indentation and positioning

### f-strings

The *f-strings* comes in handy when writing strings containing multiple
variables of basic types. With basic types, I mean types link `str` and
`int` that hold simple values, in contrast to a `dict` where we will have
to reference the items using the `["key"]` syntax.

```python
thing = "luna"
nickname = "moon"
print(f"There is a {thing} that looks really nice. It's called {nickname}")
#eval
```

I think it makes more sense to use f-strings when writing long, multilined
strings, as you can easily read what variables will be filled in where,
without having to read the arguments passed to e.g. `str.format()`.
They act very similar, but we dont' have to call the function to format
the string, which is a nice bonus.

An example function to create a systemd service file:

```python
def systemd_service(*, binary, description=""):
    return f'''
    [Unit]
    Description={description}

    [Service]
    Type=simple
    ExecStart={binary}
    Restart=always
    RestartSec=10
    '''

svc = systemd_service(binary="/bin/nc -lkuv 127.0.0.1 3000",
                      description="Listen for UDP connections on port 3000")
print(svc)
#eval
```

### % (percent) formatting

Very good to keep syntax clean.

#### Pitfalls

The syntax actually expects a tuple after the `%`, so the following will fail

```python
s = [1, 2, 3]
```

## Cheap Type Safety

### Type assertions

I am not leaning on type safety when writing Python, but some times I still
want to have explicit type checks to keep my head clean at write-time, while
at the same time improving the readability of the function body.

```python
def generate_cert(*, name, hostnames):
    assert isinstance(name, str), "name should be a string"
    assert isinstance(hostnames, list), "hostnames should be a list"

    domains = ",".join(hostnames)
    cmd = f"certbot certonly --domains {domains} --cert-name {hostnames[0]}"
    return subprocess.check_output(cmd)
```

This immediately helps understand what the expected types of the arguments
are, while at the same time throwing `AssertionError`s if the assertions
fail. This can save you time and potential headaches while writing or debugging.

**Warning:** `assert` statements are ignored when `__debug__ == False`,
i.e. when Python is started with optimisations enabled (`python -O`).

## Exceptions

If you are simply re-throwing the exception you've caught, you don't have to
assign it to a variable if you're not going to process it.

Instead of the following:

```python
try:
    throwing_code("some value")
except ValueError as e:
    graceful_cleanup()
    raise e
```

We can instead re-throw the caught exception in the exception handler implicitly by
not passing any arguments to `raise`.

```python
try:
    throwing_code("some value")
except ValueError:
    graceful_cleanup()
    raise
```

See the logging section for neat ways to log exceptions.

## Classes

### Get a `dict` of class members/variables

For this purpose, there's a bultin module [`inspect`](https://docs.python.org/3/library/inspect.html).

For an example class `Page` I have a bound method `render` that uses the class instance variables to
render a template file.

```python
# demo function
render_template = lambda a, b: print(a, b)

class Page:
    site = "blog.stigok.com"

    def __init__(self, title, body):
        self.title = title
        self.body = body

    def render(self):
        render_template("page.tpl", dict(site=self.site, title=self.title, body=self.body))

p = Page(title="hello", body="world")
p.render()
#eval
```

As you can see, I am passing all the members to the `render_template` function.
In this example, the amount of variables is no big deal, but if it's more it might be a problem.
To avoid having to manually pass all the variables to the `render` function, while
at the same time making the `render` method more generic, I can generate a `dict` of the variables
and pass that along instead.

```python
import inspect

# demo function
render_template = lambda a, b: print(a, b)

class Page:
    site = "blog.stigok.com"
    license = "copyleft"
    contact = "blog@example.com"
    foo = 1
    bar = 2

    def __init__(self, title, body):
        self.title = title
        self.body = body

    def render(self):
        # Ignore member names that start with underscore or are bound methods
        data = dict(
            (k, v) for k, v in inspect.getmembers(self) if not k.startswith("_") and not inspect.ismethod(v)
        )
        render_template("page.tpl", data)

p = Page(title="hello", body="world")
p.render()
#eval
```

`inspect.getmembers` will return all members, including Python's builtins, so it needs some filtering
to fit your needs.


## Logging
### String interpolation

Something I often see is using `str.format(...)` or `"%s:%s" % (a, b)`
expressions when using the `logging` library.

```python
import logging

name = "world"
age = 42
logging.error("hello {0:s}, I am {1:d} years old".format(name, age))
#eval
```

However, string formatting is already built into `logging.*` functions,
but without the positional syntax. Very often we don't need to log the
same variable more than once in a single log line, so we can get away
with using the supported `%<type>` syntax.

```python
import logging

name = "world"
age = 42
logging.error("hello %s, I am %d years old", name, age)
#eval
```

### Exceptions

In the following lines we are catching the exception and assigning it to `e`.

```python
import logging

def generate_cert(domain):
    raise ValueError("invalid domain: %s" % domain)

for host in ["blog.stigok.com", "www.stigok.com"]:
    logging.info("Generating cert(s) for %s", host)
    try:
        output = generate_cert(host)
    except ValueError as e:
        logging.error("Failed to generate cert:")
        logging.error(e)
        break
#eval
```

We don't actually have to assign the exception to `e` if we only want to print it.
`logging.exception()` lets us write an error message and automatically adds the
exception message to the log line.

```python
import logging

def generate_cert(domain):
    raise ValueError("invalid domain: %s" % domain)

for host in ["blog.stigok.com", "www.stigok.com"]:
    logging.info("Generating cert(s) for %s", host)
    try:
        output = generate_cert(host)
    except ValueError:
        logging.exception("Failed to generate certs")
        break
#eval
```

This gives us clean code and a clean error message -- with a stack trace,
and we don't assign the `ValueError` to a variable.

Note that `logging.exception()` should only be called inside an `except`
block (exception handler), and that we should not pass the exception itself
as an argument.

```python
# Don't do this
except subprocess.CalledProcessError as e:
    logging.exception("Failed to generate certs", e)
```

## References
- https://docs.python.org/3/library/
