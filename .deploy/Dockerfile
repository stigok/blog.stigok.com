FROM nginx:alpine

# Path is relative to context, which is the repo root
COPY ./build/ /usr/share/nginx/html
COPY ./.deploy/default.conf /etc/nginx/conf.d/default.conf
