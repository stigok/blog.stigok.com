---
layout: post
title: "Serve multiple subdomains with different root single config nginx"
date: 2017-06-30 15:57:15 +0200
categories: nginx
redirect_from:
  - /post/serve-multiple-subdomains-with-different-root-single-config-nginx
---

This is a configuration I've been very happy with for serving multiple subdomains with different content with a single nginx configuration file.

    # Subdomains
    server {
            listen 80; 
            listen 443 ssl;
            
            server_name "~^(?<subdomain>\w+)\.example\.com$";
            
            ssl_certificate         /etc/letsencrypt/live/example.com/fullchain.pem;
            ssl_certificate_key     /etc/letsencrypt/live/example.com/privkey.pem;
            
            if ($https != 'on') {
                    return 301 https://$subdomain.example.com$request_uri;
            }
            
            location / { 
                    root /var/www/example.com/subdomains/$subdomain;
                    index index.html;
            }
            
            access_log /var/log/nginx/example.com/access.log;
            error_log  /var/log/nginx/example.com/error.log error;
    }