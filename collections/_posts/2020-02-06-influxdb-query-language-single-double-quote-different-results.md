---
layout: post
title:  "Different query results with single and double quotes in InfluxDB queries"
date:   2020-02-06 11:26:43 +0100
categories: influxdb
---

I was getting different results for seemingly identical queries in `influx` CLI and
in my Go code. The problem appeared to be a difference between single and double quote
handling in the InfluxDB query language.

## Problem

The below query returns an error since `subscription` without any quotes is
actually an InfluxDB keyword. Fair enough -- I received an error!

```
SELECT * FROM azureCosts
WHERE subscription = 'Datahub' AND
      costsReport = 'total' AND
      time = 1580860800000000000
LIMIT 1

ERR: error parsing query: found SUBSCRIPTION, expected identifier, string, number, bool at line 1, char 32
```

This one, single-quoting all tag names and string values returns an empty result
set, but there's no error:

```
SELECT * FROM azureCosts
WHERE 'subscription' = 'Datahub' AND
      'costsReport' = 'total' AND
      'time' = 1580860800000000000
LIMIT 1
```

The next one, double-quoting `subscription` and the string values also returns an
empty result set:

```
SELECT * FROM azureCosts
WHERE "subscription" = "Datahub" AND
      costsReport = "total" AND
      time = 1580860800000000000
LIMIT 1
```

While this last one worked as I intended:

```
SELECT * FROM azureCosts
WHERE "subscription" = 'main' AND
      costsReport = 'total' AND
      time = 1580860800000000000
LIMIT 1
```

The above query describes: get me the first row
where `subscription` tag equals the string value `main` **and** `costReport` tag equals
the string value `total` **and** `time` is equal to the given timestamp, specified in nanoseconds.

After I finally figured it out, I knew [what to search for][search-res]. A blog post from InfluxData
pointed me on to the [InfluxDB v1.7 FAQ][faq], stating the following:

> Single quote string values (for example, tag values) but do not single quote identifiers (database names, retention policy names, user names, measurement names, tag keys, and field keys).

> Double quote identifiers if they start with a digit, contain characters other than [A-z,0-9,_], or if they are an InfluxQL keyword. Double quotes are not required for identifiers if they don’t fall into one of those categories but we recommend double quoting them anyway.

## Debugging tips

You can see incoming queries in the `influxdb` systemd unit logs as they come in

```
# export DB=my_database
# journalctl -fu influxdb | grep --line-buffered $DB | grep -oP 'query=\K.+'
"SELECT * FROM cloudcost_dev.autogen.azureCosts WHERE \"subscription\" = 'Datahub' AND costsReport = 'total'"
```

[search-res]: https://duckduckgo.com/?q=influxdb+single+quote+double+quote
[blogpost-lead]: https://www.influxdata.com/blog/tldr-influxdb-tech-tips-july-21-2016/
[faq]: https://docs.influxdata.com/influxdb/v1.7/troubleshooting/frequently-asked-questions/#when-should-i-single-quote-and-when-should-i-double-quote-in-queries
