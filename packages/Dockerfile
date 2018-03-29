FROM alpine:latest

RUN apk --update --no-cache add \
  ca-certificates \
  youtube-dl

ADD bin/squirrelbot.docker /usr/local/bin/squirrelbot

CMD ["/usr/local/bin/squirrelbot"]
