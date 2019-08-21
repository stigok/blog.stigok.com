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

This yielded a list of values similar to this:

```
Ola Nordmann|+4798765432
John Doe|+11800000000
Ola Nordmann|+4784422107
```

Some contacts are returned with multiple numbers, so I want to group them up
and create vCard records for each of my contacts.

I wrote a script in Python to help me with this:

```python
"""
Prints a vCard collection from a file with list of names and numbers.
Expects a CSV (or pipe-separated) file containing a name and a number for each record.


Author: Stig Kolstad (stig@stigok.com), Aug 2019
License: https://creativecommons.org/licenses/by/4.0/


## Example input file contents:

Ola Nordmann|+4798765432
John Doe|+11800000000
Ola Nordmann|+4784422107

## Resulting output string:

BEGIN:VCARD
VERSION:2.1
N:;Ola Nordmann;;;
TEL;CELL:+4798765432
TEL;CELL:+4784422107
END:VCARD
BEGIN:VCARD
VERSION:2.1
N:;John Doe;;;
TEL;CELL:+11800000000
END:VCARD
"""
import argparse
import os
from collections import defaultdict

def main():
    parser = argparse.ArgumentParser()
    parser.add_argument('filename', help="Separated values file")
    parser.add_argument('-s', '--separator', help="Value separator (default |)", default='|')
    args = parser.parse_args()

    FILE = args.filename
    SEP  = args.separator

    contacts = defaultdict(list)

    # Collect all numbers for each distinct name into lists
    with open(FILE, mode='r') as f:
        for line in f:
            name, number = line.strip().split(SEP)
            contacts[name].append(number)

    # Output VCF formatted records for each contact
    for name, numbers in contacts.items():
        print(vcfstr(name, numbers))


def vcfstr(name=None, numbers=None):
    """Returns a VCF vCard 2.1 formatted string"""
    f = [ "BEGIN:VCARD", "VERSION:2.1", "N:;%s;;;" % name ]
    for n in numbers:
        f.append("TEL;CELL:%s" % n)
    f.append("END:VCARD")

    return os.linesep.join(f)

if __name__ == "__main__":
    main()
```

Then I ran this script with the results returned from the sqlite query and saved it to a vcf file.

```
$ python vcf.py contacts.txt > old-contacts.vcf
$ adb push old-contacts.vcf /data/Download/
```

Then imported it on my phone again, through the settings menu in the Contacts app

![Contact app import prompt](https://public.stigok.com/img/2019-08-21-023611.png)

After a little second, the contacts were successfully imported!

![Import successful](https://public.stigok.com/img/2019-08-21-024103.png)

## References
- <https://en.wikipedia.org/wiki/VCard#vCard_2.1>
