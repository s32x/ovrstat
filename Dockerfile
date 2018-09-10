FROM alpine:latest
RUN apk add --no-cache ca-certificates
ADD bin/ovrstat /usr/local/bin/ovrstat
ADD /web /web/
CMD ovrstat