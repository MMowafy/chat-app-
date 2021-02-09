FROM ruby:2.6.3-alpine

RUN apk update && apk upgrade && apk add ruby ruby-json ruby-io-console ruby-bundler ruby-irb ruby-bigdecimal tzdata postgresql-dev && apk add nodejs && apk add curl-dev ruby-dev build-base libffi-dev && apk add build-base libxslt-dev libxml2-dev ruby-rdoc mysql-dev sqlite-dev

RUN mkdir /app
WORKDIR /app

COPY Gemfile Gemfile.lock ./
RUN gem install ovirt-engine-sdk -v '4.3.0' --source 'https://rubygems.org/'
RUN bundle install --binstubs

COPY . .

RUN chmod +x /app/config/run/startup.sh
RUN chmod +x /app/config/run/prepare-db.sh
RUN chmod +x /app/config/run/asset-pre-compile.sh
RUN chmod +x /app/config/run/wait-for-services.sh

EXPOSE 3000

ENTRYPOINT ["sh", "./config/run/startup.sh"]
