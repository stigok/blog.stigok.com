---
layout: post
title:  "Buffered read split by two consecutive newlines using bufio.Scanner in Go"
date:   2020-12-04 22:29:38 +0100
categories: go golang adventofcode2020
excerpt: Normal bufio.NewScanner splits by single newline. This splits by double newlines.
#proccessors: pymd
---

## Preface

During the [Advent of Code 2020](https://adventofcode.com/2020), some of the challenges have so
far contained input data separated by two consecutive newlines. I chose to try and figure out
how to make use of `bufio.Scanner` to give me the multiline items one by one.

## Golang's `bufio.Scanner`

The builtin `bufio.NewScanner` uses `bufio.ScanLines` that buffers the output
by a single newline. Meaning that given a `Reader` interface reading from
a stream containing

```
hello, world
foo bar
baz 42
```

will yield three lines over three consecutive calls to `scanner.Scan()` using the following code:

```go
text := `hello, world
foo bar
baz 42`
r := strings.NewReader(text)
scanner := bufio.NewScanner(r)
i := 1
for scanner.Scan() {
  fmt.Println(i, scanner.Text())
  i++
}
```
<small>Code from <https://golang.org/pkg/bufio/#example_Scanner_lines></small>

However, when you'd want to split a string by two consective newlines instead, a new `bufio.Scanner`
must be written.

## Writing a custom bufio.Scanner

I first took a look at [the original implementation](https://github.com/golang/go/blob/2975b27bbd3d4e85a2488ac289e112bc0dedfebe/src/bufio/scan.go#L336-L364)
It's much easier to write a new implementation when having a full working example,
especially when it already behaves almost exactly as I'd want it to.

```go
// dropCR drops a terminal \r from the data.
func dropCR(data []byte) []byte {
  if len(data) > 0 && data[len(data)-1] == '\r' {
    return data[0 : len(data)-1]
  }
  return data
}

// ScanLines is a split function for a Scanner that returns each line of
// text, stripped of any trailing end-of-line marker. The returned line may
// be empty. The end-of-line marker is one optional carriage return followed
// by one mandatory newline. In regular expression notation, it is `\r?\n`.
// The last non-empty line of input will be returned even if it has no
// newline.
func ScanLines(data []byte, atEOF bool) (advance int, token []byte, err error) {
  if atEOF && len(data) == 0 {
    return 0, nil, nil
  }
  if i := bytes.IndexByte(data, '\n'); i >= 0 {
    // We have a full newline-terminated line.
    return i + 1, dropCR(data[0:i]), nil
  }
  // If we're at EOF, we have a final, non-terminated line. Return it.
  if atEOF {
    return len(data), dropCR(data), nil
  }
  // Request more data.
  return 0, nil, nil
}
```

Updating the function to split by two consecutive newlines instead shouldn't be too hard.
Additionally, I want it to join the single newlines together, so that it yields
a concatenated string without newlines when calling `scanner.Text()`. Meaning that `foo\nbar\n\nbaz`
yields two strings; `foo bar` and `baz`. Hence, standalone newlines replaced with
spaces.

To make things easier for myself, I am using regex lookup instead of slice indices. This lets
me replace both `\n` and `\r` in a single operation, while at the same time giving me the
indices that I need to return on each invokation. It will probably be more expensive
perfomance-wise, but computers are fast and my input is short, so I shouldn't worry about that.

This led me to the following implementation:

```go
var (
	patEols  = regexp.MustCompile(`[\r\n]+`)
	pat2Eols = regexp.MustCompile(`[\r\n]{2}`)
)

// Modified version of Go's builtin bufio.ScanLines to return strings separated by
// two newlines (instead of one). Returns a string without newlines in it, and trims
// spaces from start and end.
// https://github.com/golang/go/blob/master/src/bufio/scan.go#L344-L364
func ScanTwoConsecutiveNewlines(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if loc := pat2Eols.FindIndex(data); loc != nil && loc[0] >= 0 {
		// Replace newlines within string with a space
		s := patEols.ReplaceAll(data[0:loc[0]+1], []byte(" "))
		// Trim spaces and newlines from string
		s = bytes.Trim(s, "\n ")
		return loc[1], s, nil
	}

	if atEOF {
		// Replace newlines within string with a space
		s := patEols.ReplaceAll(data, []byte(" "))
		// Trim spaces and newlines from string
		s = bytes.Trim(s, "\r\n ")
		return len(data), s, nil
	}

	// Request more data.
	return 0, nil, nil
}
```

Take a look at my [utility repository](https://github.com/stigok/go-utils) for updated code, along with [a test for this function](https://github.com/stigok/go-utils/blob/main/bufio_test.go).

## References
- <https://godoc.org/bufio>
