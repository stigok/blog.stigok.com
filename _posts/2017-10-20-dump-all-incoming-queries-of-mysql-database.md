---
layout: post
title: "Dump all incoming queries of mysql database"
date: 2017-10-20 17:19:47 +0200
categories: mysql debugging
redirect_from:
  - /post/dump-all-incoming-queries-of-mysql-database
---

I have a mysql database I want to dump all incoming queries from. There is a plugin in a Wordpress installation that doesn't do what it's supposed to do, and I want to try and figure out what it *actually* does.

My database server instance is running in a docker container and is only serving a single application. This method might not be suitable for a busy server, or may not even be a good idea to use on a production system.

### Enable global query logging on a running mysql instance

Start a `mysql` shell, logging in as a user with SUPER privileges (e.g. `root`)

    # mysql -u root --password='<password>'
    mysql> SET GLOBAL general_log_file = 'all-queries.log';
    Query OK, 0 rows affected (0.00 sec)

    mysql> SET GLOBAL general_log = 'ON';
    Query OK, 0 rows affected (0.00 sec)

### Monitor executed queries

Use tail to get a streaming log

    # tail -f /var/lib/mysql/all-queries.log

Example output as follows:

    2017-10-20T15:18:22.501267Z	 1615 Connect	wp@172.18.0.4 on  using TCP/IP
    2017-10-20T15:18:22.501645Z	 1615 Query	SET NAMES utf8mb4
    2017-10-20T15:18:22.502003Z	 1615 Query	SET NAMES 'utf8mb4' COLLATE 'utf8mb4_unicode_520_ci'
    2017-10-20T15:18:22.502318Z	 1615 Query	SELECT @@SESSION.sql_mode
    2017-10-20T15:18:22.502639Z	 1615 Query	SET SESSION sql_mode='NO_ZERO_IN_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_AUTO_CREATE_USER,NO_ENGINE_SUBSTITUTION'
    2017-10-20T15:18:22.502945Z	 1615 Init DB	wp
    2017-10-20T15:18:22.504770Z	 1615 Query	SELECT option_name, option_value FROM wp_options WHERE autoload = 'yes'
    2017-10-20T15:18:22.514748Z	 1615 Query	SELECT * FROM wp_users WHERE user_login = 'os'
    2017-10-20T15:18:22.515395Z	 1615 Query	SELECT user_id, meta_key, meta_value FROM wp_usermeta WHERE user_id IN (1) ORDER BY umeta_id ASC
    2017-10-20T15:18:22.518410Z	 1615 Query	SELECT option_value FROM wp_options WHERE option_name = 'active_sitewide_plugins' LIMIT 1
    2017-10-20T15:18:22.522076Z	 1615 Query	SHOW TABLES LIKE 'wp_termmeta'
    2017-10-20T15:18:22.522723Z	 1615 Query	SELECT option_value FROM wp_options WHERE option_name = 'piklist_core' LIMIT 1
    2017-10-20T15:18:22.523365Z	 1615 Query	SELECT DISTINCT meta_key FROM wp_postmeta WHERE meta_key LIKE '\_\_%'
    2017-10-20T15:18:22.523804Z	 1615 Query	SELECT DISTINCT meta_key FROM wp_termmeta WHERE meta_key LIKE '\_\_%'
    2017-10-20T15:18:22.524079Z	 1615 Query	SELECT DISTINCT meta_key FROM wp_commentmeta WHERE meta_key LIKE '\_\_%'
    2017-10-20T15:18:22.524289Z	 1615 Query	SELECT DISTINCT meta_key FROM wp_usermeta WHERE meta_key LIKE '\_\_%'
    2017-10-20T15:18:22.525340Z	 1615 Query	SELECT   wp_posts.* FROM wp_posts  WHERE 1=1  AND wp_posts.ID = 4 AND wp_posts.post_type = 'page'  ORDER BY wp_posts.post_date DESC
    2017-10-20T15:18:22.525890Z	 1615 Query	SELECT post_id, meta_key, meta_value FROM wp_postmeta WHERE post_id IN (4) ORDER BY meta_id ASC
    2017-10-20T15:18:22.530311Z	 1615 Query	SELECT option_value FROM wp_options WHERE option_name = 'can_compress_scripts' LIMIT 1
    2017-10-20T15:18:22.531956Z	 1615 Query	SELECT * FROM wp_posts WHERE ID = 38 LIMIT 1
    2017-10-20T15:18:22.532493Z	 1615 Query	SELECT * FROM wp_posts WHERE ID = 35 LIMIT 1
    2017-10-20T15:18:22.533586Z	 1615 Query	SELECT * FROM wp_posts WHERE ID = 40 LIMIT 1
    2017-10-20T15:18:22.534265Z	 1615 Query	SELECT * FROM wp_posts WHERE ID = 61 LIMIT 1
    2017-10-20T15:18:22.534761Z	 1615 Query	SELECT * FROM wp_posts WHERE ID = 67 LIMIT 1
    2017-10-20T15:18:22.535398Z	 1615 Query	SELECT * FROM wp_posts WHERE ID = 63 LIMIT 1
    2017-10-20T15:18:22.535854Z	 1615 Query	SELECT * FROM wp_posts WHERE ID = 65 LIMIT 1
    2017-10-20T15:18:22.536367Z	 1615 Query	SELECT * FROM wp_posts WHERE ID = 57 LIMIT 1
    2017-10-20T15:18:22.538579Z	 1615 Quit

## References
- https://dev.mysql.com/doc/refman/5.7/en/query-log.html
- https://dev.mysql.com/doc/refman/5.7/en/server-system-variables.html#sysvar_general_log_file