---
layout: post
title: "Export PostgreSQL table data as JSON with psql"
date: 2017-09-07 22:59:06 +0200
categories: postgres json docker psql
redirect_from:
  - /post/export-postgresql-table-data-as-json-with-psql
---

I had to export some data from a Postgres database into a Mongo database. Postgres apparently (yay!) lets you export table rows as JSON.

> *CREDITS:* This solution comes from the dba.stackexchange answer at the bottom.

First I'm connecting to my remote database using Docker and opening up a psql prompt (`-W` forces psql to ask me for the password)

    docker run -it --rm postgres psql -h $hostname -p $port -U $username -d $dbname -W
    psql (9.6.4, server 9.5.5)
    SSL connection (protocol: TLSv1.2, cipher: ECDHE-RSA-AES256-GCM-SHA384, 
    bits: 256, compression: off)
    Type "help" for help.
    
    dbname=>

> From now on, the psql prompt is denoted as `%`.

If I go straight away to try and export a row to JSON using `row_to_json`, I will get output that is not ideal for exporting:

    % SELECT row_to_json(row) FROM (SELECT id, satisfaction_degree FROM poll LIMIT 3) row;
                row_to_json                
    -------------------------------------------
     {"id":1067,"satisfaction_degree":"FINE"}
     {"id":1068,"satisfaction_degree":"FINE"}
     {"id":1069,"satisfaction_degree":"GREAT"}
    (3 rows)

Turning off output alignment with `\a` and trying again:

    % \a
    Output format is unaligned.
    % SELECT row_to_json(row) FROM (SELECT id, satisfaction_degree FROM poll LIMIT 3) row;
    row_to_json
    {"id":1067,"satisfaction_degree":"FINE"}
    {"id":1068,"satisfaction_degree":"FINE"}
    {"id":1069,"satisfaction_degree":"GREAT"}
    (3 rows)

The header and footer is still there. Turn those off as well with `\t` (Tuples only):

    % \t
    Tuples only is on.
    % SELECT row_to_json(row) FROM (SELECT id, satisfaction_degree FROM poll LIMIT 3) row;
    {"id":1067,"satisfaction_degree":"FINE"}
    {"id":1068,"satisfaction_degree":"FINE"}
    {"id":1069,"satisfaction_degree":"GREAT"}

Great! Now I want to export all rows in table as JSON to a file. Setting output to a file and re-executing the query:

    % \o poll-backup.json
    % SELECT row_to_json(row) FROM (SELECT id, satisfaction_degree FROM poll LIMIT 3) row;

There is no output to the terminal now. Everything is dumped into the output file. Now, **before exiting the container**, open up a new terminal to pull out the outout file from the ephemeral container volume. (`eloquent_clarke` is the name of my psql container)

    $ docker exec -it eloquent_clarke cat /poll-backup.json > ~/tmp/poll-backup.json

Verify the contents of the file:

    $ cat ~/tmp/poll-backup.json
    {"id":1067,"satisfaction_degree":"FINE"}
    {"id":1068,"satisfaction_degree":"FINE"}
    {"id":1069,"satisfaction_degree":"GREAT"}

As you can see, this file is not valid JSON, as it's not an array of objects, but merely a list of JSON objects separated by newlines. It's fine by me. If you have a simple solution to converting the results to valid JSON, please write a comment.

## References
- [https://dba.stackexchange.com/questions/90482/export-postgres-table-as-json](https://dba.stackexchange.com/questions/90482/export-postgres-table-as-json)