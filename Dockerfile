FROM alpine:latest
RUN apk add --no-cache ca-certificates
ADD bin/ovrstat /usr/local/bin/ovrstat
ADD /web /web/
EXPOSE 8080
CMD ovrstat