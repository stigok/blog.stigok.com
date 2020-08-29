---
layout: post
title:  "Capture output of running a shell command in Ruby"
date:   2020-08-29 20:08:16 +0200
categories: jekyll ruby
---

## Preface

I am creating a plugin for Jekyll to perform extra processing of some of my
posts. Specifically, for some of my Python posts, I want to include the actual
output of the Python snippets I've embedded in the Markdown.

I need to
- Create a program to process my Python code
- Create a plugin for Jekyll in Ruby to execute the program
- Come up with a way for my plugin to know which posts to process

This is only about the Ruby command execution.

## Run shell commands from Ruby and capture output

The simplest way to do this is to enclose the command in backticks (`` ` ``), and it will
return a string with the outputs.

```ruby
msg = `echo hello`
puts msg
# Prints: hello
```

## Execute a program from shell with string input to stdin in Ruby

There are some different functions we can use for this purpose in `Open3`,
depending on what control over the input and ouput you need.

I wanted do pass in a string to `STDIN` to simplify my program. That can be
done using the `capture*` variants.

```ruby
cmd = "cat > file.txt"
post_content = "hello world"
stdout_str, exit_code = Open3.capture2(cmd, :stdin_data=>post_content)
```

## Capture stdout and stderr of shell command in Ruby

In the code above I'm only caring about the stdout.
But if you need to capture `STDERR` as well, either forward `STDERR` to `STDIN` in
your shell command

```ruby
stdout_str, exit_code = Open3.capture2("my_program 2>&1")
```

**or** let Ruby take care of exactly the same thing

```ruby
stdout_and_stderr, status = Open3.capture2e("my_program")
```

**or** use `capture3` if you want the outputs separated.

```ruby
post_content = "this input is piped to stdin"
stdout_str, stderr_str, exit_code = Open3.capture3(cmd, :stdin_data=>post_content)
```

Note that the `capture*` functions operates with strings, while `popen*`
operates with pipes. See the [documentation for `Open3`][1] for a complete reference
and list of available functions.

## References
- <https://stackoverflow.com/questions/26601232/preprocessing-markup-files-in-jekyll>
- <https://jekyllrb.com/docs/plugins/generators/>
- <https://stackoverflow.com/questions/2232/how-to-call-shell-commands-from-ruby>
- <https://ruby-doc.org/stdlib-2.4.1/libdoc/open3/rdoc/Open3.html#method-c-capture3>
- <https://jekyllrb.com/docs/plugins/converters/>
- <https://blog.yossarian.net/2019/06/09/Pipelines-in-Ruby-with-Open3>

[1]: https://ruby-doc.org/stdlib-2.4.1/libdoc/open3/rdoc/Open3.html#method-c-capture3
