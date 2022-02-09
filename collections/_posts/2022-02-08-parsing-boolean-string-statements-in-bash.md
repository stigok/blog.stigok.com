---
layout: post
title:  "Parsing boolean string statements in bash"
date:   2022-02-08 14:57:08 +0100
categories: bash
excerpt: How to convert yes, no, true and false into a boolean in bash.
#proccessors: pymd
---

## Preface

While writing a custom GitHub Action with a Docker entrypoint written in bash,
I want the users to be able to pass boolean variables as `true` and `false`,
`yes` and `no`, `1` and `0`, and en empty string being considered false.

## Using a function with a regular expression

I first went with `return 0` and `return 1` in the function, relying on the
function's exit code, but I felt this increased my cognitive load due to
`0` being success, i.e. true, and `1` being false.
I found that instead echoing a string gives for easier reading.

```bash
# Returns a string `true` if the string is considered a boolean true,
# otherwise `false`. An empty value is considered false.
function str_bool {
  local str="${1:-false}"
  local pat='^(true|1|yes)$'
  if [[ "$str" =~ $pat ]]
  then
    echo 'true'
  else
    echo 'false'
  fi
}
```

I can now configure my variables as booleans!

```bash
enable_debug_logging=$(str_bool "${ENABLE_DEBUG_LOGGING:-}")

if [ "$enable_debug_logging" = "true" ]
then
  # ...
  echo "Debug logging enabled!"
fi
```
