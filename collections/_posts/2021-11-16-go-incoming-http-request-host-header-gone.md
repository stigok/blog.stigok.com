---
layout: post
title:  "Go incoming http.Request - the Host header is gone!"
date:   2021-11-16 23:39:56 +0100
categories: golang
excerpt: I thought Heroku messed up the Host header of incoming requests. I was wrong.
#proccessors: pymd
---

## Preface

I created a `http.Handler` to redirect from HTTP to HTTPS for an app I'm running in Heroku.
I am relying on the `X-Forwarded-Proto` to determine the protocol.

```go
func (app *config) httpsOnlyHandler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        proto := r.Header.Get("X-Forwarded-Proto")
        host := r.Header.Get("Host")

        if app.ForceHttps && proto != "https" {
                http.Redirect(w, r, fmt.Sprintf("https://%s%s", host, r.URL.RequestURI()), http.StatusSeeOther)
                return
        }

        next.ServeHTTP(w, r)
    })
}
```

All nice and dandy, and my test (apparently flaky) tells me my handler is working fine:

```go
func TestHttpsOnlyHandler(t *testing.T) {
    okHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("ok"))
    })

    t.Run("redirects on http request", func(t *testing.T) {
        is := is.New(t)

        w := httptest.NewRecorder()
        r, err := http.NewRequest("GET", "/", nil)
        is.NoErr(err)

        r.Header.Set("X-Forwarded-Proto", "http")

        app := config{ForceHttps: true}
        handler := app.httpsOnlyHandler(okHandler)
        handler.ServeHTTP(w, r)
        resp := w.Result()

        is.Equal(resp.StatusCode, http.StatusSeeOther)
    })
}
```

Then I pushed the app to Heroku, brave as I am. I wrote a test! It's passing! So why not
push it straight out to production? In retrospect the answer is obvious: because chances
are you'll hit problems you'd otherwise discovered in a staging environment.

The app broke as it attempted to redirect `http://domain.tld/?user=Z1sVTkHYy2S7Alngeen77EXz` to
`https://?user=Z1sVTkHYy2S7Alngeen77EXz`.

I spent time searching for *heroku incoming host header* to figure out why on earth they
would remove the incoming `Host` header, or what I'd have to do to make it appear again.
Of course they didn't do this, but my assumptions went haywire the late evening.

I read the whole article for [Heroku HTTP Routing](https://devcenter.heroku.com/articles/http-routing),
which was a good and informative read, but didn't mention anything about `Host` headers.

I spun up another app, thankfully with another HTTP library, or else I probably wouldn't
have noticed, that the `Host` header was present on all incoming requests.
What was I missing?

I went to double check the documentation for `http.Request.Redirect` to see if I somehow
mixed up the function call arguments -- I didn't.

I continued the read downwards to the [struct documentation for `http.Request`](https://pkg.go.dev/net/http#Request)
and read all about `http.Request.Header` where a comment appeared:

```go
// For incoming requests, the Host header is promoted to the
// Request.Host field and removed from the Header map.
```

And I'm like: **why**? What a bad surprise. At least now I know. It was a --good-- reminder
to always read the docs, but somehow I feel it wasn't on me this time, as I think this is
a strange thing to do. I assume they have their reasons and move on after this rant.

The reason it worked in the other app I spun up was that I used `github.com/gofiber/fiber`
where they do not use this practice, and instead exposes the `Host` header in both
`fiber.Context.Hostname()` and in `fiber.Context.Get("Host")`.

## Solution

To get the `Host` header from an incoming HTTP request, use `r.Host`. A call to
`r.Header.Get("Host")` will always return an empty string for an incoming request.

```go
func (app *config) httpsOnlyHandler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        proto := r.Header.Get("X-Forwarded-Proto")

        if app.ForceHttps && proto != "https" {
            http.Redirect(w, r, fmt.Sprintf("https://%s%s", r.Host, r.URL.RequestURI()), http.StatusSeeOther)
            return
        }

        next.ServeHTTP(w, r)
    })
}
```

To write a valid test for this, make sure to assert that the contents of the `Location` header
is as expcted, in addition to the HTTP redirect status code.

```go
t.Run("redirects on http request", func(t *testing.T) {
    is := is.New(t)

    w := httptest.NewRecorder()
    r, err := http.NewRequest("GET", "/foo?bar=baz", nil)
    is.NoErr(err)

    r.Header.Set("X-Forwarded-Proto", "http")
    r.Host = "example.com" // <--

    app := config{ForceHttps: true}
    handler := app.httpsOnlyHandler(okHandler)
    handler.ServeHTTP(w, r)
    resp := w.Result()

    is.Equal(resp.StatusCode, http.StatusSeeOther)
    is.Equal(resp.Header.Get("Location"), "https://example.com/foo?bar=baz") // <--
})
```

## References
- <https://pkg.go.dev/net/http#Request>
