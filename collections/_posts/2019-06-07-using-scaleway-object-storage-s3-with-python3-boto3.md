---
layout: post
title:  "Using Scaleway's Object Storage S3 API with boto3 in Python 3"
date:   2019-06-07 00:11:15 +0200
categories: s3 scaleway
---

## Introduction

I needed to find a library for Python 3 to browse the S3 API of Scaleway's Object Storage.
`boto3` is Amazon's own project, bringing full support for the S3 protocol.

In turn, I'm going to use this to periodically purge old backup files from a backup bucket.

## How to

Install `boto3` using `pip`

```
# pip install boto3
```

The below Python 3 program will list the name of all the current buckets you have.

```
import os
import boto3

session  = boto3.Session(region_name="nl-ams")
resource = session.resource('s3',
                endpoint_url="https://s3.nl-ams.scw.cloud",
                aws_access_key_id=os.getenv("ACCESS_KEY_ID"),
                aws_secret_access_key=os.getenv("SECRET_ACCESS_KEY"))

for bucket in resource.buckets.all():
    print(bucket.name)
```

It appears to be *very important* to specify the `region_name` when initialising
the session object. When not specified, the `endpoint_url` is rewritten into
`https://s3.nl-ams.amazonaws.com/`, and the request will fail.

## References
- https://blog.ruanbekker.com/blog/2018/04/03/using-python-boto3-and-dreamhosts-dreamobjects-to-interact-with-their-object-storage-offering/
- https://boto3.amazonaws.com/v1/documentation/api/latest/reference/core/session.html
