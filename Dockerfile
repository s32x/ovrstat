FROM alpine:latest
RUN apk add --no-cache ca-certificates
ADD bin/server /usr/local/bin/server
ADD /web /web/
CMD server