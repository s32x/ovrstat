# ============================== BINARY BUILDER ==============================
FROM golang:1.11.5 as builder

# Copy in the source
WORKDIR /go/src/s32x.com/ovrstat
COPY / .

# Dependencies
RUN apt-get update -y && \
    apt-get upgrade -y
RUN GO111MODULE=on go mod vendor

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/server

# =============================== FINAL IMAGE ===============================
FROM alpine:latest

# Dependencies
RUN apk add --no-cache ca-certificates

# Static files
COPY service/static /service/static

# Binary
COPY --from=builder /go/src/s32x.com/ovrstat/bin/server /usr/local/bin/server
CMD ["server"]