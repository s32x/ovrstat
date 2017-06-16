FROM alpine:latest
RUN apk add --no-cache ca-certificates
ADD bin/ovrstat /usr/local/bin/ovrstat
CMD ["ovrstat"]
