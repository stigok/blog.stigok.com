---
layout: post
title:  "Using watchdog and sd-notify functionality for systemd in Python 3"
date:   2020-01-26 10:42:31 +0100
categories: systemd python
---

I [published my first package to Python Package Index (PyPI)][pkg] today after my
package has been tested for a while.

Publishing packages to PyPI was pretty simple using the [official guide][pkg-guide],
but [I created a Makefile][makefile] to help me out the next time around.

The package itself is a tiny client library to interface with the systemd watchdog functionality
and [sd-notify(3)][manpage]. I wrote a [similar post for a C integration on the Hackeriet blog][hackeriet-blog]
some time ago.

> A systemd service of Type=notify waits for the executable program to send a notification message to systemd before it is considered activated. Up until the service is active, its state is starting. systemctl start <svc> will block until the service is active, or failed.
>
> Similarly, a service which has WatchdogSec set will expect to receive a notification message no less than at every specified time interval. If no message has been received, systemd will kill the process with SIGABRT and place the service in a failed state.

## Example

Install the library

```
$ pip install sd-notify
```

Create a Python script file with the following code taken from the project readme.

```
# /home/sshow/tmp/sd-notify-test.py
import time
import sd_notify

notify = sd_notify.Notifier()
if not notify.enabled():
    # Then it's probably not running is systemd with watchdog enabled
    raise Exception("Watchdog not enabled")

# Report a status message
notify.status("Initialising my service...")
time.sleep(3)

# Report that the program init is complete
notify.ready()
notify.status("Waiting for web requests...")
time.sleep(3)

# Report an error to the service manager
notify.notify_error("An irrecoverable error occured!")
# The service manager will probably kill the program here
time.sleep(3)
```

Then create a [systemd.service][] file to run and watchdog this script. I created
this in my systemd user directory.

```
$ systemctl --user cat sd-notify-test.service
# /home/sshow/.config/systemd/user/sd-notify-test.service
[Unit]
Description=Watchdog test service

[Service]
ExecStart=/usr/bin/python3 /home/sshow/tmp/sd-notify-test.py
Type=notify
WatchdogSec=15
Restart=on-failure
RestartSec=10

[Install]
WantedBy=default.target
```

Then start the service

```
$ systemctl --user daemon-reload
$ systemctl --user start sd-notify-test
```

You can follow the log output of the user service with `journalctl`

```
$ journalctl --user-unit sd-notify-test
jan. 26 11:11:04 systemd[620]: Stopped Watchdog test service.
jan. 26 11:11:04 systemd[620]: Starting Watchdog test service...
jan. 26 11:11:07 systemd[620]: Started Watchdog test service.
jan. 26 11:11:10 systemd[620]: sd-notify-test.service: Watchdog request (last status: An irrecoverable error occured!)!
jan. 26 11:11:10 systemd[620]: sd-notify-test.service: Killing process 358538 (python3) with signal SIGABRT.
jan. 26 11:11:10 systemd[620]: sd-notify-test.service: Main process exited, code=dumped, status=6/ABRT
jan. 26 11:11:10 systemd[620]: sd-notify-test.service: Failed with result 'watchdog'.
jan. 26 11:11:10 systemd-coredump[358542]: Process 358538 (python3) of user 1000 dumped core.

                                           Stack trace of thread 358538:...
```

Another way to look at a service status message is by running `systemctl`

```
$ systemctl --user status sd-notify-test.service
● sd-notify-test.service - Watchdog test service
     Loaded: loaded (/home/sshow/.config/systemd/user/sd-notify-test.service; disabled; vendor preset: enabled)
     Active: active (running) since Thu 2020-01-26 11:15:46 CET; 1s ago
   Main PID: 358886 (python3)
     Status: "Waiting for web requests..."
     CGroup: /user.slice/user-1000.slice/user@1000.service/sd-notify-test.service
             └─358886 /usr/bin/python3 /home/sshow/tmp/sd-notify-test.py

jan. 26 11:15:43 <hostname> systemd[620]: Stopped Watchdog test service.
jan. 26 11:15:43 <hostname> systemd[620]: Starting Watchdog test service...
jan. 26 11:15:46 <hostname> systemd[620]: Started Watchdog test service.
```

If you have any issues, please post to the [issue tracker][issues] of the [source repository][repo].

[pkg]: https://pypi.org/project/sd-notify/
[pkg-guide]: https://packaging.python.org/tutorials/packaging-projects/
[makefile]: https://github.com/stigok/sd-notify/commit/cef51ccd7edfe882e8c624f3aadeceafffeccabf#diff-b67911656ef5d18c4ae36cb6741b7965
[manpage]: http://man7.org/linux/man-pages/man3/sd_notify.3.html
[hackeriet-blog]: https://blog.hackeriet.no/systemd-service-type-notify-and-watchdog-c/
[systemd.service]: http://man7.org/linux/man-pages/man5/systemd.unit.5.html
[issues]: https://github.com/stigok/sd-notify/issues
[repo]: https://github.com/stigok/sd-notify
