# ============================== BINARY BUILDER ==============================
FROM golang:latest as builder

# Copy in the source
COPY . /src
WORKDIR /src

# Dependencies
RUN apt-get update -y 
RUN apt-get upgrade -y

# Vendor, Test and Build the Binary
RUN go mod vendor
RUN go test ./...
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o ./bin/server

# =============================== FINAL IMAGE ===============================
FROM alpine:latest

# Dependencies
RUN apk update 
RUN apk add --no-cache ca-certificates

# Static files and Binary
COPY --from=builder /src/bin/server /usr/local/bin/server
CMD ["server"]