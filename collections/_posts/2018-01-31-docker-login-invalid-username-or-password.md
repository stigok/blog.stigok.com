---
layout: post
title: "Docker login invalid username or password"
date: 2018-01-31 15:47:29 +0100
categories: docker
redirect_from:
  - /post/docker-login-invalid-username-or-password
---

So I was trying really hard to log in with `docker login`

    $ docker login
    Login with your Docker ID to push and pull images from Docker Hub. If you don't have a Docker ID, head over to https://hub.docker.com to create one.
    Username: email@example.com
    Password: 
    Error response from daemon: Get https://registry-1.docker.io/v2/: unauthorized: incorrect username or password

So this is what happened too many times for me now. Even changing password to something without any symbols did not work either. Eventually it was determined that it is in fact **not the email address** I use to log in with it asks for, but my **username**.

You can find your username by logging in to [https://hub.docker.com/](https://hub.docker.com/), then opening up your profile page by using the menu up in the right corner. Underneath the profile picture, the username appears.

    $ docker login
    Login with your Docker ID to push and pull images from Docker Hub. If you don't have a Docker ID, head over to https://hub.docker.com to create one.
    Username: myactualusername
    Password:
    Login Succeeded