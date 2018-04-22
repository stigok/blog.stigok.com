---
layout: post
title: "Gogs: Key access denied with SSH"
date: 2017-03-28 15:09:03 +0200
categories: gogs ssh git
redirect_from:
  - /post/gogs-key-access-denied-with-ssh
---

## Update (2017-12-28)

After adding a deployment key to a repo, I had to manually go to *Admin panel*, then running *Rewrite '.ssh/authorized_keys' file*.
Then I was successfully authenticated again.

## Original post

I created a new gogs user for my buildbot but I was unable to clone the user's
own repos even though I had added the SSH key of my current UNIX user.

    $ git clone ssh://git@git.example.com/buildbot/foobar.git
    Cloning into 'foobar'...
    Gogs: Key access denied
    fatal: Could not read from remote repository.

    Please make sure you have the correct access rights

I think one of my other gogs users was already registered with the same SSH key
without my knowing. So it worked correctly after I created a new SSH key pair.
(Be aware that the below command **deletes** a user's SSH keys and may be disastrous.)

    $ sudo -u buildbot -i bash -l
    buildbot$ rm ~/.ssh/id_rsa ~/.ssh/id_rsa.pub
    buildbot$ ssh-keygen -b 4069 -t

Then the Gogs profile with the new public key

    buildbot$ cat ~/.ssh/id_rsa.pub

And it all worked as expected again

    $ git clone ssh://git@git.example.com/buildbot/foobar.git
    Cloning into 'foobar'...
    The authenticity of host '[git.example.com] ([ba:be:ca:fe::1337])' can't be established.
    RSA key fingerprint is SHA256:7xm2gVtHc41V0nh8NussuA7lNAzQwh5f9yn+5Mcpdmg.
    Are you sure you want to continue connecting (yes/no)? yes
    Warning: Permanently added '[git.example.com],[ba:be:ca:fe::1337]' (RSA) to the list of known hosts.
    warning: You appear to have cloned an empty repository.
    Checking connectivity... done.

... I just have to confirm before posting a bug.