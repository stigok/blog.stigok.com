#!/bin/bash
set -ex
docker build --network host -f .deploy/Dockerfile --tag blog.stigok.com .
docker run --rm -it --network=host blog.stigok.com
