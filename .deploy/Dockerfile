FROM ubuntu:20.10 AS base
RUN apt-get update \
    && apt-get install -y --no-install-recommends \
               ruby-full build-essential zlib1g-dev \
               git python3 \
    && rm -rf /var/lib/apt/lists/* \
    && groupadd -g 1000 jekyll \
    && useradd -mu 1000 -g jekyll jekyll

USER jekyll
WORKDIR /home/jekyll
ENV GEM_HOME=/home/jekyll/gems \
    PATH="/home/jekyll/gems/bin:${PATH}"
RUN gem install jekyll bundler

ADD Gemfile Gemfile.lock ./
RUN bundle install

ADD . .
RUN bundle exec jekyll build -d _site -t

FROM nginx:1.19-alpine
COPY --from=base /home/jekyll/_site /usr/share/nginx/html
COPY ./.deploy/default.conf /etc/nginx/conf.d/default.conf
