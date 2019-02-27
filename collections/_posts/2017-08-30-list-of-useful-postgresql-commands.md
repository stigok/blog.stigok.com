---
layout: post
title: "List of useful PostgreSQL commands"
date: 2017-08-30 11:22:29 +0200
categories: postgres sql
redirect_from:
  - /post/list-of-useful-postgresql-commands
---

From time to time I come across PostgreSQL databases that I have to interact with. But the command syntax is somewhat different than what I'm used to from other databases like MSSQL or MySQL. This is a collection of commands I use, but often forget.

## psql

    $ psql -U user -d dbname -c "# SQL command" -f ./script.sql

### Commands

Select database

    db=# \c databasename

List tables in current database

    db=# \dt

Show all available commands

    db=# \?

### Queries

Change database user password

    ALTER USER user_name WITH PASSWORD 'new_password';