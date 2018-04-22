---
layout: post
title: "gcloud kubectl kompose error message troubleshooting"
date: 2018-02-28 12:14:29 +0100
categories: k8s debugging docker
redirect_from:
  - /post/gcloud-kubectl-kompose-error-message-troubleshooting
---

## kompose

### up

    Error while deploying application: k.Transform failed: Unable to build Docker image for service hostmon: Unable to build image. For more output, use -v or --verbose when converting.: dial unix /var/run/docker.sock: connect: no such file or directory

You haven't started the Docker service. `systemctl start docker`

    DEBU Pushing Docker image 'stigok/hostmon' 
    INFO Pushing image 'stigok/hostmon:latest' to registry 'docker.io' 
    INFO Attempting authentication credentials 'https://index.docker.io/v1/ 
    ERRO Unable to push image 'stigok/hostmon:latest' to registry 'docker.io'. Error: Get https://registry-1.docker.io/v2/: dial tcp: lookup registry-1.docker.io: no such host 
    FATA Error while deploying application: k.Transform failed: Unable to push Docker image for service hostmon: unable to push docker image(s). Check that `docker login` works successfully on the command lin

You have problems with your DNS, or, less probable, Docker.io is down.