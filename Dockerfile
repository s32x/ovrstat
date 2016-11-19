FROM alpine:3.4
RUN apk add --no-cache ca-certificates
ADD ovrstat /usr/local/bin/ovrstat
EXPOSE 7000
CMD ["ovrstat"]
