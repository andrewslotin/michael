FROM alpine:3.3
MAINTAINER Andrew Slotin <andrew.slotin@gmail.com>

RUN apk add -U ca-certificates

EXPOSE 8081
ADD slack-deploy-command /bin/server

CMD ["/bin/server", "-h", "0.0.0.0", "-p", "8081"]
