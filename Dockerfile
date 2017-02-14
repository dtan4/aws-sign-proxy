FROM alpine:3.5

RUN apk add --no-cache --update \
      ca-certificates

COPY bin/aws-sign-proxy /aws-sign-proxy

CMD ["/aws-sign-proxy"]
