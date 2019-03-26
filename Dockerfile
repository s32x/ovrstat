# ============================== BINARY BUILDER ==============================
FROM golang:latest as builder

# Copy in the source
COPY . /src
WORKDIR /src

# Dependencies
RUN apt-get update -y && \
    apt-get upgrade -y
RUN GO111MODULE=on go mod vendor

# Vendor, Test and Build the Binary
RUN GO111MODULE=on go mod vendor
RUN go test ./...
RUN CGO_ENABLED=0 go build -o ./bin/server

# =============================== FINAL IMAGE ===============================
FROM alpine:latest

# Dependencies
RUN apk update && apk add --no-cache ca-certificates

# Static files and Binary
COPY --from=builder /src/static /static
COPY --from=builder /src/bin/server /usr/local/bin/server
CMD ["server"]