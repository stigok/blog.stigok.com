# stats

Logs visits to `{post}`s. Visits are a hash of the `User-Agent` header and the
remote address of the incoming HTTP request.

Multiple visits to a `{post}` with the same hash will only be counted once.

## Build and run

Build and start the app

```shell
$ go build
$ DATABASE_NAME=/data/sqlite.db LISTEN_ADDRESS=0.0.0.0:8000 ./stats
```

Or build with Docker

```shell
$ docker build . -t stats:latest
$ docker run stats:latest
```

## API

- `GET /visit/{post}/hit`
  - `{post}` is a token to record visits to
  - Returns JSON if `Accept` contains `application/json`
    - `{"success": true, "error": ""}`
  - Returns text if `Accept` contains `text/plain`
  - Returns a 1x1 `image/gif` if `Accept` header does not match any of the above.

- `GET /visit/{post}/get`
  - `{post}` is a token to return visit count for
  - If `{post}` doesn't exist, return a count of `0`
  - Returns JSON if `Accept` contains `application/json`
    - `{"post": "{post}", "count": 0, "success": true, "error": ""}`
  - Returns `text/plain` if `Accept` header does not match any of the above.
