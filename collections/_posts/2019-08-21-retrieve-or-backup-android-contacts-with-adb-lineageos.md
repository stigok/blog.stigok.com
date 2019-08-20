---
layout: post
title:  "Retrieve or backup Android (LineageOS) contacts with adb"
date:   2019-08-21 00:17:33 +0200
categories: android lineageos sqlite
---

I have a broken phone that doesn't register touch events anymore.
I have a new phone now, but the old one has all the phone numbers.
I tried to find a way to get the contacts out of there using `adb`.

For adb to work, *adb debugging* must be (or have been, in my case)
enabled on the device. To enable developer settings, touch "Build number"
in "About phone" 15 times. (See elsewhere for more information)

I needed root shell to access the `/data` folder on the device.

```
$ adb root
$ adb shell
```

Then in the root shell, I tried to find some files related to contacts

```
# find . -iname *contact* 2>/dev/null
[...]
./data/data/com.android.providers.contacts/databases/contacts2.db
[...]
```

Opening the *.db* file in a text editor showed that it was a sqlite3 database file.
I backed it up with `adb pull`, then opened it up with `sqlite3`.

```
$ adb pull ./data/data/com.android.providers.contacts/databases ~/contacts
$ sqlite3 -readonly ~/contacts/contacts2.db
```

In the sqlite shell, you can use `.tables` to list all tables in the database,
and `.schema [name]` to dump the schema.
The interesting tables (as far as I could find) were `contacts_raw` and `phone_lookup`.
I created a query to dump all the names and numers in the database.

```sqlite
SELECT DISTINCT c.display_name, n.normalized_number
FROM raw_contacts AS c
LEFT JOIN phone_lookup AS n
ON c._id=n.raw_contact_id
ORDER BY c.display_name;
```

Now I had a list of all the numbers. How I'll go about importing them, I
don't yet know. Maybe manually...
