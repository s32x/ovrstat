FROM alpine:latest
RUN apk add --no-cache ca-certificates
ADD bin/server /usr/local/bin/server
CMD server