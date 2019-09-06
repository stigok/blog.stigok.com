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

### Global variable definition pitfalls

```go
import "fmt"

type CustomInt int

var i CustomInt

func main() {
	fmt.Printf("type %T, value %v, address %p\n", i, i, &i)
	// type main.CustomInt, value 0, address 0x1c4be4

	i = 21
	fmt.Printf("type %T, value %v, address %p\n", i, i, &i)
	//type main.CustomInt, value 21, address 0x1c4be4

	i := 21
	fmt.Printf("type %T, value %v, address %p\n", i, i, &i)
	//type int, value 21, address 0x414030

	i = 7
	fmt.Printf("type %T, value %v, address %p\n", i, i, &i)
	//type int, value 7, address 0x414030
}
```

Taking a closer look at the code above, a custom type `CustomInt` has been
defined. It's an `int` itself, and a variable `i` of type `CustomInt` is
also defined.

In the `main()` function, the type, value and memory address of the objects
are printed. What I found to be causing a bit of trouble for me, is that I
am allowed to do `i := 21`, even though the variable already exists (in the
outer scope), and the `CustomInt` type is replaced with an `int`.

I feel like I'm missing a mesage from the compiler, warning me about
redefining a global variable. It doesn't even care that `i` is unused,
even if it's not exported, like in the below example:

```go
package main

import "fmt"

type CustomInt int

var i CustomInt

func main() {
	i := "Hello"
	fmt.Printf("type %T, value %v, address %p\n", i, i, &i)
}
```

I had done a similar thing in a project i was working with this week, causing
some disturbance in my mind. It looked something like the following. Can you
spot the mistake?

```go
package main

import "fmt"
import "time"

var msgs chan string

func printer() {
	for {
		select {
		case m := <-msgs:
			fmt.Println(m)
		case <-time.After(1 * time.Millisecond):
			fmt.Println("... no new messages arrived")
		}
	}
}

func main() {
	msgs := make(chan string)

	go printer()
	
	func(c chan string) {
		for i := 0; i < 10; i++ {
			msgs <- fmt.Sprintf("Hello %v", i)
		}
	}(msgs)
}
```

[ref](https://www.reddit.com/r/golang/comments/bi2k7o/a_variable_which_is_declared_outside_of_the_main/elxmo12/)

### Useful videos and presentation

- [GopherCon 2019: Mat Ryer - How I Write HTTP Web Services after Eight Years](https://www.youtube.com/watch?v=rWBSMsLG8po)

## References
- <https://tour.golang.org>
- 1 <https://www.youtube.com/watch?v=29LLRKIL_TI>
- <https://dave.cheney.net/2016/04/27/dont-just-check-errors-handle-them-gracefully>

[tour]: https://tour.golang.org
[1]: https://www.youtube.com/watch?v=29LLRKIL_TI
