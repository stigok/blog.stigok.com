---
layout: post
title:  "My jump into golang (notes)"
date:   2019-07-05 16:18:40 +0200
categories: golang notes
---

## Introduction

> Disclaimer: This post is not complete, hence does not promise to be a good
> and valid point of reference in terms of quality for its readers.

I'm learning Go, as one of my latest projects involve doing code review on a
code base written in it. A very interesting language, with a lot of similarities
to other languages I've worked with, like Python, JavaScript and C.

I'm taking "[A Tour of Go][tour]", and taking notes as I go. This post will probably
grow over time. The tour itself is warmly recommended. A very nice way to
learn a new programming language.

## Random notes

- Don't require a broad interface when a narrow one will suffice [1][]

### Print to stderr

```go
fmt.Fprintln(os.Stderr, err)
```

### Implement custom errors

```go
type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string {
    return fmt.Sprintf("Cannot Sqrt negative number: %f", e)
}

func Sqrt(x float64) (float64, error) {
    if x < 0 {
        return 0, ErrNegativeSqrt(x)
    }
    return 0, nil
}
```

### Wrap existing errors to avoid custom errors

```go
import "github.com/pkg/errors"

func ReadConfig() ([]byte, error) {
        home := os.Getenv("HOME")
        config, err := ReadFile(filepath.Join(home, ".settings.xml"))
        return config, errors.Wrap(err, "could not read config")
}
```

[ref](https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully)

### Rot13Decode

```go
func Rot13Decode(encoded []byte) []byte {
    decoded := make([]byte, len(encoded))
    for i, b := range encoded {
        switch true {
        case b >  0x6D: // m
            decoded[i] = b - 13
        case b >= 0x61: // a
            decoded[i] = b + 13
        case b >  0x4D: // M
            decoded[i] = b - 13
        case b >= 0x41: // A
            decoded[i] = b + 13
        }
    }
    return decoded
}
```

### Filter arrays using higher order functions

```go
package main

import "fmt"
import "strings"

func filter(s []string, f func(string) bool) []string {
	var r []string
	for _, v := range s {
		if f(v) == true {
			r = append(r, v)
		}
	}
	return r
}

func main() {
	words := []string{"apple", "orange", "kiwi"}
	ewords := filter(words, func(w string) bool {
		return strings.Contains(w, "e")
	})
	fmt.Printf("Words with an 'e': %v", ewords)
}
```

[ref](https://golangbot.com/first-class-functions/)

## References
- <https://tour.golang.org>
- 1 <https://www.youtube.com/watch?v=29LLRKIL_TI>
- <https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully>

[tour]: https://tour.golang.org
[1]: https://www.youtube.com/watch?v=29LLRKIL_TI
