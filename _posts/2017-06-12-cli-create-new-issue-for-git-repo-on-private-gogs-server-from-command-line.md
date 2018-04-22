---
layout: post
title: "CLI create new issue for git repo on private Gogs server from command line"
date: 2017-06-12 13:05:19 +0200
categories: git gogs cli issue
redirect_from:
  - /post/cli-create-new-issue-for-git-repo-on-private-gogs-server-from-command-line
---

I created a utility for creating new issues in git repos on my Gogs installation. It uses `git remote get-url origin` in the current working directory to get the server address of your Gogs server. If origin doesn't point to a gogs instance, but instead is something totally different, I don't know what will happen, but maybe you will find out.

<https://gist.github.com/stigok/fe8056d2672eaa49fdc5490dd10438d0>

## Usage

Enter repository

    $ cd ~/mygitrepos/bar-foo

Add origin if it's not already there

    $ git remote add origin ssh://my-gogs-server:10022/sshow/bar-foo

Create a new issue

    $ gissue Rename repo to foo-bar
    https://my-gogs-server/sshow/bar-foo/issues/1

If everything works out, it prints the URL to the new issue. If it doesn't, hopefully it will print some errors to point you in the right direction. Please leave a comment if you think something should be fixed.