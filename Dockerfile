FROM ruby:2.4


RUN mkdir -p /opt/resource
COPY assets/check.rb    /opt/resource/check
COPY assets/in.rb       /opt/resource/in
COPY assets/out.rb      /opt/resource/out
COPY assets/common.rb   /opt/resource/common.rb
COPY Gemfile            /opt/resource/Gemfile
COPY Gemfile.lock       /opt/resource/Gemfile.lock

WORKDIR /opt/resource
RUN bundle install
