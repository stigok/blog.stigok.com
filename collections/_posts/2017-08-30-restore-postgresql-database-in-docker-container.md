---
layout: post
title: "Restore PostgreSQL database in Docker container"
date: 2017-08-30 11:08:33 +0200
categories: postgres sql docker backup
redirect_from:
  - /post/restore-postgresql-database-in-docker-container
---

I have a backup file `~/tmp/backup.sql` which was created using `pg_dump` from a remote Postgres server. I now want to verify that the backup actually works, and that it can be imported into a new database.

![List tables of a specific PostgreSQL database](https://public.stigok.com/img/1504084167642549943.png)

## Restore database

Create a PostgreSQL server container, and mount the current folder containing the database backup.

    $ cd ~/tmp/backup.sql
    $ docker run --rm -v $(pwd):/backup postgres

Still on the host machine, execute commands on the container with `docker exec`. I created the dump without `CREATE DATABASE` statement, so I will have to create a new target database first. Find the name of the new container with `docker ls`. In this example, the name is *determined_pare*.

    $ docker exec determined_pare psql -U postgres -c "CREATE DATABASE newdb"

You can verify that the database was created by listing all available databases:

    $ docker exec determined_pare psql -U postgres -l

Restore the backed up database to the new target database using the `.sql` file which was mounted earlier at `/backup`

    $ docker exec determined_pare psql -U postgres -d newdb -f /backup/backup.sql

Verify that the tables were indeed created

    $ docker exec determined_pare psql -U postgres -d newdb -c "\dt"

I have also written a [list of useful PostgreSQL commands](https://blog.stigok.com/post/list-of-useful-postgresql-commands) that I want to remember.

## References
- https://stackoverflow.com/questions/12445608/psql-list-all-tables#12455382
- https://stackoverflow.com/questions/10335561/use-database-name-command-in-postgresql#10338367