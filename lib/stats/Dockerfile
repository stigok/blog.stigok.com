# Build container
FROM golang:1.15-buster AS build

WORKDIR /usr/src/app
COPY . .

RUN go get -v .
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o app .

# Runtime container
FROM debian:buster-slim
RUN groupadd --system --gid 999 app \
    && useradd --system --uid 999 --gid 999 --no-create-home app
RUN apt update \
    && apt install --yes --no-install-recommends --no-install-suggests \
       sqlite3

COPY --from=build /usr/src/app/app /usr/local/bin/

USER app
CMD ["app"]
