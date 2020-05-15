FROM nginx:1.18-alpine
COPY . /usr/share/nginx/html
# Path is relative to context, which is the repo root
COPY ./.deploy/default.conf /etc/nginx/conf.d/default.conf
