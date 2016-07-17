FROM golang
MAINTAINER Steven Wolfe "steven@swolfe.me"

# Add main src files
ADD . /go/src/github.com/sdwolfe32/ovrstat

# Install lib dependancies
RUN go get github.com/sdwolfe32/ovrstat/goow

# Install core dependancy library
RUN go install github.com/sdwolfe32/ovrstat/goow

# Get and install the rest
RUN go get github.com/sdwolfe32/ovrstat

# Install API
RUN go install github.com/sdwolfe32/ovrstat

ENTRYPOINT /go/bin/ovrstat
EXPOSE 7000
