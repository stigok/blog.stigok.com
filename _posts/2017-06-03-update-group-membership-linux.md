---
layout: post
title: "Update group membership Linux"
date: 2017-06-03 00:11:01 +0200
categories: linux group
redirect_from:
  - /post/update-group-membership-linux
---

Added myself to a new group, but I need to refresh the groups for the current session to realize the new membership. Logging out and back in again to the current SSH session does not solve it. `newgrp` can help you out. Takes one argument; the group you think you should be a part of:

    # gpasswd -a myusername wheel
    $ group
    myusername
    $ newgrp wheel
    $ group
    myusername wheel