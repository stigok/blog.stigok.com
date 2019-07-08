---
layout: post
title:  "My jump into golang (notes)"
date:   2019-07-05 16:18:40 +0200
categories: golang
---

## Introduction

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

## References
- <https://tour.golang.org>
- 1 <https://www.youtube.com/watch?v=29LLRKIL_TI>

[tour]: https://tour.golang.org
[1]: https://www.youtube.com/watch?v=29LLRKIL_TI
