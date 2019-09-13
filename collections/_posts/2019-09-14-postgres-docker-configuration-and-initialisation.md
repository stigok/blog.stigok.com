---
layout: post
title:  "Alter Postgres 9.4 configuration options using official postgres docker image"
date:   2019-09-14 00:21:26 +0200
categories: postgres docker
---

I wanted to make my postgres instance log all queries it received.

Postgres 9.4 can be configured using SQL statements like below:

```sql
SELECT set_config('log_statement', 'all', true);
ALTER SYSTEM SET log_statement = 'all';
```

I can place executable *.sh* and *.sql* files in a special folder inside
the container, */docker-entrypoint-initdb.d*, that will be run when the
container database initialises. These scripts will be executed before the
server is listening for requests, the first time it is started.

Mount a configuration file to enable logging of all executed queries:

```terminal
$ echo "ALTER SYSTEM SET log_statement = 'all';" > conf.sql
$ docker run --name psql-test -v $(pwd)/conf.sql:/docker-entrypoint-initdb.d/conf.sql postgres
```

Now, when I execute a query, it shows up in my logs like:

```terminal
$ docker logs psql-test -f
LOG:  execute <unnamed>:
                        SELECT u.uid, u.username
                        FROM users AS u
                        WHERE u.username = $1
                        LIMIT 1
DETAIL:  parameters: $1 = 'sshow'
```

## References
- https://stackoverflow.com/questions/30848670/how-to-customize-the-configuration-file-of-the-official-postgresql-docker-image
- https://hub.docker.com/_/postgres
- https://www.postgresql.org/docs/9.4/runtime-config-short.html
- https://www.postgresql.org/docs/9.4/sql-createuser.html
- https://stackoverflow.com/questions/722221/how-to-log-postgresql-queries/16589931#16589931

