#!/bin/bash
set -ex
docker build --network host -f .deploy/Dockerfile --tag blog.stigok.com .
docker run --rm -it -p 127.0.0.1:8080:80 blog.stigok.com
