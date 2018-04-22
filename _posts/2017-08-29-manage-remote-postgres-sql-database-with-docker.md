---
layout: post
title: "Manage remote Postgres SQL database with Docker"
date: 2017-08-29 16:16:54 +0200
categories: postgresql postgres pgadmin
redirect_from:
  - /post/manage-remote-postgres-sql-database-with-docker
---

I dont' want to install lots of programs to be able to manage a remote database from my local computer. I already have Docker installed, so I am planning to use an existing image for the database administration software.

Looking around the web, at first I found [adminer](), but after trying the Export feature, it threw some errors at when I imported the backup again to a new database stating:

    PDO::quote() expects parameter 1 to be string

So, instead I will use [pgAdmin](https://www.pgadmin.org/), with an existing [Docker image containing pgAdmin 4](https://github.com/thaJeztah/pgadmin4-docker).

## Getting started

As stated in the readme of the pgAdmin Docker image, we can start the container like this

    $ docker run --rm -p 5050:5050 thajeztah/pgadmin4

Then browse the web interface locally from `http://localhost:5050`. But some features will not work without mounting some additional binaries to the container (e.g. `pg_dump` and `pg_restore`). So I'm downloading all the pre-compiled binaries from [the official postgres website](https://www.enterprisedb.com/download-postgresql-binaries) and extracting them to a user readable folder.

So stopping the container, then performing some steps locally to get what I need:

    $ mkdir -p ~/tmp/postgres
    $ cd ~/tmp/postgres
    $ curl -O https://get.enterprisedb.com/postgresql/postgresql-9.6.4-1-linux-x64-binaries.tar.gz
    $ tar -xf postgresql-9.6.4-1-linux-x64-binaries.tar.gz
    $ ls pgsql/bin/
    clusterdb   ecpg               pg_ctl          pg_restore      pgbench        psql.bin
    createdb    initdb             pg_dump         pg_rewind       pltcl_delmod   reindexdb
    createlang  oid2name           pg_dumpall      pg_standby      pltcl_listmod  vacuumdb
    createuser  pg_archivecleanup  pg_isready      pg_test_fsync   pltcl_loadmod  vacuumlo
    dropdb      pg_basebackup      pg_receivexlog  pg_test_timing  postgres
    droplang    pg_config          pg_recvlogical  pg_upgrade      postmaster
    dropuser    pg_controldata     pg_resetxlog    pg_xlogdump     psql

Now mount the binaries to the container. The volume syntax here is `-v [local path]:[container mount path]`, which means I am mounting the binaries we just extracted into `/pgbin` of the container filesystem.

    $ docker run --rm -v ~/tmp/postgres/pgsql/bin:/pgbin -p 5050:5050 thajeztah/pgadmin4

While the container is running, we can verify that the files were actually mounted by executing `ls` on the container. First checking what the name of the container is:

    $ docker ps
docker ps
    CONTAINER ID        IMAGE                COMMAND                  CREATED             STATUS              PORTS                    NAMES
    908b146e9703        thajeztah/pgadmin4   "python ./usr/loca..."   4 minutes ago       Up 4 minutes        0.0.0.0:5050->5050/tcp   cocky_keller

So now we can execute `ls` on the container

    $ docker exec cocky_keller ls /pgbin
    createdb
    createlang
    droplang
    oid2name
    [...]

Now, to make use of these tools, open up a browser with http://localhost:5050 and go to `File->Preferences->Paths->Binary paths`, and enter `/pgbin` into both path fields.

![pgAdmin 4 enter PostreSQL binary path](https://public.42.fm/1504016039125943576.png)

Now I can import and export my databases.