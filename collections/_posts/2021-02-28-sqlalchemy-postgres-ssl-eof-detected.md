---
layout: post
title:  "SSL EOF errors with Python 3 SQLAlchemy and managed cloud Postgres database"
date:   2021-02-28 13:08:22 +0100
categories: python
excerpt: Every now and then I would get SSL errors in my Python application using SQLAlchemy with a managed Postgres database in Scaleway.
#proccessors: pymd
---

## Preface

Distributed systems break all the time for various reasons.
My Python web application using SQLAlchemy as the object-relational mapper (ORM),
which also manages connections and connection pools for the database, had connection drops
way too often.

```
psycopg2.OperationalError: SSL SYSCALL error: EOF detected
```

When building distributed systems, it is important to write code that assumes all external dependencies
may fail at any time. I.e. assume that networked services will fail or be unreachable -- in this
case, assume that the cloud DB will be down from time to time.

I am using a managed database in Scaleway which resides in Paris, France, and my webapp is running
in Amsterdam, Netherlands. The reason for mentioning this is to make a point in that a lot of things
can happen to the network between two sites (especially) in different geographical regions.
Servers can be moved, connections can be closed or for unknown reasons just get dropped.
What happens in the data-centers is out of our (end-users) control.

At time of writing, I'm using SQLAlchemy v1.13 with psycopg2-binary v2.8.6.

## SQLAlchemy

In SQLAlchemy there is connection pooling by default. This means that it will create new and re-use
existing connections for new database sessions as necessary. E.g. for a web application I don't
have to re-connect to the database for every incoming request, but only open it for the first request
and re-use it for the next ones, until it either breaks down or is recycled for other reasons.

Ideally you don't want to think about connection handling at all, but just let the ORM do it for you.
With the default settings in SQLAlchemy, you'll get a long way, but especially for servers outside
your control where you don't know the internal network infrastructure you need to take some extra
precautions.

When you create your engine with `sqlalchemy.create_engine` it can take some extra arguments.

```python
sqla_engine = create_engine(
    app.config["myapp.dburi"], client_encoding="utf8", pool_pre_ping=True
)
```

The reasons for the `SSL SYSCALL error: EOF detected` is that the client (ORM) thinks that the TCP
connection is still up, but the server has already hung up without saying so. The client then
starts sending a query down the pipe and when it does, it notices the connection is broken, resulting
in a sudden EOF.

What `pool_pre_ping` does is to test the connection before attempting to execute the actual query.
This comes with an extra round-trip for all queries, but at least in my small-scale application this
doesn't matter at all. Behind the scenes, it sends a query similar to `SELECT 1` to sort of *ping*
the database. If it succeeds it follows up with the actual query you wanted to send -- if it fails
it recycles the connection along with all other connections established *earlier* than the connection
it tried, and establishes a new one before sending the query again.

For this example:

```
# Engine setup
engine = create_engine("postgres://[uri]", client_encoding="utf8", pool_pre_ping=True)

# Session creation helper
Session = sessionmaker(bind=engine)

# Connect to database
session = Session()

# Start transaction with implicit commit
with session.begin():
  post = BlogPost()
  post.title = "hello, world"
  post.text = "this is an example"

  # Insert the row into the database
  session.add(post)
```

Will send something along the following SQL to the database:

```
# Check the connection
SELECT 1
# Then, if it succeeds, follow up with
INSERT INTO blog_posts (title, text) VALUES ('hello, world', 'this is an example')
```

SQLAlchemy calls this *pessimistic disconnect handling* and notes that
> [... this] does not accommodate for connections dropped in the middle of transactions or other SQL operations
but will simplify error handling of dropped connections dramatically for my application.

If the pre-ping operation fails three times in a row, a connection error will bubble up to the surface
normally.

I'm happy with this!

## References
- <https://docs.sqlalchemy.org/en/13/core/pooling.html#disconnect-handling-pessimistic>
