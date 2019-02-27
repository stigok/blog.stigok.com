---
layout: post
title: "Lost password to gogs git server docker instance"
date: 2017-03-06 17:36:33 +0100
categories: gogs git docker sqlite
redirect_from:
  - /post/lost-password-to-gogs-git-server-docker-instance
---

I lost my password to my gogs instance, a private git server with a nice github-like web interface. When generating a password for my user my password manager didn't record it, but due to my browser session settings I never needed to log in again.

Then the day came when I cleared my cookies and wanted to log back in. No password to be found. And to make things a bit harder I had disabled the password reset function. My past self must have done this to be smart. Using the git functionality itself, however, is not a problem, as I'm using SSH keys.


## Finding out how the password is stored

My gogs instance is in a docker container and the configuration and data is persisted on a docker volume.

First I want to get the ID or name of the container

{% raw %}
    $ docker ps --format 'table {{.ID}} {{.Names}}' | grep gogs
    1b3c2f0502ce gogs
{% endraw %}

Then see where the data volume is located

    $ docker inspect 1b3c2f0502ce | grep data
    "gogs-data:/data:rw"
    "Name": "gogs-data",
    "Source": "/var/lib/docker/volumes/gogs-data/_data",
    "Destination": "/data",
    "GOGS_CUSTOM=/data/gogs"
    "/data": {}

> All commands starting with # are run as root.

Looking into the `Source` location helped me find an `app.ini` file.

    # ls -al /var/lib/docker/volumes/gogs-data/_data/gogs/conf/
    app.ini

The `database` directive tells me it's using sqlite

    [database]
    DB_TYPE  = sqlite3
    HOST     = 127.0.0.1:3306
    NAME     = gogs
    USER     = root
    PASSWD   = secret
    SSL_MODE = disable
    PATH     = data/gogs.db

So now I'm gonna try and browse that database instance. Expecting
to find a hashed password for my user. My plan is then to create
a new password, hash it with the same algorithm, then log in
successfully.

## Create a new password for my user

The container most likely contains the binaries I need to do this.

    # docker exec -it gogs find / -iname "*sql*"
    /app/gogs/models/models_sqlite.go
    /app/gogs/public/plugins/codemirror-5.17.0/mode/sql
    /app/gogs/public/plugins/codemirror-5.17.0/mode/sql/sql.js

... No. So I will have to get those binaries myself on my local machine, then use them on the docker volume files.

    # apt-get update && apt-get install sqlite3

Now, I want to back up the database in case I do something really stupid.

    # cp /var/lib/docker/volumes/gogs-data/_data/gogs/data/gogs.db ~/gogs-backup.db

Then start looking around.
Commands that helped me find the table and column names in the sqlite3 CLI was `.table` to get all table names, then `.schema user` to get the table columns of the user table.

    # sqlite3 ./gogs.db
    > SELECT name,passwd,salt FROM user;
    olduser|540416a8b535194150bca3e822c73034646a313839333479666e736476393181b68a9e180992e92af0e846cf95aaa254fb0c|N1DEVGrAYa

I'm getting a hash here, but I'm not sure what hashing method is used. Think I should
[check the source code](https://github.com/gogits/gogs/blob/cd15a1797076d97261c1fd9e9ebf50fd4bb76a4c/models/user.go#L328). 

    newPasswd := pbkdf2.Key([]byte(u.Passwd), []byte(u.Salt), 10000, 50, sha256.New)

I will try to replicate that hashing scheme with `mkpasswd`, although I do not
know what to do with the `50` param here. May be that mkpasswd is the wrong tool. Since 50 might mean hash length. But let's try anyway.

    $ mkpasswd --salt=S3CRE75ALT --rounds=10000 --method=sha-256
    $5$rounds=10000$wer1ush3kf$SsXjjNLiHp3R3/YvSQK4CnYEjDEeA1diPH3WqQcaEV0

The resulting hash doesn't look like the strings in the database. Some kind of
encoding has to be applied. The old hash looks to be hex, so lets try encode into that.
I don't know how to convert string to hex in bash, so I'm using Node.js instead.

    $ echo "$5$rounds=10000$wer1ush3kf$SsXjjNLiHp3R3/YvSQK4CnYEjDEeA1diPH3WqQcaEV0" \
      | node -e "process.stdin.on('data', data => console.log(new Buffer(data).toString('hex')))"

    3d31303030302f597653514b34436e59456a4445654131646950483357715163614556300a

It's starting to look better, but the resulting string isn't long enough. I'm starting to think I could  decode the original hash instead to get some more clues. That is, if it's simply base64 encoded. After some searching, found out how to do hex conversion in bash using `xxd`.

    $ echo "540416a8b535194150bca3e822c73034646a313839333479666e736476393181b68a9e180992e92af0e846cf95aaa254fb0c" | xxd -r -p
    ��5AP���"�04�i����F� �d4ǁ���	��*��Fϕ��T�

It appears I was wrong, so looking for another approach.

## Create a new user and use that users hash and salt for my old one

Now I think it might be a better idea to enable new user account registration in the gogs config, create a new user then use the hash and salt from that user and copy it to my old one. I remember seeing this flag in the config when I was looking for the database connection info earlier. Seems to be path of
least resistance right now.

In the app.ini file, I'm setting `DISABLE_REGISTRATION` to `false` and restarting
the container.

    $ docker restart gogs

Logged in to the web interface and created a new user before copying the new hash
and salt from the new user to my old one.

    sqlite> SELECT name,passwd,salt FROM user WHERE name='testuser';
    testuser|6ed26719fe986fea3a9f997ab3c7853ac64ce5d0f5f278dbb571a50a2a2509a4540dd353506cf1f4b9b6211b04faa1a1fc28|97ZaYtIZwM

    sqlite> UPDATE user
    ...> SET passwd='6ed26719fe986fea3a9f997ab3c7853ac64ce5d0f5f278dbb571a50a2a2509a4540dd353506cf1f4b9b6211b04faa1a1fc28',
       ...> salt='97ZaYtIZwM'
       ...> WHERE name='user';

This worked like a charm.

    const path = "/trial/and/error"
